// Package ordmap implements an ordered map type.
package ordmap

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

// Range calls f sequentially for each key and value present in the map
// following the map ordering: least recently updated first.
// If f returns false, Range stops the iteration.
func (m *Map[K, V]) Range(f func(key K, val V) bool) {
	for e := m.first; e != nil; e = e.next {
		if !f(e.k, e.v) {
			return
		}
	}
}

// Keys returns a slice of all keys in the map in map order.
func (m *Map[K, V]) Keys() []K {
	if len(m.m) == 0 {
		return nil
	}
	keys := make([]K, 0, len(m.m))
	for e := m.first; e != nil; e = e.next {
		keys = append(keys, e.k)
	}
	return keys
}

// Values returns a slice of all values in the map in map order.
func (m *Map[K, V]) Values() []V {
	if len(m.m) == 0 {
		return nil
	}
	vals := make([]V, 0, len(m.m))
	for e := m.first; e != nil; e = e.next {
		vals = append(vals, e.v)
	}
	return vals
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
