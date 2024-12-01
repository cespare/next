package ordmap

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMap(t *testing.T) {
	m := new(Map[string, int])

	checkGet(t, m, "a", 0, false)
	m.Set("a", 0)
	checkAll(t, m, []keyVal[string, int]{{"a", 0}})
	m.Set("a", 3)
	checkAll(t, m, []keyVal[string, int]{{"a", 3}})

	m.Set("b", 3)
	m.Set("c", 10)
	checkAll(t, m, []keyVal[string, int]{{"a", 3}, {"b", 3}, {"c", 10}})
	m.Set("a", 3)
	checkAll(t, m, []keyVal[string, int]{{"b", 3}, {"c", 10}, {"a", 3}})
	m.Set("c", 0)
	checkAll(t, m, []keyVal[string, int]{{"b", 3}, {"a", 3}, {"c", 0}})

	// Delete non-existent.
	m.Delete("abc")
	checkAll(t, m, []keyVal[string, int]{{"b", 3}, {"a", 3}, {"c", 0}})

	// Delete front.
	m.Delete("b")
	checkAll(t, m, []keyVal[string, int]{{"a", 3}, {"c", 0}})
	// Delete middle.
	m.Set("b", 8)
	m.Delete("c")
	checkAll(t, m, []keyVal[string, int]{{"a", 3}, {"b", 8}})
	// Delete end.
	m.Delete("b")
	checkAll(t, m, []keyVal[string, int]{{"a", 3}})
	// Delete last remaining element.
	m.Delete("a")
	checkAll(t, m, nil)

	m.Set("x", 10)
	m.Set("y", 20)
	checkAll(t, m, []keyVal[string, int]{{"x", 10}, {"y", 20}})
}

func checkGet[K, V comparable](t *testing.T, m *Map[K, V], key K, want V, wantOK bool) {
	t.Helper()
	got, ok := m.Get(key)
	if got != want || ok != wantOK {
		t.Fatalf(
			"Get(%#v): got (%#v, %t); want (%#v, %t)",
			key,
			got, ok,
			want, wantOK,
		)
	}
}

type keyVal[K comparable, V any] struct {
	Key K
	Val V
}

func checkAll[K, V comparable](t *testing.T, m *Map[K, V], kvs []keyVal[K, V]) {
	t.Helper()
	var keys []K
	var vals []V
	for _, kv := range kvs {
		got, ok := m.Get(kv.Key)
		if !ok {
			t.Fatalf("Get(%#v): got !ok", kv.Key)
		}
		if got != kv.Val {
			t.Fatalf("Get(%#v): got %#v; want %#v", kv.Key, got, kv.Val)
		}
		keys = append(keys, kv.Key)
		vals = append(vals, kv.Val)
	}

	gotKeys := slices.Collect(m.Keys())
	if diff := cmp.Diff(gotKeys, keys); diff != "" {
		t.Fatalf("Keys gave incorrect sequence (-got, +want):\n%s", diff)
	}
	gotVals := slices.Collect(m.Values())
	if diff := cmp.Diff(gotVals, vals); diff != "" {
		t.Fatalf("Values gave incorrect sequence (-got, +want):\n%s", diff)
	}

	var gotKVs []keyVal[K, V]
	for k, v := range m.All() {
		gotKVs = append(gotKVs, keyVal[K, V]{k, v})
	}
	if diff := cmp.Diff(gotKVs, kvs); diff != "" {
		t.Fatalf("All gave incorrect sequence (-got, +want):\n%s", diff)
	}
}
