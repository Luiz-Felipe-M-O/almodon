package store

import "github.com/alan-b-lima/pkg/opt"

func SomeSet[T any](dst *T, src opt.Opt[T]) {
	val, ok := src.Unwrap()
	if !ok {
		return
	}

	*dst = val
}
