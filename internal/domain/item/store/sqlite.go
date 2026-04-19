package itemstore

import (
	"context"
	"database/sql"

	"github.com/alan-b-lima/almodon/internal/domain/item"
	"github.com/alan-b-lima/almodon/internal/support/service"
	"github.com/alan-b-lima/almodon/internal/support/store"
	"github.com/alan-b-lima/almodon/pkg/uuid"
	"github.com/alan-b-lima/pkg/problem"
)

const Table = `
create table if not exists Items (
	uuid      blob primary key,
	material  blob not null,
	amount    real not null,
	unit_cost real not null,
	arrival   datetime not null,
	expires   datetime not null,
	created   datetime not null,
	updated   datetime not null,

	foreign key (material) references Materials(uuid)
);`

const Views = `
create view if not exists Items_View as
	select i.uuid, m.name, m.ecampus, m.catmat, m.siads, i.material, i.amount, i.unit_cost, m.unit, i.arrival, i.expires, m.min, i.created, i.updated
	from Items i
	join Materials m on i.material = m.uuid;`

type SQLDB struct {
	db store.DBTx
}

var _ item.Store = (*SQLDB)(nil)

func New(db store.DBTx) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]item.Record, error) {
	return scan_list(s.db.QueryContext(ctx, list))
}

func (s *SQLDB) ListByMaterial(ctx context.Context, material uuid.UUID) ([]item.Record, error) {
	return scan_list(s.db.QueryContext(ctx, list_by_material, material))
}

func (s *SQLDB) ListByECampus(ctx context.Context, ecampus int) ([]item.Record, error) {
	return scan_list(s.db.QueryContext(ctx, list_by_ecampus, ecampus))
}

func (s *SQLDB) ListByCATMAT(ctx context.Context, catmat int) ([]item.Record, error) {
	return scan_list(s.db.QueryContext(ctx, list_by_catmat, catmat))
}

func (s *SQLDB) ListBySIADS(ctx context.Context, siads int) ([]item.Record, error) {
	return scan_list(s.db.QueryContext(ctx, list_by_siads, siads))
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (item.Record, error) {
	row := s.db.QueryRowContext(ctx, get, uuid.Bytes())

	var ent item.Record
	if ok, err := scan(&ent, row); err != nil {
		if ok {
			return item.Record{}, err
		}

		if err == sql.ErrNoRows {
			return item.Record{}, item.ErrNotFound
		}
		return item.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return ent, nil
}

func (s *SQLDB) Create(ctx context.Context, rec item.CreateRecord) error {
	_, err := s.db.ExecContext(ctx, create,
		rec.UUID.Bytes(),
		rec.Material.Bytes(),
		rec.Amount,
		rec.UnitCost,
		rec.Arrival,
		rec.Expires,
		rec.Created,
		rec.Updated,
	)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	return nil
}

func (s *SQLDB) UpdateAmount(ctx context.Context, uuid uuid.UUID, amount float64) error {
	res, err := s.db.ExecContext(ctx, update_amount, amount)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return item.ErrNotFound
	}

	return nil
}

func (s *SQLDB) Patch(ctx context.Context, uuid uuid.UUID, rec item.PatchRecord) error {
	res, err := s.db.ExecContext(ctx, patch,
		store.NoneNil(rec.Material),
		store.NoneNil(rec.UnitCost),
		store.NoneNil(rec.Arrival),
		store.NoneNil(rec.Expires),
		rec.Updated,
		uuid.Bytes(),
	)
	if err != nil {
		return store.ErrExec.Cause(err).Make()
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return item.ErrNotFound
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

func scan_list(rows *sql.Rows, err error) ([]item.Record, error) {
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	var ents []item.Record
	for rows.Next() {
		var ent item.Record
		if ok, err := scan(&ent, rows); err != nil {
			if ok {
				return nil, err
			}

			return nil, store.ErrQuery.Cause(err).Make()
		}

		ents = append(ents, ent)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return ents, nil
}

func scan(ent *item.Record, scanner interface{ Scan(...any) error }) (bool, error) {
	var bytes1, bytes2 []byte

	if err := scanner.Scan(
		&bytes1,
		&ent.Name,
		&ent.ECampus,
		&ent.CATMAT,
		&ent.SIADS,
		&bytes2,
		&ent.Amount,
		&ent.UnitCost,
		&ent.Unit,
		&ent.Arrival,
		&ent.Expires,
		&ent.Min,
		&ent.Created,
		&ent.Updated,
	); err != nil {
		return false, err
	}

	err := problem.Join(
		service.Set(&ent.UUID, bytes1, uuid.FromBytes),
		service.Set(&ent.Material, bytes2, uuid.FromBytes),
	)
	if err != nil {
		return true, store.ErrQuery.Cause(err).Make()
	}

	return true, nil
}

const (
	list             = `select uuid, name, ecampus, catmat, siads, material, amount, unit_cost, unit, arrival, expires, min, created, updated from Items_View`
	list_by_material = list + ` where material = ?`
	list_by_ecampus  = list + ` where ecampus = ?`
	list_by_catmat   = list + ` where catmat = ?`
	list_by_siads    = list + ` where siads = ?`
	get              = list + ` where uuid = ?`

	create        = `insert into Items (uuid, material, amount, unit_cost, arrival, expires, created, updated) values (?, ?, ?, ?, ?, ?, ?, ?)`
	update_amount = `update Items set amount = ? where uuid = ?`
	patch         = `update Items set material = coalesce(?, material), unit_cost = coalesce(?, unit_cost), arrival = coalesce(?, arrival), expires = coalesce(?, expires), updated = ? where uuid = ?`
	delete        = `delete from Items where uuid = ?`
)
