package store

import "github.com/alan-b-lima/pkg/opt"

func SomeSet[T any](dst *T, src opt.Opt[T]) {
	val, ok := src.Unwrap()
	if !ok {
		return
	}

	*dst = val
}

func NoneNil[T any](opt opt.Opt[T]) any {
	if val, ok := opt.Unwrap(); ok {
		return val
	}

	return nil
}
