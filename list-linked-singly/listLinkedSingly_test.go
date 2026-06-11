package listlinkedsingly

import "testing"

func TestListLinkedSinglySpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		s := New[int]()
		if s == nil {
			t.Fatalf("New() = nil, want non-nil")
		}
		if s.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", s.Len())
		}
	})

	t.Run("push_pop_append_deletefirst", func(t *testing.T) {
		s := New[int]()
		if _, ok := s.PopFront(); ok {
			t.Fatalf("PopFront() ok = true on empty list")
		}
		s.PushFront(2)
		s.PushFront(1)
		s.Append(3)
		if !s.DeleteFirst(func(v int) bool { return v == 2 }) {
			t.Fatalf("DeleteFirst existing value should return true")
		}
		if s.DeleteFirst(func(v int) bool { return v == 9 }) {
			t.Fatalf("DeleteFirst missing value should return false")
		}
		if v, ok := s.PopFront(); !ok || v != 1 {
			t.Fatalf("PopFront() = (%d, %v), want (1, true)", v, ok)
		}
		s.Clear()
		if s.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", s.Len())
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		s := New[int]()
		s.Append(1)
		s.Append(2)
		s.Append(3)
		got := collectSeq(s.Values())
		if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
			t.Fatalf("Values order/count = %v, want [1 2 3]", got)
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
