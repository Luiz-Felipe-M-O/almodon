package store

import (
	"context"
	"database/sql"
)

type DBTx interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)

	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func WithTx(ctx context.Context, dbtx DBTx, proc func(DBTx) error) error {
	db, ok := dbtx.(*sql.DB)
	if !ok {
		return ErrNestedTx
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return ErrTx.Cause(err).Make()
	}

	if err := proc(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return ErrTx.Cause(err).Make()
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return ErrTx.Cause(err).Make()
	}
	return nil
}
