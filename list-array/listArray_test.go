package listarray

import "testing"

func TestListArraySpec(t *testing.T) {
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

	t.Run("append_get_set_insert_delete", func(t *testing.T) {
		s := New[int](2)
		if !s.Append(10) || !s.Append(20) {
			t.Fatalf("Append should return true")
		}
		if v, ok := s.Get(1); !ok || v != 20 {
			t.Fatalf("Get(1) = (%d, %v), want (20, true)", v, ok)
		}
		if _, ok := s.Get(-1); ok {
			t.Fatalf("Get(-1) ok = true, want false")
		}
		if !s.Set(1, 30) {
			t.Fatalf("Set(1, 30) = false, want true")
		}
		if s.Set(5, 99) {
			t.Fatalf("Set(5, 99) = true, want false")
		}
		if !s.Insert(0, 5) {
			t.Fatalf("Insert at head should succeed")
		}
		if !s.Insert(s.Len(), 40) {
			t.Fatalf("Insert at tail should succeed")
		}
		if s.Insert(-1, 1) {
			t.Fatalf("Insert(-1, 1) = true, want false")
		}
		if v, ok := s.Delete(0); !ok || v != 5 {
			t.Fatalf("Delete(0) = (%d, %v), want (5, true)", v, ok)
		}
		if _, ok := s.Delete(99); ok {
			t.Fatalf("Delete(99) ok = true, want false")
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		s := New[int](4)
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
