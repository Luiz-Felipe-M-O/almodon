package promotionstore

import (
	"context"
	"database/sql"

	"github.com/alan-b-lima/almodon/internal/domain/promotion"
	"github.com/alan-b-lima/almodon/internal/support/entity"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/problem"
)

const Table = `
create table if not exists Promotions (
	uuid    blob primary key,
	user    blob not null,
	expires datetime not null,

	foreign key (user) references Users(uuid)
);
`

const Indexes = `
create index if not exists Promotions_user on Promotions(user)
`

type SQLDB struct {
	db store.DBTx
}

func New(db store.DBTx) promotion.Store {
	return &SQLDB{db: db}
}

var _ promotion.Store = (*SQLDB)(nil)

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (promotion.Record, error) {
	row := s.db.QueryRowContext(ctx, get, uuid.Bytes())

	var res promotion.Record
	if err := scan(&res, row); err != nil {
		if err == sql.ErrNoRows {
			return promotion.Record{}, promotion.ErrNotFound
		}

		return promotion.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return res, nil
}

func (s *SQLDB) GetByUser(ctx context.Context, uuid uuid.UUID) (promotion.Record, error) {
	row := s.db.QueryRowContext(ctx, get_by_user, uuid.Bytes())

	var res promotion.Record
	if err := scan(&res, row); err != nil {
		if err == sql.ErrNoRows {
			return promotion.Record{}, promotion.ErrNotFound
		}

		return promotion.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return res, nil
}

func (s *SQLDB) Create(ctx context.Context, req promotion.CreateRecord) error {
	_, err := s.db.ExecContext(ctx, create, req.UUID.Bytes(), req.User, req.Expires)
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

func (s *SQLDB) RunTx(ctx context.Context, proc func(promotion.Store) error) error {
	err := store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(New(tx))
	})
	return err
}

func scan(ent *promotion.Record, scanner interface{ Scan(...any) error }) error {
	var bytes1, bytes2 []byte

	if err := scanner.Scan(&bytes1, &bytes2, &ent.Expires); err != nil {
		return err
	}

	err := problem.Join(
		entity.Set(&ent.UUID, bytes1, uuid.FromBytes),
		entity.Set(&ent.User, bytes2, uuid.FromBytes),
	)
	if err != nil {
		return store.ErrQuery.Cause(err).Make()
	}

	return nil
}

const (
	get         = `select uuid, user, expires from Promotions where uuid = ?`
	get_by_user = `select uuid, user, expires from Promotions where user = ?`
	create      = `insert into Promotions (uuid, user, expires) values (?, ?, ?)`
	update      = `update Promotions set expires = ? where uuid = ?`
	delete      = `delete from Promotions where uuid = ?`
)
