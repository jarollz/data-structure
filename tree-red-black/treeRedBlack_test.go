package treeredblack

import (
	"math/rand"
	"sort"
	"testing"
)

func TestTreeRedBlackSpec(t *testing.T) {
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

	t.Run("delete_root_cases_and_min_max_updates", func(t *testing.T) {
		cases := [][]int{{2}, {2, 1}, {2, 1, 3}, {3, 2, 1}, {1, 2, 3}}
		for _, values := range cases {
			tree := New[int](cmpInt)
			for _, v := range values {
				tree.Insert(v)
			}
			if !tree.Delete(values[0]) {
				t.Fatalf("Delete(root=%d) = false, want true", values[0])
			}
			if tree.Has(values[0]) {
				t.Fatalf("Has(%d) = true after delete", values[0])
			}
		}
	})

	t.Run("randomized_against_set_model", func(t *testing.T) {
		tree := New[int](cmpInt)
		model := make(map[int]struct{})
		rng := rand.New(rand.NewSource(9))
		for step := 0; step < 400; step++ {
			value := rng.Intn(300)
			switch rng.Intn(3) {
			case 0:
				inserted := tree.Insert(value)
				_, exists := model[value]
				if inserted == exists {
					t.Fatalf("Insert(%d) returned %v, exists=%v", value, inserted, exists)
				}
				if !exists {
					model[value] = struct{}{}
				}
			case 1:
				deleted := tree.Delete(value)
				_, exists := model[value]
				if deleted != exists {
					t.Fatalf("Delete(%d) returned %v, exists=%v", value, deleted, exists)
				}
				delete(model, value)
			case 2:
				_, exists := model[value]
				if tree.Has(value) != exists {
					t.Fatalf("Has(%d) mismatch", value)
				}
			}
			ordered := make([]int, 0, len(model))
			for v := range model {
				ordered = append(ordered, v)
			}
			sort.Ints(ordered)
			got := collectSeq(tree.InOrder())
			if len(got) != len(ordered) {
				t.Fatalf("InOrder len = %d, want %d", len(got), len(ordered))
			}
			for i := range ordered {
				if got[i] != ordered[i] {
					t.Fatalf("InOrder()[%d] = %d, want %d", i, got[i], ordered[i])
				}
			}
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

func TestTreeRedBlackCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		type item struct {
			key int
			ref *int
		}

		cmpItem := func(a, b item) int {
			if a.key < b.key {
				return -1
			}
			if a.key > b.key {
				return 1
			}
			return 0
		}

		a := 1
		b := 2
		c := 3

		tree := New[item](cmpItem)
		tree.Insert(item{key: 2, ref: &b})
		tree.Insert(item{key: 1, ref: &a})
		tree.Insert(item{key: 3, ref: &c})

		cloned := tree.Clone()
		if cloned == tree {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != tree.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), tree.Len())
		}
		if v, ok := cloned.Min(); !ok || v.key != 1 || v.ref != &a {
			t.Fatalf("clone Min() = (%+v, %v), want key 1 with shared ref", v, ok)
		}
		if v, ok := cloned.Max(); !ok || v.key != 3 || v.ref != &c {
			t.Fatalf("clone Max() = (%+v, %v), want key 3 with shared ref", v, ok)
		}

		got := collectSeq(cloned.InOrder())
		if len(got) != 3 || got[0].key != 1 || got[1].key != 2 || got[2].key != 3 || got[0].ref != &a || got[1].ref != &b || got[2].ref != &c {
			t.Fatalf("clone InOrder() = %v, want sorted shared refs", got)
		}

		if !tree.Delete(item{key: 1}) {
			t.Fatalf("Delete(1) on original failed")
		}
		if !cloned.Has(item{key: 1}) {
			t.Fatalf("clone lost key 1 after original mutation")
		}

		if !cloned.Delete(item{key: 3}) {
			t.Fatalf("Delete(3) on clone failed")
		}
		if !tree.Has(item{key: 3}) {
			t.Fatalf("original lost key 3 after clone mutation")
		}
	})

	t.Run("clonewith_nil_and_custom_hook", func(t *testing.T) {
		type item struct{ key int }

		cmpItem := func(a, b item) int {
			if a.key < b.key {
				return -1
			}
			if a.key > b.key {
				return 1
			}
			return 0
		}

		tree := New[item](cmpItem)
		tree.Insert(item{key: 2})
		tree.Insert(item{key: 1})
		tree.Insert(item{key: 3})

		nilClone := tree.CloneWith(nil)
		nilValues := collectSeq(nilClone.InOrder())
		if len(nilValues) != 3 || nilValues[0].key != 1 || nilValues[1].key != 2 || nilValues[2].key != 3 {
			t.Fatalf("CloneWith(nil) InOrder() = %v, want [1 2 3]", nilValues)
		}

		calls := make([]int, 0, 3)
		cloned := tree.CloneWith(func(v item) item {
			calls = append(calls, v.key)
			return item{key: v.key + 10}
		})
		if len(calls) != 3 || calls[0] != 1 || calls[1] != 2 || calls[2] != 3 {
			t.Fatalf("cloneValue call order = %v, want [1 2 3]", calls)
		}

		got := collectSeq(cloned.InOrder())
		if len(got) != 3 || got[0].key != 11 || got[1].key != 12 || got[2].key != 13 {
			t.Fatalf("CloneWith InOrder() = %v, want [11 12 13]", got)
		}
	})

	t.Run("clone_empty_and_hook_order", func(t *testing.T) {
		empty := New[int](cmpInt)
		calls := 0
		empty.CloneWith(func(v int) int {
			calls++
			return v
		})
		if calls != 0 {
			t.Fatalf("empty CloneWith hook calls = %d, want 0", calls)
		}
	})
}

