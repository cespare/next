// Package set defines a Set type that holds a set of elements.
//
// TODO(caleb): NaNs?
package set

// TODO(caleb): We probably need to get rid of the dependencies.

import (
	"fmt"
	"sort"
	"strings"
)

// A Set is a set of elements of some comparable type.
// Sets are implemented using maps, and have similar performance characteristics.
// Like maps, Sets are reference types.
// That is, for Sets s1 = s2 will leave s1 and s2 pointing to the same set of elements:
// changes to s1 will be reflected in s2 and vice-versa.
// Unlike maps, the zero value of a Set is usable; there is no equivalent to make.
// As with maps, concurrent calls to functions and methods that read values are fine;
// concurrent calls to functions and methods that write values are racy.
type Set[E comparable] struct {
	m map[E]struct{}
}

// Of returns a new set containing the listed elements.
func Of[E comparable](v ...E) *Set[E] {
	if len(v) == 0 {
		return &Set[E]{nil}
	}
	m := make(map[E]struct{})
	for _, vv := range v {
		m[vv] = struct{}{}
	}
	return &Set[E]{m}
}

func (s *Set[E]) String() string {
	// Print it out in some deterministic order for now.
	// Better would be to sort numbers numerically.
	vals := make([]string, 0, s.Len())
	for v := range s.m {
		vals = append(vals, fmt.Sprint(v))
	}
	sort.Strings(vals)
	return fmt.Sprintf("set[%s]", strings.Join(vals, " "))
}

func (s *Set[E]) GoString() string {
	var v E
	typeName := fmt.Sprintf("%T", v)
	vals := make([]string, 0, s.Len())
	for v := range s.m {
		vals = append(vals, fmt.Sprintf("%#v", v))
	}
	sort.Strings(vals)
	// TODO(caleb): Technically this is slightly misleading in the case that
	// s.m != nil && len(s.m) == 0 because the result does not yield the
	// same exact thing when interpreted literally as Go code.
	return fmt.Sprintf("set.Of[%s](%s)", typeName, strings.Join(vals, ", "))
}

func (s *Set[E]) init() {
	if s.m == nil {
		s.m = make(map[E]struct{})
	}
}

// Add adds elements to a set.
func (s *Set[E]) Add(v ...E) {
	if len(v) == 0 {
		return
	}
	s.init()
	for _, vv := range v {
		s.m[vv] = struct{}{}
	}
}

// AddSet adds the elements of set s2 to s.
func (s *Set[E]) AddSet(s2 *Set[E]) {
	if len(s2.m) == 0 {
		return
	}
	s.init()
	for v2 := range s2.m {
		s.m[v2] = struct{}{}
	}
}

// Remove removes elements from a set.
// Elements that are not present are ignored.
func (s *Set[E]) Remove(v ...E) {
	for _, vv := range v {
		delete(s.m, vv)
	}
}

// RemoveSet removes the elements of set s2 from s.
// Elements present in s2 but not s are ignored.
func (s *Set[E]) RemoveSet(s2 *Set[E]) {
	for v2 := range s2.m {
		delete(s.m, v2)
	}
}

// Contains reports whether v is in the set.
func (s *Set[E]) Contains(v E) bool {
	_, ok := s.m[v]
	return ok
}

// ContainsAny reports whether any of the elements in s2 are in s.
func (s *Set[E]) ContainsAny(s2 *Set[E]) bool {
	for v2 := range s2.m {
		if _, ok := s.m[v2]; ok {
			return true
		}
	}
	return false
}

// ContainsAll reports whether all of the elements in s2 are in s.
func (s *Set[E]) ContainsAll(s2 *Set[E]) bool {
	for v2 := range s2.m {
		if _, ok := s.m[v2]; !ok {
			return false
		}
	}
	return true
}

// Slice returns the elements in the set s as a slice.
// The values will be in an indeterminate order.
func (s *Set[E]) Slice() []E {
	if len(s.m) == 0 {
		return nil
	}
	vals := make([]E, 0, len(s.m))
	for v := range s.m {
		vals = append(vals, v)
	}
	return vals
}

// Equal reports whether s and s2 contain the same elements.
func (s *Set[E]) Equal(s2 *Set[E]) bool {
	if len(s.m) != len(s2.m) {
		return false
	}
	for v := range s.m {
		if _, ok := s2.m[v]; !ok {
			return false
		}
	}
	return true
}

// Clear removes all elements from s, leaving it empty.
func (s *Set[E]) Clear() {
	for v := range s.m {
		delete(s.m, v)
	}
}

// Clone returns a copy of s.
// The elements are copied using assignment,
// so this is a shallow clone.
func (s *Set[E]) Clone() *Set[E] {
	if len(s.m) == 0 {
		return &Set[E]{nil}
	}
	m := make(map[E]struct{}, len(s.m))
	for v := range s.m {
		m[v] = struct{}{}
	}
	return &Set[E]{m}
}

// RemoveIf deletes any elements from s for which remove returns true.
func (s *Set[E]) RemoveIf(remove func(E) bool) {
	for v := range s.m {
		if remove(v) {
			delete(s.m, v)
		}
	}
}

// Len returns the number of elements in s.
func (s *Set[E]) Len() int {
	return len(s.m)
}

// Do calls f on every element in the set s,
// stopping if f returns false.
// f should not change s.
// f will be called on values in an indeterminate order.
func (s *Set[E]) Do(f func(E) bool) {
	for v := range s.m {
		if !f(v) {
			return
		}
	}
}

// Union constructs a new set containing the union of s1 and s2.
func Union[E comparable](s1, s2 *Set[E]) *Set[E] {
	// TODO(caleb): Presize?
	s := s1.Clone()
	s.AddSet(s2)
	return s
}

// Intersection constructs a new set containing the intersection of s1 and s2.
func Intersection[E comparable](s1, s2 *Set[E]) *Set[E] {
	// TODO(caleb): Presize?
	var s Set[E]
	for v := range s1.m {
		if _, ok := s2.m[v]; ok {
			s.Add(v)
		}
	}
	return &s
}

// Difference constructs a new set containing the elements of s1 that
// are not present in s2.
func Difference[E comparable](s1, s2 *Set[E]) *Set[E] {
	// TODO(caleb): Presize?
	var s Set[E]
	for v := range s1.m {
		if _, ok := s2.m[v]; !ok {
			s.Add(v)
		}
	}
	return &s
}
