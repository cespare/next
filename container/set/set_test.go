package set

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestOf(t *testing.T) {
	for _, tt := range []struct {
		vs []string
	}{
		{nil},
		{[]string{"a"}},
		{[]string{"a", "b", "c"}},
	} {
		got := Of(tt.vs...)
		check(t, got, tt.vs)
	}
}

func TestString(t *testing.T) {
	for _, tt := range []struct {
		set  *Set[string]
		want string
	}{
		{Of[string](), "set[]"},
		{emptyOf[string](), "set[]"},
		{Of("a"), "set[a]"},
		{Of("a", "b", "c"), "set[a b c]"},
	} {
		got := tt.set.String()
		if got != tt.want {
			t.Errorf(
				"&Set{%#v}.String(): got %q; want %q",
				tt.set.m, got, tt.want,
			)
		}
	}
}

func TestGoString(t *testing.T) {
	for _, tt := range []struct {
		set  *Set[string]
		want string
	}{
		{Of[string](), "set.Of[string]()"},
		{emptyOf[string](), "set.Of[string]()"},
		{Of("a"), `set.Of[string]("a")`},
		{Of("a", "b", "c"), `set.Of[string]("a", "b", "c")`},
	} {
		got := tt.set.GoString()
		if got != tt.want {
			t.Errorf(
				"&Set{%#v}.String(): got %q; want %q",
				tt.set.m, got, tt.want,
			)
		}
	}
}

// emptyOf returns a set with a non-nil, empty map.
// We don't provide a direct constructor for this state but we want to cover it
// in tests.
func emptyOf[E comparable]() *Set[E] {
	return &Set[E]{make(map[E]struct{})}
}

func TestAdd(t *testing.T) {
	var s Set[int]
	check(t, &s, nil)
	s.Add()
	check(t, &s, nil)
	s.Add(3)
	check(t, &s, []int{3})
	s.Add(3)
	check(t, &s, []int{3})
	s.Add(1, 3, 5)
	check(t, &s, []int{1, 3, 5})
	s.Add(2, 3, 4)
	check(t, &s, []int{1, 2, 3, 4, 5})
}

func TestAddSet(t *testing.T) {
	var s Set[int]
	check(t, &s, nil)
	s.AddSet(Of[int]())
	check(t, &s, nil)
	s.AddSet(Of(3))
	check(t, &s, []int{3})
	s.AddSet(Of(1, 3, 5))
	check(t, &s, []int{1, 3, 5})
	s.AddSet(Of(2, 3, 4))
	check(t, &s, []int{1, 2, 3, 4, 5})
}

func TestRemove(t *testing.T) {
	var s Set[int]
	check(t, &s, nil)
	s.Remove(3)
	check(t, &s, nil)

	s.Add(1, 2, 3, 4, 5)
	s.Remove(3)
	check(t, &s, []int{1, 2, 4, 5})
	s.Remove(2, 3, 4)
	check(t, &s, []int{1, 5})
	s.Remove(1, 3, 5)
	check(t, &s, []int{})
}

func TestRemoveSet(t *testing.T) {
	var s Set[int]
	check(t, &s, nil)
	s.RemoveSet(Of(3))
	check(t, &s, nil)

	s.Add(1, 2, 3, 4, 5)
	s.RemoveSet(Of(3))
	check(t, &s, []int{1, 2, 4, 5})
	s.RemoveSet(Of(2, 3, 4))
	check(t, &s, []int{1, 5})
	s.RemoveSet(Of(1, 3, 5))
	check(t, &s, []int{})
}

func TestContains(t *testing.T) {
	for _, tt := range []struct {
		set  *Set[int]
		v    int
		want bool
	}{
		{Of[int](), 3, false},
		{emptyOf[int](), 3, false},
		{Of(3), 3, true},
		{Of(4), 3, false},
		{Of(1, 3, 5), 3, true},
		{Of(1, 3, 5), 2, false},
	} {
		got := tt.set.Contains(tt.v)
		if got != tt.want {
			t.Errorf("%s.Contains(%d): got %t", tt.set.debug(), tt.v, got)
		}
	}
}

