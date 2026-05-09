package userstore

import (
	"context"
	"database/sql"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/pkg/problem"
)

const Table = `
create table if not exists Users (
	uuid     blob primary key,
	siape    text unique not null,
	name     text not null,
	email    text not null,
	password blob not null,
	role     text not null,
	created  datetime not null,
	updated  datetime not null
);`

const Indexes = `create unique index if not exists Users_siape on Users(siape);`

const Views = `
create view if not exists Users_View as
	select u.uuid, u.siape, u.name, u.email, u.password, iif(p.uuid is null, u.role, 'promoted-admin') as 'role', s.uuid is not null as 'logged', u.created, u.updated
	from Users u
	left join Sessions s on s.user = u.uuid
	left join Promotions p on p.user = u.uuid;`

type SQLDB struct {
	db store.DBTx
}

var _ user.Store = (*SQLDB)(nil)

func New(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]user.Record, error) {
	rows, err := s.db.QueryContext(ctx, list)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	var recs []user.Record
	for rows.Next() {
		var rec user.Record
		if err := scan(rows, &rec); err != nil {
			return []user.Record{}, store.ErrQuery.Cause(err).Make()
		}

		recs = append(recs, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return recs, nil
}

func (s *SQLDB) CountChiefs(ctx context.Context) (int, error) {
	row := s.db.QueryRowContext(ctx, count_chiefs)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, store.ErrQuery.Cause(err).Make()
	}

	return count, nil
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (user.Record, error) {
	return s.get(ctx, get, uuid.Bytes())
}

func (s *SQLDB) GetBySIAPE(ctx context.Context, siape string) (user.Record, error) {
	return s.get(ctx, get_by_siape, siape)
}

func (s *SQLDB) get(ctx context.Context, query string, args ...any) (user.Record, error) {
	row := s.db.QueryRowContext(ctx, query, args...)

	var ent user.Record
	if err := scan(row, &ent); err != nil {
		if err == sql.ErrNoRows {
			return user.Record{}, user.ErrNotFound
		}

		return user.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return ent, nil
}

func (s *SQLDB) Create(ctx context.Context, rec user.CreateRecord) error {
	_, err := s.db.ExecContext(ctx, create, rec.UUID.Bytes(), rec.SIAPE, rec.Name, rec.Email, rec.Password, rec.Role.String(), rec.Created, rec.Updated)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Patch(ctx context.Context, uuid uuid.UUID, rec user.PatchRecord) error {
	res, err := s.db.ExecContext(ctx, patch, store.NoneNil(rec.Name), store.NoneNil(rec.Email), rec.Updated, uuid.Bytes())
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	changed, err := res.RowsAffected()
	if err == nil && changed == 0 {
		return user.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, delete, uuid.Bytes())
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}
	return nil
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(user.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(&SQLDB{db: tx})
	})
}

func scan(scanner store.Scanner, ent *user.Record) error {
	var bytes []byte
	var string string

	if err := scanner.Scan(
		&bytes,
		&ent.SIAPE,
		&ent.Name,
		&ent.Email,
		&ent.Password,
		&string,
		&ent.Logged,
		&ent.Created,
		&ent.Updated,
	); err != nil {
		return err
	}

	err := problem.Join(
		service.Set(&ent.UUID, bytes, uuid.FromBytes),
		service.Set(&ent.Role, string, role_from_string),
	)
	if err != nil {
		return err
	}

	return nil
}

func role_from_string(role string) (auth.Role, error) {
	if role, ok := auth.FromString(role); ok {
		return role, nil
	}
	return auth.Unlogged, user.ErrRoleInvalid
}

const (
	list         = `select uuid, siape, name, email, password, role, logged, created, updated from Users_View`
	count_chiefs = `select count(*) from Users_View where role = 'chief'`
	get          = `select uuid, siape, name, email, password, role, logged, created, updated from Users_View where uuid = ?`
	get_by_siape = `select uuid, siape, name, email, password, role, logged, created, updated from Users_View where siape = ?`
	create       = `insert into Users (uuid, siape, name, email, password, role, created, updated) values (?, ?, ?, ?, ?, ?, ?, ?)`
	patch        = `update Users set name = coalesce(?, name), email = coalesce(?, email), updated = ? where uuid = ?`
	delete       = `delete from Users where uuid = ?`
)
