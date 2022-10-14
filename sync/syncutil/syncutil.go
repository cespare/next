// Package syncutil provides Pool and Map, generic versions of
// sync.Pool and sync.Map, respectively.
package syncutil

import "sync"

// A Pool is a set of temporary objects of type T that may be individually saved
// and retrieved.
//
// See sync.Pool for the details of Pool's behavior and caveats about its use.
type Pool[T any] struct {
	// New optionally specifies a function to generate a value when Get
	// would otherwise return the zero value.
	// It may not be changed concurrently with calls to Get.
	New func() T
	p   sync.Pool
}

// Get selects an arbitrary item from the Pool, removes it from the Pool,
// and returns it to the caller. Get may choose to ignore the pool and
// treat it as empty. Callers should not assume any relation between values
// passed to Put and values returned by Get.
//
// If Get would otherwise returns the zero value of T and p.New is non-nil,
// Get returns the result of calling p.New.
func (p *Pool[T]) Get() T {
	if v := p.p.Get(); v != nil {
		return v.(T)
	}
	if p.New != nil {
		return p.New()
	}
	var x T
	return x
}

// Put adds x to the pool.
func (p *Pool[T]) Put(x T) {
	p.p.Put(x)
}

// Map is like a Go map[K]V but is safe for concurrent use by multiple
// goroutines without additional locking or coordination. Loads, stores, and
// deletes run in amortized constant time.
//
// See sync.Map for the details of Map's behavior and caveats about its use.
type Map[K comparable, V any] struct {
	m sync.Map
}

// Delete deletes the value for a key.
func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Load returns the value stored in the map for a key,
// or the zero value of V if no value is present.
// The ok result indicates whether value was found in the map.
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	if v != nil {
		value = v.(V)
	}
	return value, ok
}

// LoadAndDelete deletes the value for a key, returning the previous value if
// any. The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, ok := m.m.LoadAndDelete(key)
	if v != nil {
		value = v.(V)
	}
	return value, ok
}

// LoadOrStore returns the existing value for the key if present. Otherwise,
// it stores and returns the given value. The loaded result is true if the
// value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, ok := m.m.LoadOrStore(key, value)
	return v.(V), ok
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the
// Map's contents: no key will be visited more than once, but if the value
// for any key is stored or deleted concurrently (including by f), Range may
// reflect any mapping for that key from any point during the Range call.
// Range does not block other methods on the receiver; even f itself may call
// any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(k, v any) bool { return f(k.(K), v.(V)) })
}

// Store sets the value for a key.
func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}
