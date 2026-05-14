package sessionstore

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/session"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/problem"
)

//go:embed sqlite.sql
var Script string

type SQLDB struct {
	db store.DBTx
}

var _ session.Store = (*SQLDB)(nil)

func New(db store.DBTx) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]session.Record, error) {
	rows, err := s.db.QueryContext(ctx, list)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	var recs []session.Record
	for rows.Next() {
		var rec session.Record
		if err := scan(&rec, rows); err != nil {
			return nil, store.ErrQuery.Cause(err).Make()
		}

		recs = append(recs, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return recs, nil
}

func (s *SQLDB) Get(ctx context.Context, token session.Token) (session.Record, error) {
	return s.get(s.db.QueryRowContext(ctx, get, token.Bytes()))
}

func (s *SQLDB) GetByUser(ctx context.Context, user uuid.UUID) (session.Record, error) {
	return s.get(s.db.QueryRowContext(ctx, get_by_user, user.Bytes()))
}

func (s *SQLDB) get(row *sql.Row) (session.Record, error) {
	var rec session.Record
	if err := scan(&rec, row); err != nil {
		if err == sql.ErrNoRows {
			return session.Record{}, session.ErrNotFound
		}

		return session.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return rec, nil
}

func (s *SQLDB) Create(ctx context.Context, rec session.CreateRecord) error {
	_, err := s.db.ExecContext(ctx, create, rec.Token.Bytes(), rec.User.Bytes(), rec.Renewed, rec.Expires, rec.Created)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Update(ctx context.Context, token session.Token, rec session.UpdateRecord) error {
	res, err := s.db.ExecContext(ctx, update, rec.Renewed, rec.Expires, token.Bytes())
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	changed, err := res.RowsAffected()
	if err == nil && changed == 0 {
		return session.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Delete(ctx context.Context, token session.Token) error {
	_, err := s.db.ExecContext(ctx, delete, token.Bytes())
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

func scan(ent *session.Record, scanner store.Scanner) error {
	var bytes1, bytes2 []byte

	err := scanner.Scan(
		&bytes1,
		&bytes2,
		&ent.Renewed,
		&ent.Expires,
		&ent.Created,
	)
	if err != nil {
		return err
	}

	return problem.Join(
		service.Set(&ent.Token, bytes1, token_from_bytes),
		service.Set(&ent.User, bytes2, uuid.FromBytes),
	)
}

func token_from_bytes(bytes []byte) (session.Token, error) {
	if len(bytes) != session.TokenLen {
		return session.Token{}, session.ErrInvalidToken
	}

	return session.Token(bytes), nil
}

const (
	list        = `select token, user, renewed, expires, created from Sessions`
	get         = list + ` where token = ?`
	get_by_user = list + ` where user = ?`

	create         = `insert into Sessions (token, user, renewed, expires, created) values (?, ?, ?, ?, ?)`
	update         = `update Sessions set renewed = ?, expires = ? where token = ?`
	delete         = `delete from Sessions where token = ?`
	delete_expired = `delete from Sessions where expires < ?`
)
