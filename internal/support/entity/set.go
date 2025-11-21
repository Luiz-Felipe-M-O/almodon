package entity

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/opt"
)

func Set[D, S any](dst *D, src S, proc func(S) (D, error)) error {
	val, err := proc(src)
	if err != nil {
		return err
	}

	*dst = val
	return nil
}

func SetWithUpdate[D, S any](dst *D, src S, proc func(S) (D, error), update *time.Time) error {
	if err := Set(dst, src, proc); err != nil {
		return err
	}

	if update != nil {
		*update = time.Now()
	}

	return nil
}

func SomeThen[F, R any](dst *opt.Opt[R], src opt.Opt[F], fn func(F) (R, error)) error {
	val, ok := src.Unwrap()
	if !ok {
		return nil
	}

	res, err := fn(val)
	if err != nil {
		return err
	}

	*dst = opt.Some(res)
	return nil
}
