package closer

import "errors"

type Closer interface {
	Close() error
}

type CloserFunc func() error

func (f CloserFunc) Close() error { return f() }

type Bundle struct {
	cleanup []Closer
}

func (b *Bundle) Bundle(a any) bool {
	if closer, ok := a.(Closer); ok {
		b.cleanup = append(b.cleanup, closer)
		return true
	}

	return false
}

func (b *Bundle) BundleFunc(f func() error) {
	b.Bundle(CloserFunc(f))
}

func (b *Bundle) BundleMany(a ...any) {
	for _, v := range a {
		b.Bundle(v)
	}
}

func (b *Bundle) Reset() {
	clear(b.cleanup)
	b.cleanup = b.cleanup[:0]
}

func (b *Bundle) Close() error {
	errs := make([]error, 0, len(b.cleanup))

	for _, closer := range b.cleanup {
		errs = append(errs, closer.Close())
	}

	return errors.Join(errs...)
}

func CloseMany(a ...any) error {
	errs := make([]error, 0, len(a))

	for _, v := range a {
		if closer, ok := v.(Closer); ok {
			errs = append(errs, closer.Close())
		}
	}

	return errors.Join(errs...)
}