func TestContainsAny(t *testing.T) {
	for _, tt := range []struct {
		s1   *Set[int]
		s2   *Set[int]
		want bool
	}{
		{Of[int](), Of[int](), false},
		{Of[int](), Of(3), false},
		{Of[int](), emptyOf[int](), false},
		{emptyOf[int](), Of[int](), false},
		{emptyOf[int](), Of(3), false},
		{emptyOf[int](), emptyOf[int](), false},
		{Of(3), Of[int](), false},
		{Of(3), emptyOf[int](), false},
		{Of(3), Of(3), true},
		{Of(3), Of(2), false},
		{Of(3), Of(2, 3), true},
		{Of(1, 3, 5), Of(3), true},
		{Of(1, 3, 5), Of(2, 3), true},
		{Of(1, 3, 5), Of(2, 4), false},
	} {
		got := tt.s1.ContainsAny(tt.s2)
		if got != tt.want {
			t.Errorf(
				"%s.ContainsAny(%s): got %t",
				tt.s1.debug(), tt.s2.debug(), got,
			)
		}
	}
}

func TestContainsAll(t *testing.T) {
	for _, tt := range []struct {
		s1   *Set[int]
		s2   *Set[int]
		want bool
	}{
		{Of[int](), Of[int](), true},
		{Of[int](), Of(3), false},
		{Of[int](), emptyOf[int](), true},
		{emptyOf[int](), Of[int](), true},
		{emptyOf[int](), Of(3), false},
		{emptyOf[int](), emptyOf[int](), true},
		{Of(3), Of[int](), true},
		{Of(3), emptyOf[int](), true},
		{Of(3), Of(3), true},
		{Of(3), Of(2), false},
		{Of(3), Of(2, 3), false},
		{Of(1, 3, 5), Of(3), true},
		{Of(1, 3, 5), Of(1, 5), true},
		{Of(1, 3, 5), Of(2, 3), false},
		{Of(1, 3, 5), Of(2, 4), false},
	} {
		got := tt.s1.ContainsAll(tt.s2)
		if got != tt.want {
			t.Errorf(
				"%s.ContainsAll(%s): got %t",
				tt.s1.debug(), tt.s2.debug(), got,
			)
		}
	}
}

func TestSlice(t *testing.T) {
	for _, tt := range []struct {
		s    *Set[int]
		want []int
	}{
		{Of[int](), []int(nil)},
		{emptyOf[int](), []int(nil)},
		{Of(3), []int{3}},
		{Of(1, 3, 5), []int{1, 3, 5}},
	} {
		got := tt.s.Slice()
		sort.Ints(got)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s.Slice(): got %v; want %v", tt.s.debug(), got, tt.want)
		}
	}
}

func TestEqual(t *testing.T) {
	for _, tt := range []struct {
		s1   *Set[int]
		s2   *Set[int]
		want bool
	}{
		{Of[int](), Of[int](), true},
		{Of[int](), emptyOf[int](), true},
		{Of[int](), Of(3), false},
		{emptyOf[int](), Of[int](), true},
		{emptyOf[int](), emptyOf[int](), true},
		{emptyOf[int](), Of(3), false},
		{Of(3), Of[int](), false},
		{Of(3), emptyOf[int](), false},
		{Of(3), Of(3), true},
		{Of(3), Of(4), false},
		{Of(3), Of(3, 4), false},
		{Of(1, 3), Of(3), false},
		{Of(1, 3), Of(1, 3), true},
		{Of(1, 3), Of(1, 5), false},
	} {
		got := tt.s1.Equal(tt.s2)
		if got != tt.want {
			t.Errorf("%s.Equal(%s): got %t", tt.s1.debug(), tt.s2.debug(), got)
		}
	}
}

func TestClear(t *testing.T) {
	s := Of[string]()
	s.Clear()
	check(t, s, nil)

	s = emptyOf[string]()
	s.Clear()
	check(t, s, []string{})

	s = Of("a", "b", "c")
	s.Clear()
	check(t, s, []string{})
}

func BenchmarkClear(b *testing.B) {
	s := Of[int]()
	for i := 0; i < b.N; i++ {
		s.Add(1, 2, 3, 4, 5)
		s.Clear()
	}
}

