package entity

import "github.com/alan-b-lima/pkg/opt"

func Set[D, S any](dst *D, src S, proc func(S) (D, error)) error {
	val, err := proc(src)
	if err != nil {
		return err
	}

	*dst = val
	return nil
}

func SetOpt[F, R any](dst *opt.Opt[R], src opt.Opt[F], proc func(F) (R, error)) error {
	val, ok := src.Unwrap()
	if !ok {
		return nil
	}

	res, err := proc(val)
	if err != nil {
		return err
	}

	*dst = opt.Some(res)
	return nil
}
