// Package ordmap implements an ordered map type.
package ordmap

import "iter"

// TODO(caleb): This list-based approach looks pretty cache-inefficient.
// Add some benchmarks; optimize.

// Map is like a Go map[K]V but is ordered: it retains the insertion/update
// ordering where less recently updated elements precede more recently updated
// elements.
type Map[K comparable, V any] struct {
	m     map[K]*element[K, V]
	first *element[K, V]
	last  *element[K, V]
}

type element[K comparable, V any] struct {
	k    K
	v    V
	prev *element[K, V]
	next *element[K, V]
}

// Get returns the value stored in the map for a key,
// or the zero value of V if no value is present.
// The ok result indicates whether the key was found in the map.
func (m *Map[K, V]) Get(key K) (val V, ok bool) {
	if e, ok := m.m[key]; ok {
		return e.v, true
	}
	return val, false
}

// Set sets the value for a key.
func (m *Map[K, V]) Set(key K, v V) {
	if e, ok := m.m[key]; ok {
		e.v = v
		m.listMoveToEnd(e)
		return
	}
	if m.m == nil {
		m.m = make(map[K]*element[K, V])
	}
	e := &element[K, V]{k: key, v: v}
	m.listAppend(e)
	m.m[key] = e
}

// Delete deletes the value for a key.
func (m *Map[K, V]) Delete(key K) {
	e, ok := m.m[key]
	if !ok {
		return
	}
	m.listDelete(e)
	delete(m.m, key)
}

// All returns an iterator over key-value pairs in the map.
// The iteration order follows the map ordering: least recently updated first.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.first; e != nil; e = e.next {
			if !yield(e.k, e.v) {
				return
			}
		}
	}
}

// Keys returns an iterator over keys in the map.
// The iteration order follows the map ordering: least recently updated first.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for e := m.first; e != nil; e = e.next {
			if !yield(e.k) {
				return
			}
		}
	}
}

// Values returns an iterator over values in the map.
// The iteration order follows the map ordering: least recently updated first.
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for e := m.first; e != nil; e = e.next {
			if !yield(e.v) {
				return
			}
		}
	}
}

func (m *Map[K, V]) listAppend(e *element[K, V]) {
	if m.first == nil {
		m.first = e
	} else {
		m.last.next = e
	}
	e.prev = m.last
	e.next = nil
	m.last = e
}

func (m *Map[K, V]) listMoveToEnd(e *element[K, V]) {
	if m.last == e {
		return
	}
	prev, next := e.prev, e.next
	if prev == nil {
		m.first = next
	} else {
		prev.next = next
	}
	next.prev = prev
	m.last.next = e
	e.prev = m.last
	e.next = nil
	m.last = e
}

func (m *Map[K, V]) listDelete(e *element[K, V]) {
	if e.prev == nil {
		m.first = e.next
	} else {
		e.prev.next = e.next
	}
	if e.next == nil {
		m.last = e.prev
	} else {
		e.next.prev = e.prev
	}
	e.prev = nil
	e.next = nil
}
