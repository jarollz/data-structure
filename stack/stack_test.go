package stack

import "testing"

func TestStackSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		s := New[int](0)
		if s == nil {
			t.Fatalf("New(0) = nil, want non-nil")
		}
		if s.Cap() < 16 {
			t.Fatalf("Cap() = %d, want >= 16", s.Cap())
		}
		if s.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", s.Len())
		}
	})

	t.Run("push_pop_peek_clear", func(t *testing.T) {
		s := New[int](2)
		if _, ok := s.Pop(); ok {
			t.Fatalf("Pop() ok = true on empty stack")
		}
		if _, ok := s.PeekTop(); ok {
			t.Fatalf("PeekTop() ok = true on empty stack")
		}
		if !s.Push(1) || !s.Push(2) {
			t.Fatalf("Push should return true")
		}
		if v, ok := s.PeekTop(); !ok || v != 2 {
			t.Fatalf("PeekTop() = (%d, %v), want (2, true)", v, ok)
		}
		if v, ok := s.Pop(); !ok || v != 2 {
			t.Fatalf("Pop() = (%d, %v), want (2, true)", v, ok)
		}
		s.Clear()
		if s.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", s.Len())
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		s := New[int](4)
		s.Push(1)
		s.Push(2)
		s.Push(3)
		got := collectSeq(s.Values())
		if len(got) != 3 || got[0] != 3 || got[1] != 2 || got[2] != 1 {
			t.Fatalf("Values order/count = %v, want [3 2 1]", got)
		}

		count := 0
		for range s.Values() {
			count++
			if count == 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("early-stop count = %d, want 2", count)
		}

		s.Clear()
		emptyCount := 0
		for range s.Values() {
			emptyCount++
		}
		if emptyCount != 0 {
			t.Fatalf("empty Values count = %d, want 0", emptyCount)
		}
		t.Log("mutation during iteration is not safe by contract")
	})
}