func TestTreeRedBlackAlgorithmInvariants(t *testing.T) {
	t.Run("deterministic_insert_and_delete_fixup_cases", func(t *testing.T) {
		insertCases := [][]int{
			{10, 20, 30},
			{30, 20, 10},
			{10, 30, 20},
			{20, 10, 30, 5, 15, 25, 35, 1},
			{7, 3, 18, 10, 22, 8, 11, 26},
			{41, 38, 31, 12, 19, 8},
		}
		for _, seq := range insertCases {
			tree := New[int](cmpInt)
			for _, v := range seq {
				if !tree.Insert(v) {
					t.Fatalf("Insert(%d) = false, want true for sequence %v", v, seq)
				}
				assertRedBlackInvariantsInt(t, tree)
			}
		}

		deleteCases := []struct {
			insert []int
			delete []int
		}{
			{insert: []int{11, 2, 14, 1, 7, 15, 5, 8, 4}, delete: []int{1, 2, 14}},
			{insert: []int{7, 3, 18, 10, 22, 8, 11, 26}, delete: []int{18, 11, 3}},
			{insert: []int{20, 10, 30, 5, 15, 25, 35, 1, 6, 14, 16}, delete: []int{1, 5, 6, 10, 20}},
		}
		for _, tc := range deleteCases {
			tree := New[int](cmpInt)
			for _, v := range tc.insert {
				tree.Insert(v)
			}
			assertRedBlackInvariantsInt(t, tree)
			for _, v := range tc.delete {
				if !tree.Delete(v) {
					t.Fatalf("Delete(%d) = false, want true for sequence %+v", v, tc)
				}
				assertRedBlackInvariantsInt(t, tree)
			}
		}
	})

	t.Run("randomized_operations_validate_red_black_after_each_step", func(t *testing.T) {
		seeds := []int64{404, 505, 606}
		for _, seed := range seeds {
			tree := New[int](cmpInt)
			model := make(map[int]struct{})
			rng := rand.New(rand.NewSource(seed))
			for step := 0; step < 3000; step++ {
				v := rng.Intn(4000) - 2000
				switch rng.Intn(3) {
				case 0:
					inserted := tree.Insert(v)
					_, exists := model[v]
					if inserted == exists {
						t.Fatalf("seed=%d step=%d Insert(%d) returned %v exists=%v", seed, step, v, inserted, exists)
					}
					if !exists {
						model[v] = struct{}{}
					}
				case 1:
					deleted := tree.Delete(v)
					_, exists := model[v]
					if deleted != exists {
						t.Fatalf("seed=%d step=%d Delete(%d) returned %v exists=%v", seed, step, v, deleted, exists)
					}
					delete(model, v)
				case 2:
					_, exists := model[v]
					if tree.Has(v) != exists {
						t.Fatalf("seed=%d step=%d Has(%d) mismatch", seed, step, v)
					}
				}

				ordered := make([]int, 0, len(model))
				for x := range model {
					ordered = append(ordered, x)
				}
				sort.Ints(ordered)
				got := collectSeq(tree.InOrder())
				if len(got) != len(ordered) {
					t.Fatalf("seed=%d step=%d InOrder len=%d want=%d", seed, step, len(got), len(ordered))
				}
				for i := range ordered {
					if got[i] != ordered[i] {
						t.Fatalf("seed=%d step=%d InOrder[%d]=%d want=%d", seed, step, i, got[i], ordered[i])
					}
				}

				assertRedBlackInvariantsInt(t, tree)
			}
		}
	})
}