func TestClone(t *testing.T) {
	for _, tt := range []struct {
		s    *Set[int]
		want *Set[int]
	}{
		{Of[int](), Of[int]()},
		{emptyOf[int](), Of[int]()},
		{Of(3), Of(3)},
		{Of(1, 3, 5), Of(1, 3, 5)},
	} {
		got := tt.s.Clone()
		if !equal(got, tt.want) {
			t.Errorf(
				"%s.Clone(): got %s; want %s",
				tt.s.debug(), got.debug(), tt.want.debug(),
			)
		}
	}

	s1 := Of(1, 3, 5)
	s2 := s1.Clone()
	s2.Remove(3)
	s2.Add(4)
	check(t, s1, []int{1, 3, 5})
	check(t, s2, []int{1, 4, 5})
}

func TestRemoveIf(t *testing.T) {
	removeAll := func(int) bool { return true }
	removeNone := func(int) bool { return false }
	removeOdd := func(n int) bool { return n%2 != 0 }
	for _, tt := range []struct {
		s      *Set[int]
		remove func(int) bool
		want   []int
	}{
		{Of[int](), removeAll, nil},
		{emptyOf[int](), removeAll, []int{}},
		{Of(3), removeAll, []int{}},
		{Of(3), removeNone, []int{3}},
		{Of(3), removeOdd, []int{}},
		{Of(1, 2, 3, 4, 5), removeAll, []int{}},
		{Of(1, 2, 3, 4, 5), removeNone, []int{1, 2, 3, 4, 5}},
		{Of(1, 2, 3, 4, 5), removeOdd, []int{2, 4}},
	} {
		tt.s.RemoveIf(tt.remove)
		check(t, tt.s, tt.want)
	}
}

func TestLen(t *testing.T) {
	for _, tt := range []struct {
		s    *Set[int]
		want int
	}{
		{Of[int](), 0},
		{emptyOf[int](), 0},
		{Of(3), 1},
		{Of(1, 2, 3), 3},
	} {
		if got := tt.s.Len(); got != tt.want {
			t.Errorf("%s.Len(): got %d; want %d", tt.s.debug(), got, tt.want)
		}
	}
}

func TestDo(t *testing.T) {
	s := Of[int]()
	s.Do(func(n int) bool {
		t.Fatalf("%s.Do yielded element %d", s.debug(), n)
		panic("unreached")
	})
	s = emptyOf[int]()
	s.Do(func(n int) bool {
		t.Fatalf("%s.Do yielded element %d", s.debug(), n)
		panic("unreached")
	})

	s = Of(1, 2, 3, 4, 5)
	var saw []int
	s.Do(func(n int) bool {
		saw = append(saw, n)
		return true
	})
	sort.Ints(saw)
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(saw, want) {
		t.Errorf("%s.Do: saw elements %v; want %v", s.debug(), saw, want)
	}

	iter := 0
	s.Do(func(int) bool {
		iter++
		switch {
		case iter < 2:
			return true
		case iter == 2:
			return false
		default:
			t.Fatalf("%s.Do called fn after it returned false", s.debug())
			panic("unreached")
		}
	})
	if iter != 2 {
		t.Errorf("%s.Do called fn %d times", s, iter)
	}
}

func TestUnion(t *testing.T) {
	for _, tt := range []struct {
		s1   *Set[int]
		s2   *Set[int]
		want *Set[int]
	}{
		{Of[int](), Of[int](), Of[int]()},
		{Of[int](), emptyOf[int](), Of[int]()},
		{Of[int](), Of(3), Of(3)},
		{emptyOf[int](), Of[int](), Of[int]()},
		{emptyOf[int](), emptyOf[int](), Of[int]()},
		{emptyOf[int](), Of(3), Of(3)},
		{Of(3), Of[int](), Of(3)},
		{Of(3), emptyOf[int](), Of(3)},
		{Of(3), Of(3), Of(3)},
		{Of(3), Of(4), Of(3, 4)},
		{Of(3, 4), Of(3), Of(3, 4)},
		{Of(3), Of(3, 4), Of(3, 4)},
		{Of(3, 4), Of(3, 5), Of(3, 4, 5)},
		{Of(3, 4), Of(5, 6), Of(3, 4, 5, 6)},
	} {
		got := Union(tt.s1, tt.s2)
		if !equal(got, tt.want) {
			t.Errorf(
				"Union(%s, %s): got %s; want %s",
				tt.s1.debug(), tt.s2.debug(), got.debug(), tt.want.debug(),
			)
		}
	}
}

