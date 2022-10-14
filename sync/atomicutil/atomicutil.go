// Package atomicutil provides Value, a generic version of atomic.Value.
package atomicutil

import "sync/atomic"

// TODO(caleb): Inheriting the prohibition from atomic.Value against storing
// values of different concrete types (of T is an interface type) is kind of
// weird. But always boxing (for example, by using atomic.Pointer[T] underneath
// rather than atomic.Value) isn't great because that forces allocs for types
// that otherwise wouldn't need them (float64, maps, etc).
//
// Figure this out.
//
// Also explain the restriction better in the docs.

// A Value provides an atomic load and store of a value of type T.
// The zero value of a Value returns the zero value of T from Load.
// Once Store has been called, a Value must not be copied.
//
// A Value must not be copied after first use.
type Value[T any] struct {
	v atomic.Value
}

// Load returns the value set by the most recent Store.
// It returns the zero value of T if there has been no call to Store
// for this Value.
func (v *Value[T]) Load() T {
	if x := v.v.Load(); x != nil {
		return x.(T)
	}
	var x T
	return x
}

// Store sets the value of the Value to val.
//
// If T is an interface type, all calls to Store for a given Value must use
// non-nil values of the same concrete type. Store of an inconsistent type panics,
// as does Store(nil).
func (v *Value[T]) Store(val T) {
	v.v.Store(val)
}

// Swap stores new into the Value and returns the previous value.
// It returns the zero value of T if the Value is empty.
//
// If T is an interface type, all calls to Swap for a given Value must use
// non-nil values of the same concrete type. Store of an inconsistent type panics,
// as does Store(nil).
func (v *Value[T]) Swap(new T) (old T) {
	if x := v.v.Swap(new); x != nil {
		return x.(T)
	}
	return old
}