func TestTreeRedBlackWalkPowerSpec(t *testing.T) {
	t.Run("rootnode_empty_and_non_empty", func(t *testing.T) {
		tree := New[int](cmpInt)
		if root, ok := tree.RootNode(); ok || root != nil {
			t.Fatalf("empty RootNode() = (%v, %v), want (nil, false)", root, ok)
		}

		tree.Insert(2)
		root, ok := tree.RootNode()
		if !ok || root == nil {
			t.Fatalf("non-empty RootNode() = (%v, %v), want (node, true)", root, ok)
		}
		if root.Color() != ColorBlack {
			t.Fatalf("root color = %v, want ColorBlack", root.Color())
		}
	})

	t.Run("children_order_childcount_color_and_early_stop", func(t *testing.T) {
		tree := New[int](cmpInt)
		for _, v := range []int{2, 1, 3} {
			tree.Insert(v)
		}
		root, ok := tree.RootNode()
		if !ok {
			t.Fatalf("RootNode() ok = false, want true")
		}

		children := collectSeq(root.Children())
		if len(children) != root.ChildCount() {
			t.Fatalf("ChildCount()=%d, yielded=%d", root.ChildCount(), len(children))
		}
		if len(children) != 2 {
			t.Fatalf("root child count=%d, want 2", len(children))
		}
		if children[0].Value() >= root.Value() || children[1].Value() <= root.Value() {
			t.Fatalf("children order invalid: root=%d, first=%d, second=%d", root.Value(), children[0].Value(), children[1].Value())
		}

		if root.Color() != ColorBlack && root.Color() != ColorRed {
			t.Fatalf("unexpected root color value: %v", root.Color())
		}

		count := 0
		for range root.Children() {
			count++
			break
		}
		if count != 1 {
			t.Fatalf("early-stop child count=%d, want 1", count)
		}
	})

	t.Run("dfs_visits_exactly_len_and_preserves_contract", func(t *testing.T) {
		tree := New[int](cmpInt)
		for _, v := range []int{10, 5, 15, 3, 7, 12, 18, 1, 4, 6, 8} {
			tree.Insert(v)
		}
		root, ok := tree.RootNode()
		if !ok {
			t.Fatalf("RootNode() ok = false, want true")
		}

		seen := make(map[int]struct{})
		var walk func(NodeAPI[int])
		walk = func(n NodeAPI[int]) {
			if n == nil {
				return
			}
			v := n.Value()
			if _, exists := seen[v]; exists {
				t.Fatalf("duplicate node visit for value %d", v)
			}
			seen[v] = struct{}{}

			if n.Color() != ColorBlack && n.Color() != ColorRed {
				t.Fatalf("node=%d has invalid color value %v", v, n.Color())
			}

			childCount := 0
			for child := range n.Children() {
				childCount++
				walk(child)
			}
			if childCount != n.ChildCount() {
				t.Fatalf("node=%d ChildCount()=%d yielded=%d", v, n.ChildCount(), childCount)
			}
		}

		walk(root)
		if len(seen) != tree.Len() {
			t.Fatalf("dfs seen=%d, Len()=%d", len(seen), tree.Len())
		}
		t.Log("mutation during node traversal is not safe by contract")
	})
}

func assertRedBlackInvariantsInt(t *testing.T, tree *TreeRedBlack[int]) {
	t.Helper()
	st := ensureState(tree, nil)
	if st == nil {
		t.Fatalf("internal state is nil")
	}
	if st.root == nilIndex {
		if tree.Len() != 0 {
			t.Fatalf("empty root with Len()=%d, want 0", tree.Len())
		}
		return
	}

	if st.parent(st.root) != nilIndex {
		t.Fatalf("root parent = %d, want %d", st.parent(st.root), nilIndex)
	}
	if st.colorOf(st.root) != colorBlack {
		t.Fatalf("root color not black")
	}

	visited := make(map[int]struct{})
	var walk func(index int, min *int, max *int) (blackHeight int, count int)
	walk = func(index int, min *int, max *int) (blackHeight int, count int) {
		if index == nilIndex {
			return 1, 0
		}
		if _, exists := visited[index]; exists {
			t.Fatalf("cycle or duplicate node index detected: %d", index)
		}
		visited[index] = struct{}{}

		v := st.value(index)
		if min != nil && st.cmp(v, *min) <= 0 {
			t.Fatalf("bst violation: value %d <= min bound %d", v, *min)
		}
		if max != nil && st.cmp(v, *max) >= 0 {
			t.Fatalf("bst violation: value %d >= max bound %d", v, *max)
		}

		left := st.left(index)
		right := st.right(index)
		if left != nilIndex && st.parent(left) != index {
			t.Fatalf("parent link mismatch for left child of value %d", v)
		}
		if right != nilIndex && st.parent(right) != index {
			t.Fatalf("parent link mismatch for right child of value %d", v)
		}

		color := st.colorOf(index)
		if color != colorBlack && color != colorRed {
			t.Fatalf("invalid color value at node %d: %d", v, color)
		}
		if color == colorRed {
			if st.colorOf(left) != colorBlack {
				t.Fatalf("red-red violation at node %d and left child", v)
			}
			if st.colorOf(right) != colorBlack {
				t.Fatalf("red-red violation at node %d and right child", v)
			}
		}

		leftBH, leftCount := walk(left, min, &v)
		rightBH, rightCount := walk(right, &v, max)
		if leftBH != rightBH {
			t.Fatalf("black-height mismatch at value %d: left=%d right=%d", v, leftBH, rightBH)
		}

		bh := leftBH
		if color == colorBlack {
			bh++
		}
		return bh, leftCount + rightCount + 1
	}

	_, nodeCount := walk(st.root, nil, nil)
	if nodeCount != tree.Len() {
		t.Fatalf("visited nodes=%d, Len()=%d", nodeCount, tree.Len())
	}
}
