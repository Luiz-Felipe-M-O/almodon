package promotionstore

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/promotion"
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

func New(db store.DBTx) promotion.Store {
	return &SQLDB{db: db}
}

var _ promotion.Store = (*SQLDB)(nil)

func (s *SQLDB) List(ctx context.Context) ([]promotion.Record, error) {
	rows, err := s.db.QueryContext(ctx, list)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	var recs []promotion.Record
	for rows.Next() {
		var rec promotion.Record
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

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (promotion.Record, error) {
	return s.get(s.db.QueryRowContext(ctx, get, uuid.Bytes()))
}

func (s *SQLDB) GetByUser(ctx context.Context, uuid uuid.UUID) (promotion.Record, error) {
	return s.get(s.db.QueryRowContext(ctx, get_by_user, uuid.Bytes()))
}

func (s *SQLDB) get(row *sql.Row) (promotion.Record, error) {
	var rec promotion.Record
	if err := scan(&rec, row); err != nil {
		if err == sql.ErrNoRows {
			return promotion.Record{}, promotion.ErrNotFound
		}

		return promotion.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return rec, nil
}

func (s *SQLDB) Create(ctx context.Context, req promotion.CreateRecord) error {
	_, err := s.db.ExecContext(ctx, create, req.UUID.Bytes(), req.User.Bytes(), req.Expires)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) Update(ctx context.Context, uuid uuid.UUID, req promotion.UpdateRecord) error {
	_, err := s.db.ExecContext(ctx, update, req.Expires, uuid.Bytes())
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

func (s *SQLDB) RunTx(ctx context.Context, proc func(promotion.Store) error) error {
	err := store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(New(tx))
	})
	return err
}

func scan(ent *promotion.Record, scanner store.Scanner) error {
	var bytes1, bytes2 []byte

	if err := scanner.Scan(&bytes1, &bytes2, &ent.Expires); err != nil {
		return err
	}

	return problem.Join(
		service.Set(&ent.UUID, bytes1, uuid.FromBytes),
		service.Set(&ent.User, bytes2, uuid.FromBytes),
	)
}

const (
	list        = `select uuid, user, expires from Promotions`
	get         = list + ` where uuid = ?`
	get_by_user = list + ` where user = ?`

	create         = `insert into Promotions (uuid, user, expires) values (?, ?, ?)`
	update         = `update Promotions set expires = ? where uuid = ?`
	delete         = `delete from Promotions where uuid = ?`
	delete_expired = `delete from Promotions where expires <= ?`
)
