// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the THIRD-PARTY-LICENSE.txt file.

package heap

import (
	"math/rand"
	"testing"
)

type intElem struct {
	v int
	i int
}

func newIntHeap() *Heap[*intElem] {
	h := New(func(e0, e1 *intElem) bool { return e0.v < e1.v })
	h.SetIndex = func(e *intElem, i int) { e.i = i }
	return h
}

func verify(t *testing.T, h *Heap[*intElem]) {
	t.Helper()
	checkInvariant(t, h, 0)
	for i, e := range h.s {
		if e.i != i {
			t.Fatalf("heap index not in sync: [%d].i = %d", i, e.i)
		}
	}
}

func checkInvariant(t *testing.T, h *Heap[*intElem], i int) {
	t.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.Less(h.s[j1], h.s[i]) {
			t.Fatalf("heap invariant invalidated [%d] = %v > [%d] = %v", i, h.s[i], j1, h.s[j1])
		}
		checkInvariant(t, h, j1)
	}
	if j2 < n {
		if h.Less(h.s[j2], h.s[i]) {
			t.Fatalf("heap invariant invalidated [%d] = %v > [%d] = %v", i, h.s[i], j2, h.s[j2])
		}
		checkInvariant(t, h, j2)
	}
}

func TestHeapInit0(t *testing.T) {
	h := newIntHeap()
	s := make([]*intElem, 20)
	for i := range s {
		s[i] = &intElem{v: 0} // all elements the same
	}
	h.Init(s)
	verify(t, h)

	for i := 1; h.Len() > 0; i++ {
		e := h.Pop()
		verify(t, h)
		if e.v != 0 {
			t.Errorf("%d.th pop got %d; want 0", i, e.v)
		}
	}
}

func TestHeap1(t *testing.T) {
	h := newIntHeap()
	s := make([]*intElem, 20)
	for i := range s {
		s[i] = &intElem{v: 20 - i}
	}
	h.Init(s)
	verify(t, h)

	for i := 1; h.Len() > 0; i++ {
		e := h.Pop()
		verify(t, h)
		if e.v != i {
			t.Errorf("%d.th pop got %d; want %d", i, e.v, i)
		}
	}
}

func TestHeap(t *testing.T) {
	h := newIntHeap()
	verify(t, h)

	var s []*intElem
	for i := 20; i > 10; i-- {
		s = append(s, &intElem{v: i})
	}
	h.Init(s)
	verify(t, h)

	for i := 10; i > 0; i-- {
		h.Push(&intElem{v: i})
		verify(t, h)
	}

	for i := 1; h.Len() > 0; i++ {
		e := h.Pop()
		if i < 20 {
			h.Push(&intElem{v: 20 + i})
		}
		verify(t, h)
		if e.v != i {
			t.Errorf("%d.th pop got %d; want %d", i, e.v, i)
		}
	}
}

func TestHeapRemove0(t *testing.T) {
	h := newIntHeap()
	for i := 0; i < 10; i++ {
		h.Push(&intElem{v: i})
	}
	verify(t, h)

	for h.Len() > 0 {
		i := h.Len() - 1
		e := h.Remove(i)
		if e.v != i {
			t.Errorf("Remove(%d) got %d; want %d", i, e.v, i)
		}
		verify(t, h)
	}
}

func TestHeapRemove1(t *testing.T) {
	h := newIntHeap()
	for i := 0; i < 10; i++ {
		h.Push(&intElem{v: i})
	}
	verify(t, h)

	for i := 0; h.Len() > 0; i++ {
		e := h.Remove(0)
		if e.v != i {
			t.Errorf("Remove(0) got %d; want %d", e.v, i)
		}
		verify(t, h)
	}
}

