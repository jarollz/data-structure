package maptreeredblack

import "testing"

func TestMapTreeRedBlackSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		m := New[int, string](cmpInt)
		if m == nil {
			t.Fatalf("New(cmp) = nil, want non-nil")
		}
		if m.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", m.Len())
		}
	})

	t.Run("put_get_delete_has_min_max_clear", func(t *testing.T) {
		m := New[int, string](cmpInt)
		if _, _, ok := m.Min(); ok {
			t.Fatalf("Min() ok = true on empty map")
		}
		if _, _, ok := m.Max(); ok {
			t.Fatalf("Max() ok = true on empty map")
		}
		m.Put(2, "b")
		m.Put(1, "a")
		m.Put(3, "c")
		m.Put(2, "bb")
		if v, ok := m.Get(2); !ok || v != "bb" {
			t.Fatalf("Get(2) = (%q, %v), want (bb, true)", v, ok)
		}
		if !m.Has(1) {
			t.Fatalf("Has(1) = false, want true")
		}
		if !m.Delete(1) {
			t.Fatalf("Delete(1) = false, want true")
		}
		if m.Delete(9) {
			t.Fatalf("Delete(9) = true, want false")
		}
		m.Clear()
		if m.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", m.Len())
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		m := New[int, string](cmpInt)
		m.Put(3, "c")
		m.Put(1, "a")
		m.Put(2, "b")

		got := collectSeq2(m.All())
		if len(got) != 3 {
			t.Fatalf("All count = %d, want 3", len(got))
		}

		last := -1
		count := 0
		for k := range m.All() {
			if count > 0 && k < last {
				t.Fatalf("All keys not ascending: last=%d now=%d", last, k)
			}
			last = k
			count++
			if count == 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("early-stop count = %d, want 2", count)
		}

		m.Clear()
		emptyCount := 0
		for range m.All() {
			emptyCount++
		}
		if emptyCount != 0 {
			t.Fatalf("empty All count = %d, want 0", emptyCount)
		}
		t.Log("mutation during iteration is not safe by contract")
	})
}
