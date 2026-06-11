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

	t.Run("id_stability_and_removed_holes", func(t *testing.T) {
		tree := New[string]("root")
		a, _ := tree.AddChild(0, "a")
		b, _ := tree.AddChild(0, "b")
		c, _ := tree.AddChild(a, "c")
		if !tree.RemoveSubtree(a) {
			t.Fatalf("RemoveSubtree(a) = false, want true")
		}
		if _, ok := tree.Get(a); ok {
			t.Fatalf("Get(a) ok = true after remove")
		}
		if _, ok := tree.Get(c); ok {
			t.Fatalf("Get(c) ok = true after remove")
		}
		d, ok := tree.AddChild(0, "d")
		if !ok || d <= b {
			t.Fatalf("AddChild(root, d) = (%d, %v), want new id > %d", d, ok, b)
		}
		if parent, ok := tree.Parent(d); !ok || parent != 0 {
			t.Fatalf("Parent(d) = (%d, %v), want (0, true)", parent, ok)
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

func TestTreeGeneralCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		type item struct {
			label string
			ref   *int
		}

		rootRef := 0
		aRef := 1
		bRef := 2
		cRef := 3

		tree := New[item](item{label: "root", ref: &rootRef})
		childA, ok := tree.AddChild(0, item{label: "a", ref: &aRef})
		if !ok {
			t.Fatalf("AddChild(root, a) failed")
		}
		childB, ok := tree.AddChild(0, item{label: "b", ref: &bRef})
		if !ok {
			t.Fatalf("AddChild(root, b) failed")
		}
		grandChild, ok := tree.AddChild(childA, item{label: "c", ref: &cRef})
		if !ok {
			t.Fatalf("AddChild(a, c) failed")
		}
		if !tree.RemoveSubtree(childA) {
			t.Fatalf("RemoveSubtree(childA) failed")
		}

		cloned := tree.Clone()
		if cloned == tree {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != tree.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), tree.Len())
		}

		if v, ok := cloned.Get(0); !ok || v.label != "root" || v.ref != &rootRef {
			t.Fatalf("clone Get(0) = (%+v, %v), want root value with shared ref", v, ok)
		}
		if v, ok := cloned.Get(childB); !ok || v.label != "b" || v.ref != &bRef {
			t.Fatalf("clone Get(childB) = (%+v, %v), want b value with shared ref", v, ok)
		}
		if _, ok := cloned.Get(childA); ok {
			t.Fatalf("clone Get(childA) ok = true, want false for removed node")
		}
		if _, ok := cloned.Get(grandChild); ok {
			t.Fatalf("clone Get(grandChild) ok = true, want false for removed node")
		}
		if parent, ok := cloned.Parent(childB); !ok || parent != 0 {
			t.Fatalf("clone Parent(childB) = (%d, %v), want (0, true)", parent, ok)
		}
		if count := cloned.ChildCount(0); count != 1 {
			t.Fatalf("clone ChildCount(0) = %d, want 1", count)
		}

		nextID, ok := cloned.AddChild(0, item{label: "d"})
		if !ok || nextID != 4 {
			t.Fatalf("clone AddChild(root, d) = (%d, %v), want (4, true)", nextID, ok)
		}
		if tree.Len() != 2 {
			t.Fatalf("original Len after clone mutation = %d, want 2", tree.Len())
		}
		if _, ok := tree.Get(nextID); ok {
			t.Fatalf("original Get(nextID) ok = true, want false")
		}
	})

	t.Run("clonewith_nil_and_custom_hook", func(t *testing.T) {
		type item struct{ label string }

		tree := New[item](item{label: "root"})
		childA, _ := tree.AddChild(0, item{label: "a"})
		tree.AddChild(childA, item{label: "leaf"})
		tree.AddChild(0, item{label: "b"})

		nilClone := tree.CloneWith(nil)
		nilValues := collectSeq(nilClone.PreOrder())
		if len(nilValues) != 4 || nilValues[0].label != "root" || nilValues[1].label != "a" || nilValues[2].label != "leaf" || nilValues[3].label != "b" {
			t.Fatalf("CloneWith(nil) PreOrder() = %v, want [root a leaf b]", nilValues)
		}

		calls := make([]string, 0, 4)
		cloned := tree.CloneWith(func(v item) item {
			calls = append(calls, v.label)
			return item{label: v.label + "!"}
		})
		if len(calls) != 4 || calls[0] != "root" || calls[1] != "a" || calls[2] != "leaf" || calls[3] != "b" {
			t.Fatalf("cloneValue call order = %v, want [root a leaf b]", calls)
		}

		got := collectSeq(cloned.PreOrder())
		if len(got) != 4 || got[0].label != "root!" || got[1].label != "a!" || got[2].label != "leaf!" || got[3].label != "b!" {
			t.Fatalf("CloneWith PreOrder() = %v, want [root! a! leaf! b!]", got)
		}
	})

	t.Run("clone_empty_after_root_removal_and_hook_calls", func(t *testing.T) {
		tree := New[int](1)
		if !tree.RemoveSubtree(0) {
			t.Fatalf("RemoveSubtree(0) = false, want true")
		}
		calls := 0
		cloned := tree.CloneWith(func(v int) int {
			calls++
			return v
		})
		if calls != 0 {
			t.Fatalf("empty CloneWith hook calls = %d, want 0", calls)
		}
		if cloned.Len() != 0 {
			t.Fatalf("cloned Len() = %d, want 0", cloned.Len())
		}
	})
}