func TestHeapRemove2(t *testing.T) {
	h := newIntHeap()
	for i := 0; i < 10; i++ {
		h.Push(&intElem{v: i})
	}
	verify(t, h)

	m := make(map[int]struct{})
	for h.Len() > 0 {
		m[h.Remove((h.Len()-1)/2).v] = struct{}{}
		verify(t, h)
	}

	if len(m) != 10 {
		t.Errorf("len(m) = %d; want %d", len(m), 10)
	}
	for i := 0; i < len(m); i++ {
		if _, ok := m[i]; !ok {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func TestHeapFix(t *testing.T) {
	h := newIntHeap()
	verify(t, h)

	for i := 200; i > 0; i -= 10 {
		h.Push(&intElem{v: i})
	}
	verify(t, h)

	if h.s[0].v != 10 {
		t.Fatalf("got head %d; want 10", h.s[0].v)
	}
	h.s[0].v = 210
	h.Fix(0)
	verify(t, h)

	for i := 100; i > 0; i-- {
		j := rand.Intn(h.Len())
		if i&1 == 0 {
			h.s[j].v *= 2
		} else {
			h.s[j].v /= 2
		}
		h.Fix(j)
		verify(t, h)
	}
}

type myHeap []int

func (h *myHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *myHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *myHeap) Len() int {
	return len(*h)
}

func (h *myHeap) Pop() (v any) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *myHeap) Push(v any) {
	*h = append(*h, v.(int))
}

func (h myHeap) verify(t *testing.T, i int) {
	t.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.Less(j1, i) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h[i], j1, h[j1])
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if h.Less(j2, i) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h[i], j1, h[j2])
			return
		}
		h.verify(t, j2)
	}
}

func TestInit0(t *testing.T) {
	h := new(myHeap)
	for i := 20; i > 0; i-- {
		h.Push(0) // all elements are the same
	}
	Init(h)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := Pop(h).(int)
		h.verify(t, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	h := new(myHeap)
	for i := 20; i > 0; i-- {
		h.Push(i) // all elements are different
	}
	Init(h)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := Pop(h).(int)
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func Test(t *testing.T) {
	h := new(myHeap)
	h.verify(t, 0)

	for i := 20; i > 10; i-- {
		h.Push(i)
	}
	Init(h)
	h.verify(t, 0)

	for i := 10; i > 0; i-- {
		Push(h, i)
		h.verify(t, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x := Pop(h).(int)
		if i < 20 {
			Push(h, 20+i)
		}
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	h := new(myHeap)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		x := Remove(h, i).(int)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		h.verify(t, 0)
	}
}

func TestRemove1(t *testing.T) {
	h := new(myHeap)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for i := 0; h.Len() > 0; i++ {
		x := Remove(h, 0).(int)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		h.verify(t, 0)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	h := new(myHeap)
	for i := 0; i < N; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		m[Remove(h, (h.Len()-1)/2).(int)] = true
		h.verify(t, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func BenchmarkDup(b *testing.B) {
	const n = 10000
	h := make(myHeap, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			Push(&h, 0) // all elements are the same
		}
		for h.Len() > 0 {
			Pop(&h)
		}
	}
}

func TestFix(t *testing.T) {
	h := new(myHeap)
	h.verify(t, 0)

	for i := 200; i > 0; i -= 10 {
		Push(h, i)
	}
	h.verify(t, 0)

	if (*h)[0] != 10 {
		t.Fatalf("Expected head to be 10, was %d", (*h)[0])
	}
	(*h)[0] = 210
	Fix(h, 0)
	h.verify(t, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(h.Len())
		if i&1 == 0 {
			(*h)[elem] *= 2
		} else {
			(*h)[elem] /= 2
		}
		Fix(h, elem)
		h.verify(t, 0)
	}
}

type benchElem struct {
	v int
}

type benchHeap []*benchElem

func (h *benchHeap) Less(i, j int) bool {
	return (*h)[i].v < (*h)[j].v
}

func (h *benchHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *benchHeap) Len() int {
	return len(*h)
}

func (h *benchHeap) Pop() (v any) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *benchHeap) Push(v any) {
	*h = append(*h, v.(*benchElem))
}

func BenchmarkHeapInterface(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	var h benchHeap
	for i := 0; i < 1000; i++ {
		h.Push(&benchElem{v: rng.Intn(1000)})
	}
	Init(&h)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Push(&h, &benchElem{v: rng.Intn(1000)})
		_ = Pop(&h).(*benchElem)
	}
}

func BenchmarkHeapGeneric(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	h := New(func(e0, e1 *benchElem) bool { return e0.v < e1.v })
	s := make([]*benchElem, 1000)
	for i := range s {
		s[i] = &benchElem{v: rng.Intn(1000)}
	}
	h.Init(s)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Push(&benchElem{v: rng.Intn(1000)})
		_ = h.Pop()
	}
}
