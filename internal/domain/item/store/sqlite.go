package itemstore

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/alan-b-lima/almodon/internal/domain/item"
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

var _ item.Store = (*SQLDB)(nil)

func New(db *sql.DB) *SQLDB {
	return &SQLDB{db: db}
}

func (s *SQLDB) List(ctx context.Context) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListByMaterial(ctx context.Context, uuid uuid.UUID) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list_by_material, uuid.Bytes())
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListByECampus(ctx context.Context, ecampus int) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list_by_ecampus, ecampus)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListByCATMAT(ctx context.Context, catmat int) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list_by_catmat, catmat)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) ListBySIADS(ctx context.Context, siads int) ([]item.Record, error) {
	rows, err := s.db.QueryContext(ctx, list_by_siads, siads)
	if err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}
	defer rows.Close()

	return scan_list(rows)
}

func (s *SQLDB) Get(ctx context.Context, uuid uuid.UUID) (item.Record, error) {
	row := s.db.QueryRowContext(ctx, get, uuid.Bytes())

	ent, err := scan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return item.Record{}, item.ErrNotFound
		}

		return item.Record{}, store.ErrQuery.Cause(err).Make()
	}

	return ent, nil
}

func (s *SQLDB) History(ctx context.Context, uuid uuid.UUID) (item.HistoryRecord, error) {
	var hist item.HistoryRecord
	err := s.run_tx(ctx, func(s *SQLDB) error {
		rec, err := s.Get(ctx, uuid)
		if err != nil {
			return err
		}

		hist = item.HistoryRecord{
			UUID:    rec.UUID,
			Version: rec.Version,
			Created: rec.Created,
			Updated: rec.Updated,
		}

		rows, err := s.db.QueryContext(ctx, history, uuid.Bytes())
		if err != nil {
			return store.ErrQuery.Cause(err).Make()
		}
		defer rows.Close()

		recs := make([]item.PastRecord, 0, hist.Version)
		for rows.Next() {
			var rec item.PastRecord
			if err := scan_true(rows, &rec); err != nil {
				return store.ErrQuery.Cause(err).Make()
			}

			recs = append(recs, rec)
		}
		if err := rows.Err(); err != nil {
			return store.ErrQuery.Cause(err).Make()
		}

		hist.Versions = recs
		return nil
	})
	if err != nil {
		return item.HistoryRecord{}, err
	}

	return hist, nil
}

func (s *SQLDB) Create(ctx context.Context, ent item.Entity) error {
	return s.run_tx(ctx, func(s *SQLDB) error {
		version := 1

		if _, err := s.db.ExecContext(ctx, create_true, ent.UUID.Bytes(), version, ent.Material.Bytes(), ent.Amount, int64(ent.UnitCost), ent.Expires, ent.Updated); err != nil {
			return store.ErrQuery.Cause(err).Make()
		}

		if _, err := s.db.ExecContext(ctx, create_surf, ent.UUID.Bytes(), version, ent.Created); err != nil {
			return store.ErrQuery.Cause(err).Make()
		}

		return nil
	})
}

func (s *SQLDB) Patch(ctx context.Context, uuid uuid.UUID, ent item.PatchEntity) error {
	return s.run_tx(ctx, func(s *SQLDB) error {
		rec, err := s.Get(ctx, uuid)
		if err != nil {
			return err
		}

		version := rec.Version + 1

		if _, err := s.db.ExecContext(ctx, create_true,
			uuid.Bytes(),
			version,
			store.Or(ent.Material, rec.Material).Bytes(),
			store.Or(ent.Amount, rec.Amount),
			store.Or(ent.UnitCost, rec.UnitCost),
			store.Or(ent.Expires, rec.Expires),
			ent.Updated,
		); err != nil {
			return store.ErrQuery.Cause(err).Make()
		}

		if _, err := s.db.ExecContext(ctx, update_surf, version, uuid.Bytes()); err != nil {
			return store.ErrQuery.Cause(err).Make()
		}

		return nil
	})
}

func (s *SQLDB) Delete(ctx context.Context, uuid uuid.UUID) error {
	return s.run_tx(ctx, func(s *SQLDB) error {
		if _, err := s.db.ExecContext(ctx, delete_surf, uuid.Bytes()); err != nil {
			return store.ErrQuery.Cause(err).Make()
		}

		if _, err := s.db.ExecContext(ctx, delete_true, uuid.Bytes()); err != nil {
			return store.ErrQuery.Cause(err).Make()
		}

		return nil
	})
}

func (s *SQLDB) RunTx(ctx context.Context, proc func(item.Store) error) error {
	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(&SQLDB{db: tx})
	})
}

func (s *SQLDB) run_tx(ctx context.Context, proc func(*SQLDB) error) error {
	if _, ok := s.db.(*sql.DB); !ok {
		return proc(s)
	}

	return store.WithTx(ctx, s.db, func(tx store.DBTx) error {
		return proc(&SQLDB{db: tx})
	})
}

func scan(scanner store.Scanner) (item.Record, error) {
	var bytes1, bytes2 []byte

	var rec item.Record
	if err := scanner.Scan(
		&bytes1,
		&rec.Version,
		&rec.Name,
		&rec.ECampus,
		&rec.CATMAT,
		&rec.SIADS,
		&bytes2,
		&rec.Amount,
		(*int64)(&rec.UnitCost),
		&rec.Unit,
		&rec.Expires,
		&rec.Min,
		&rec.Created,
		&rec.Updated,
	); err != nil {
		return item.Record{}, err
	}

	if err := problem.Join(
		service.Set(&rec.UUID, bytes1, uuid.FromBytes),
		service.Set(&rec.Material, bytes2, uuid.FromBytes),
	); err != nil {
		return item.Record{}, err
	}

	return rec, nil
}

func scan_list(rows *sql.Rows) ([]item.Record, error) {
	var recs []item.Record
	for rows.Next() {
		ent, err := scan(rows)
		if err != nil {
			return nil, store.ErrQuery.Cause(err).Make()
		}

		recs = append(recs, ent)
	}
	if err := rows.Err(); err != nil {
		return nil, store.ErrQuery.Cause(err).Make()
	}

	return recs, nil
}

func scan_true(scanner store.Scanner, rec *item.PastRecord) error {
	var bytes []byte

	if err := scanner.Scan(
		&rec.Version,
		&bytes,
		&rec.Amount,
		(*int64)(&rec.UnitCost),
		&rec.Expires,
		&rec.Created,
	); err != nil {
		return err
	}

	return service.Set(&rec.Material, bytes, uuid.FromBytes)
}

const (
	list             = `select uuid, version, name, ecampus, catmat, siads, material, amount, unit_cost, unit, expires, min, created, updated from Items_View`
	list_by_material = list + ` where material = ?`
	list_by_ecampus  = list + ` where ecampus = ?`
	list_by_catmat   = list + ` where catmat = ?`
	list_by_siads    = list + ` where siads = ?`

	get     = list + ` where uuid = ?`
	history = `select version, material, amount, unit_cost, expires, created from Items_History_View where uuid = ?`

	create_true = `insert into Items_true (uuid, version, material, amount, unit_cost, expires, created) values (?, ?, ?, ?, ?, ?, ?)`
	create_surf = `insert into Items (uuid, version, created) values (?, ?, ?)`
	update_surf = `update Items set version = ? where uuid = ?`

	delete_true = `delete from Items_true where uuid = ?`
	delete_surf = `delete from Items where uuid = ?`
)
