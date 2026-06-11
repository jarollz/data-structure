package maphash

import "testing"

func TestMapHashSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		m := New[int, string](0, hashInt, eqInt)
		if m == nil {
			t.Fatalf("New(0, hash, eq) = nil, want non-nil")
		}
		if m.Cap() < 16 {
			t.Fatalf("Cap() = %d, want >= 16", m.Cap())
		}
		if m.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", m.Len())
		}
	})

	t.Run("put_get_delete_has_clear_loadfactor", func(t *testing.T) {
		m := New[int, string](4, hashInt, eqInt)
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(1, "aa")
		if v, ok := m.Get(1); !ok || v != "aa" {
			t.Fatalf("Get(1) = (%q, %v), want (aa, true)", v, ok)
		}
		if !m.Has(2) {
			t.Fatalf("Has(2) = false, want true")
		}
		if !m.Delete(2) {
			t.Fatalf("Delete(2) = false, want true")
		}
		if m.Delete(9) {
			t.Fatalf("Delete(9) = true, want false")
		}
		if lf := m.LoadFactor(); lf <= 0 {
			t.Fatalf("LoadFactor() = %f, want > 0", lf)
		}
		m.Clear()
		if m.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", m.Len())
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		m := New[int, string](4, hashInt, eqInt)
		m.Put(1, "a")
		m.Put(2, "b")
		m.Put(3, "c")

		got := collectSeq2(m.All())
		if len(got) != 3 {
			t.Fatalf("All count = %d, want 3", len(got))
		}

		count := 0
		for range m.All() {
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
		t.Log("All order is unspecified; mutation during iteration is not safe by contract")
	})
}
