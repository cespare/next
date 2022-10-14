package atomicutil

import "testing"

func TestInconsistent(t *testing.T) {
	t.Skip("fail")
	var v Value[any]
	v.Store(3)
	v.Store(4)
	v.Store("a")
}
