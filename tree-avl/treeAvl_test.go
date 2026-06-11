package treeavl

import "testing"

func TestTreeAvlSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		tree := New[int](cmpInt)
		if tree == nil {
			t.Fatalf("New(cmp) = nil, want non-nil")
		}
		if tree.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", tree.Len())
		}
	})

	t.Run("insert_delete_has_min_max_clear", func(t *testing.T) {
		tree := New[int](cmpInt)
		if _, ok := tree.Min(); ok {
			t.Fatalf("Min() ok = true on empty tree")
		}
		if _, ok := tree.Max(); ok {
			t.Fatalf("Max() ok = true on empty tree")
		}
		if !tree.Insert(2) || !tree.Insert(1) || !tree.Insert(3) {
			t.Fatalf("Insert new values should return true")
		}
		if tree.Insert(2) {
			t.Fatalf("duplicate Insert should return false")
		}
		if !tree.Has(1) {
			t.Fatalf("Has(1) = false, want true")
		}
		if !tree.Delete(2) {
			t.Fatalf("Delete(2) = false, want true")
		}
		if tree.Delete(9) {
			t.Fatalf("Delete(9) = true, want false")
		}
		tree.Clear()
		if tree.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", tree.Len())
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		tree := New[int](cmpInt)
		tree.Insert(3)
		tree.Insert(1)
		tree.Insert(2)
		got := collectSeq(tree.InOrder())
		if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
			t.Fatalf("InOrder order/count = %v, want [1 2 3]", got)
		}

		count := 0
		for range tree.InOrder() {
			count++
			if count == 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("early-stop count = %d, want 2", count)
		}

		tree.Clear()
		emptyCount := 0
		for range tree.InOrder() {
			emptyCount++
		}
		if emptyCount != 0 {
			t.Fatalf("empty InOrder count = %d, want 0", emptyCount)
		}
		t.Log("mutation during iteration is not safe by contract")
	})
}
