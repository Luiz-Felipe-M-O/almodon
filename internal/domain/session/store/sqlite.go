package sessionstore

import (
	"context"
	"database/sql"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/problem"
)

const Table = `
create table if not exists Sessions (
	uuid    blob primary key,
	user    blob not null,
	renewed int not null,
	expires datetime not null,
	created datetime not null,

	foreign key (user) references Users(uuid)
)
`

const Indexes = `
create index if not exists Sessions_user on Sessions(user)
`

type SQLDB struct {
	db store.DBTx
}

var _ session.Store = (*SQLDB)(nil)

func New(db store.DBTx) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (session.Record, error) {
	res := s.db.QueryRowContext(ctx, get, uuid.Bytes())

	var r session.Record

	if ok, err := scan(&r, res); err != nil {
		if ok {
			return session.Record{}, err
		}

		if err == sql.ErrNoRows {
			return session.Record{}, session.ErrNotFound
		}
		return session.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return r, nil
}

func (s *SQLDB) Create(ctx context.Context, rec session.CreateRecord) error {
	_, err := s.db.ExecContext(ctx, create, rec.UUID.Bytes(), rec.User.Bytes(), rec.Renewed, rec.Expires, rec.Created)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Update(ctx context.Context, uuid uuid.UUID, rec session.UpdateRecord) error {
	_, err := s.db.ExecContext(ctx, update, rec.Renewed, rec.Expires, uuid.Bytes())
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

func (s *SQLDB) DeleteExpired(ctx context.Context, deadline time.Time) error {
	_, err := s.db.ExecContext(ctx, delete_expired, deadline)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(session.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(New(tx))
	})
}

func scan(ent *session.Record, scanner interface{ Scan(...any) error }) (bool, error) {
	var bytes1, bytes2 []byte

	if err := scanner.Scan(&bytes1, &bytes2, &ent.Renewed, &ent.Expires, &ent.Created); err != nil {
		return false, err
	}

	err := problem.Join(
		service.Set(&ent.UUID, bytes1, uuid.FromBytes),
		service.Set(&ent.User, bytes2, uuid.FromBytes),
	)
	if err != nil {
		return true, err
	}

	return true, nil
}

const (
	get            = `select uuid, user, renewed, expires, created from Sessions where uuid = ?`
	create         = `insert into Sessions (uuid, user, renewed, expires, created) values (?, ?, ?, ?, ?)`
	update         = `update Sessions set renewed = ?, expires = ? where uuid = ?`
	delete         = `delete from Sessions where uuid = ?`
	delete_expired = `delete from Sessions where expires < ?`
)
