package store

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrDB = problem.Imp(problem.UnexpectedError, "database-error").Message("unexpected error while accessing database")

	ErrQuery = problem.Imp(problem.UnexpectedError, "query-error").Message("unexpected error while executing query")
	ErrExec  = problem.Imp(problem.UnexpectedError, "exec-error").Message("unexpected error while executing command")

	ErrTx       = problem.Imp(problem.UnexpectedError, "transaction-error").Message("unexpected error while processing transaction")
	ErrNestedTx = problem.New(problem.UnexpectedError, "open-transaction-inside-transaction", "cannot open a transaction inside another transaction", nil, nil)

	ErrNotExtendable = problem.New(problem.UnexpectedError, "not-extendable", "store does not support transaction extension", nil, nil)
	ErrNotTx         = problem.New(problem.UnexpectedError, "not-transaction", "store is not a transaction", nil, nil)
	ErrIllegalJoin   = problem.New(problem.UnexpectedError, "illegal-join", "cannot join transactions from different pools", nil, nil)
)
