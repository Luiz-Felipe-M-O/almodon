package store

import "github.com/alan-b-lima/pkg/problem"

var (
	ErrDB = problem.Imp(problem.UnexpectedError, "database-error").Message("unknown error while accessing database")

	ErrQuery = problem.Imp(problem.UnexpectedError, "query-error").Message("unknown error while executing query")
	ErrExec  = problem.Imp(problem.UnexpectedError, "exec-error").Message("unknown error while executing command")

	ErrTx       = problem.Imp(problem.UnexpectedError, "transaction-error").Message("unknown error while processing transaction")
	ErrNestedTx = problem.New(problem.UnexpectedError, "open-transaction-inside-transaction", "cannot open a transaction inside another transaction", nil, nil)
)
