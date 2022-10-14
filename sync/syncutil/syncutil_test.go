package syncutil

import "testing"

func TestPool(t *testing.T) {
	// Just make sure it kind of works. This test fails if the Pool
	// implementation is a no-op, though.
	p := &Pool[int]{
		New: func() int { return 3 },
	}
	if got := p.Get(); got != 3 {
		t.Errorf("Pool.Get(): got %d; want 3 (from New)", got)
	}
	for i := 0; i < 10; i++ {
		p.Put(5)
		if p.Get() == 5 {
			return
		}
	}
	t.Fatal("Pool never gave back pooled value (5)")
}

// TODO(caleb): Map tests.
