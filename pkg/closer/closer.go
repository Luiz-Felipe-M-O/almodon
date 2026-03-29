// Package closer implements functionality for managing and closing resources.
package closer

import "errors"

// Closer is an interface that defines a method for closing resources.
type Closer interface {
	Close() error
}

// CloserFunc is a function type that implements the Closer interface.
type CloserFunc func() error

// Close implements the Closer interface for CloserFunc.
func (f CloserFunc) Close() error { return f() }

// Bundle is a struct that manages a collection of Closers and provides methods
// to add, reset, and close them.
type Bundle struct {
	cleanup []Closer
}

// Bundle adds a Closer to the Bundle if the provided value implements the
// Closer interface. It returns true if the value was added.
func (b *Bundle) Bundle(a any) bool {
	if closer, ok := a.(Closer); ok {
		b.cleanup = append(b.cleanup, closer)
		return true
	}

	return false
}

// BundleFunc adds a function as a Closer to the Bundle.
func (b *Bundle) BundleFunc(f func() error) {
	b.Bundle(CloserFunc(f))
}

// BundleMany adds multiple values as Closers to the Bundle if they implement
// the Closer interface.
func (b *Bundle) BundleMany(a ...any) {
	for _, v := range a {
		b.Bundle(v)
	}
}

// Reset clears all Closers from the Bundle without closing them.
func (b *Bundle) Reset() {
	clear(b.cleanup)
	b.cleanup = b.cleanup[:0]
}

// Close calls the Close method on all Closers in the Bundle and returns any
// errors that occur.
// 
// The bundle will conclude all closers even if some of them return an error.
// The errors will be collected and returned as a single error using
// [errors.Join].
func (b *Bundle) Close() error {
	errs := make([]error, 0, len(b.cleanup))

	for _, closer := range b.cleanup {
		errs = append(errs, closer.Close())
	}

	return errors.Join(errs...)
}

// CloseMany is a helper function that takes multiple values and calls the Close
// method on those that implement the Closer interface. It returns any errors
// that occur.
//
// Similar to [Bundle.Close], this function will attempt to close all provided
// values that implement the Closer interface, even if some of them return an
// error. The errors will be collected and returned as a single error using
// [errors.Join].
func CloseMany(a ...any) error {
	errs := make([]error, 0, len(a))

	for _, v := range a {
		if closer, ok := v.(Closer); ok {
			errs = append(errs, closer.Close())
		}
	}

	return errors.Join(errs...)
}