func TestIntersection(t *testing.T) {
	for _, tt := range []struct {
		s1   *Set[int]
		s2   *Set[int]
		want *Set[int]
	}{
		{Of[int](), Of[int](), Of[int]()},
		{Of[int](), emptyOf[int](), Of[int]()},
		{Of[int](), Of(3), Of[int]()},
		{emptyOf[int](), Of[int](), Of[int]()},
		{emptyOf[int](), emptyOf[int](), Of[int]()},
		{emptyOf[int](), Of(3), Of[int]()},
		{Of(3), Of[int](), Of[int]()},
		{Of(3), emptyOf[int](), Of[int]()},
		{Of(3), Of(3), Of(3)},
		{Of(3), Of(4), Of[int]()},
		{Of(3, 4), Of(3), Of(3)},
		{Of(3), Of(3, 4), Of(3)},
		{Of(3, 4), Of(3, 5), Of(3)},
		{Of(3, 4, 5), Of(3, 4, 5), Of(3, 4, 5)},
		{Of(3, 4), Of(5, 6), Of[int]()},
	} {
		got := Intersection(tt.s1, tt.s2)
		if !equal(got, tt.want) {
			t.Errorf(
				"Intersection(%s, %s): got %s; want %s",
				tt.s1.debug(), tt.s2.debug(), got.debug(), tt.want.debug(),
			)
		}
	}
}

func TestDifference(t *testing.T) {
	for _, tt := range []struct {
		s1   *Set[int]
		s2   *Set[int]
		want *Set[int]
	}{
		{Of[int](), Of[int](), Of[int]()},
		{Of[int](), emptyOf[int](), Of[int]()},
		{Of[int](), Of(3), Of[int]()},
		{emptyOf[int](), Of[int](), Of[int]()},
		{emptyOf[int](), emptyOf[int](), Of[int]()},
		{emptyOf[int](), Of(3), Of[int]()},
		{Of(3), Of[int](), Of(3)},
		{Of(3), emptyOf[int](), Of(3)},
		{Of(3), Of(3), Of[int]()},
		{Of(3), Of(4), Of(3)},
		{Of(3, 4), Of(3), Of(4)},
		{Of(3), Of(3, 4), Of[int]()},
		{Of(3, 4), Of(3, 5), Of(4)},
		{Of(3, 4, 5), Of(3, 4, 5), Of[int]()},
		{Of(3, 4), Of(5, 6), Of(3, 4)},
	} {
		got := Difference(tt.s1, tt.s2)
		if !equal(got, tt.want) {
			t.Errorf(
				"Difference(%s, %s): got %s; want %s",
				tt.s1.debug(), tt.s2.debug(), got.debug(), tt.want.debug(),
			)
		}
	}
}

func (s *Set[E]) debug() string {
	var v E
	typeName := fmt.Sprintf("%T", v)
	if s.m == nil {
		return fmt.Sprintf("set[%s](nil)", typeName)
	}
	if len(s.m) == 0 {
		return fmt.Sprintf("set[%s](empty)", typeName)
	}
	return s.String()
}

// equal reports whether s1 and s2 are equal, distinguishing between nil and
// empty internal maps.
func equal[E comparable](s1, s2 *Set[E]) bool {
	if (s1.m == nil) != (s2.m == nil) {
		return false
	}
	return s1.Equal(s2)
}

func check[E comparable](t *testing.T, s *Set[E], want []E) {
	t.Helper()
	if want == nil {
		if s.m != nil {
			t.Errorf("got %s; want s.m == nil", s)
			return
		}
		return
	}
	if s.m == nil {
		t.Errorf("got s.m == nil; want %v", want)
		return
	}

	wantSet := make(map[E]struct{})
	for _, w := range want {
		if _, ok := wantSet[w]; ok {
			panic("want contains duplicate elements")
		}
		wantSet[w] = struct{}{}
	}

	if len(s.m) != len(want) {
		t.Errorf("got %s; want %v", s, want)
		return
	}
	for _, w := range want {
		if _, ok := s.m[w]; !ok {
			t.Errorf("got %s; want %v", s, want)
			return
		}
	}
}
