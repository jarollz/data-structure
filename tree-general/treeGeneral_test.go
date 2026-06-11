package treegeneral

import "testing"

func TestTreeGeneralSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		tree := New[string]("root")
		if tree == nil {
			t.Fatalf("New(root) = nil, want non-nil")
		}
		if tree.Len() != 1 {
			t.Fatalf("Len() = %d, want 1", tree.Len())
		}
		if v, ok := tree.Get(0); !ok || v != "root" {
			t.Fatalf("Get(0) = (%q, %v), want (root, true)", v, ok)
		}
	})

	t.Run("add_get_parent_childcount_remove", func(t *testing.T) {
		tree := New[int](10)
		if _, ok := tree.AddChild(999, 1); ok {
			t.Fatalf("AddChild on invalid parent should fail")
		}
		child, ok := tree.AddChild(0, 20)
		if !ok || child < 1 {
			t.Fatalf("AddChild(0, 20) = (%d, %v), want valid id and true", child, ok)
		}
		if v, ok := tree.Get(child); !ok || v != 20 {
			t.Fatalf("Get(child) = (%d, %v), want (20, true)", v, ok)
		}
		if p, ok := tree.Parent(child); !ok || p != 0 {
			t.Fatalf("Parent(child) = (%d, %v), want (0, true)", p, ok)
		}
		if c := tree.ChildCount(0); c != 1 {
			t.Fatalf("ChildCount(0) = %d, want 1", c)
		}
		if !tree.RemoveSubtree(child) {
			t.Fatalf("RemoveSubtree(child) = false, want true")
		}
		if tree.RemoveSubtree(999) {
			t.Fatalf("RemoveSubtree(999) = true, want false")
		}
		if !tree.RemoveSubtree(0) {
			t.Fatalf("RemoveSubtree(0) = false, want true")
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		tree := New[int](1)
		a, _ := tree.AddChild(0, 2)
		b, _ := tree.AddChild(0, 3)
		tree.AddChild(a, 4)
		tree.AddChild(b, 5)

		got := collectSeq(tree.PreOrder())
		if len(got) != 5 || got[0] != 1 {
			t.Fatalf("PreOrder count/order = %v, want root-first traversal", got)
		}

		count := 0
		for range tree.PreOrder() {
			count++
			if count == 3 {
				break
			}
		}
		if count != 3 {
			t.Fatalf("early-stop count = %d, want 3", count)
		}

		tree.RemoveSubtree(0)
		emptyCount := 0
		for range tree.PreOrder() {
			emptyCount++
		}
		if emptyCount != 0 {
			t.Fatalf("empty PreOrder count = %d, want 0", emptyCount)
		}
		t.Log("mutation during iteration is not safe by contract")
	})
}
