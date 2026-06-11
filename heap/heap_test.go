package heap

import "testing"

func TestHeapSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		h := New[int](0, cmpInt)
		if h == nil {
			t.Fatalf("New(0, cmp) = nil, want non-nil")
		}
		if h.Cap() < 16 {
			t.Fatalf("Cap() = %d, want >= 16", h.Cap())
		}
		if h.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", h.Len())
		}
	})

	t.Run("push_pop_peek_clear", func(t *testing.T) {
		h := New[int](2, cmpInt)
		if _, ok := h.PopTop(); ok {
			t.Fatalf("PopTop() ok = true on empty heap")
		}
		if _, ok := h.PeekTop(); ok {
			t.Fatalf("PeekTop() ok = true on empty heap")
		}
		h.Push(3)
		h.Push(1)
		h.Push(2)
		if v, ok := h.PeekTop(); !ok || v != 1 {
			t.Fatalf("PeekTop() = (%d, %v), want (1, true)", v, ok)
		}
		if v, ok := h.PopTop(); !ok || v != 1 {
			t.Fatalf("PopTop() = (%d, %v), want (1, true)", v, ok)
		}
		h.Clear()
		if h.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", h.Len())
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		h := New[int](4, cmpInt)
		h.Push(3)
		h.Push(1)
		h.Push(2)
		got := collectSeq(h.Values())
		if len(got) != 3 {
			t.Fatalf("Values count = %d, want 3", len(got))
		}

		count := 0
		for range h.Values() {
			count++
			if count == 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("early-stop count = %d, want 2", count)
		}

		h.Clear()
		emptyCount := 0
		for range h.Values() {
			emptyCount++
		}
		if emptyCount != 0 {
			t.Fatalf("empty Values count = %d, want 0", emptyCount)
		}
		t.Log("Values order is internal array order; mutation during iteration is not safe by contract")
	})
}
