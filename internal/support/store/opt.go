package store

import "github.com/alan-b-lima/pkg/opt"

func NoneNil[T any](opt opt.Opt[T]) any {
	if val, ok := opt.Unwrap(); ok {
		return val
	}

	return nil
}

func Or[T any](opt opt.Opt[T], def T) T {
	if val, ok := opt.Unwrap(); ok {
		return val
	}

	return def
}
