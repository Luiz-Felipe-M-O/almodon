package userstore

import (
	"context"
	"database/sql"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/opt"
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
);
`

const Indexes = `
create unique index if not exists Users_siape on Users(siape);
`

const Views = `
create view if not exists Users_View as
	select u.uuid, u.siape, u.name, u.email, u.password, iif(p.uuid is null, u.role, "promoted-admin") as "role", s.uuid is not null as "logged", u.created, u.updated
	from Users u
	left join Sessions s on s.user = u.uuid
	left join Promotions p on p.user = u.uuid;
`

type SQLDB struct {
	db store.DBTx
}

var _ user.Store = (*SQLDB)(nil)

func New(db store.DBTx) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]user.Record, error) {
	rows, err := s.db.QueryContext(ctx, list)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	var ents []user.Record
	for rows.Next() {
		var ent user.Record
		if ok, err := scan(&ent, rows); err != nil {
			if ok {
				return []user.Record{}, err
			}

			return []user.Record{}, store.ErrQuery.Cause(err).Make()
		}

		ents = append(ents, ent)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return ents, nil
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
	row := s.db.QueryRowContext(ctx, get, uuid.Bytes())

	var ent user.Record
	if ok, err := scan(&ent, row); err != nil {
		if ok {
			return user.Record{}, err
		}

		if err == sql.ErrNoRows {
			return user.Record{}, user.ErrNotFound
		}
		return user.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return ent, nil
}

func (s *SQLDB) GetBySIAPE(ctx context.Context, siape string) (user.Record, error) {
	row := s.db.QueryRowContext(ctx, get_by_siape, siape)

	var ent user.Record
	if ok, err := scan(&ent, row); err != nil {
		if ok {
			return user.Record{}, err
		}

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
	_, err := s.db.ExecContext(ctx, patch, none_nil(rec.Name), none_nil(rec.Email), rec.Updated, uuid.Bytes())
	if err != nil {
		return store.ErrExec.Cause(err).Make()
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

func scan(ent *user.Record, scanner interface{ Scan(...any) error }) (bool, error) {
	var bytes []byte
	var string string

	if err := scanner.Scan(&bytes, &ent.SIAPE, &ent.Name, &ent.Email, &ent.Password, &string, &ent.Logged, &ent.Created, &ent.Updated); err != nil {
		return false, err
	}

	err := problem.Join(
		service.Set(&ent.UUID, bytes, uuid.FromBytes),
		service.Set(&ent.Role, string, role_from_string),
	)
	if err != nil {
		return true, store.ErrQuery.Cause(err).Make()
	}

	return true, nil
}

func role_from_string(role string) (auth.Role, error) {
	if role, ok := auth.FromString(role); ok {
		return role, nil
	}
	return auth.Unlogged, user.ErrRoleInvalid
}

func none_nil[T any](opt opt.Opt[T]) any {
	if val, ok := opt.Unwrap(); ok {
		return val
	}
	return nil
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
