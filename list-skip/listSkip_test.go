package listskip

import "testing"

func TestListSkipSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		s := New[int](0, cmpInt)
		if s == nil {
			t.Fatalf("New(0, cmp) = nil, want non-nil")
		}
		if s.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", s.Len())
		}
	})

	t.Run("insert_delete_has_clear", func(t *testing.T) {
		s := New[int](4, cmpInt)
		if !s.Insert(2) || !s.Insert(1) || !s.Insert(3) {
			t.Fatalf("Insert new values should return true")
		}
		if s.Insert(2) {
			t.Fatalf("duplicate Insert should return false")
		}
		if !s.Has(1) {
			t.Fatalf("Has(1) = false, want true")
		}
		if !s.Delete(2) {
			t.Fatalf("Delete(2) = false, want true")
		}
		if s.Delete(9) {
			t.Fatalf("Delete(9) = true, want false")
		}
		s.Clear()
		if s.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", s.Len())
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		s := New[int](6, cmpInt)
		s.Insert(3)
		s.Insert(1)
		s.Insert(2)
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
