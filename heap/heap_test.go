package heap

import (
	"math/rand"
	"sort"
	"testing"
)

func TestHeapSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		h := New[int](0, cmpInt)
		if h == nil {
			t.Fatalf("New(0, cmp) = nil, want non-nil")
		}
		if h.Cap() < 16 {
			t.Fatalf("Cap() = %d, want >= 16", h.Cap())
		}
		if h.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", h.Len())
		}
	})

	t.Run("push_pop_peek_clear", func(t *testing.T) {
		h := New[int](2, cmpInt)
		if _, ok := h.PopTop(); ok {
			t.Fatalf("PopTop() ok = true on empty heap")
		}
		if _, ok := h.PeekTop(); ok {
			t.Fatalf("PeekTop() ok = true on empty heap")
		}
		h.Push(3)
		h.Push(1)
		h.Push(2)
		if v, ok := h.PeekTop(); !ok || v != 1 {
			t.Fatalf("PeekTop() = (%d, %v), want (1, true)", v, ok)
		}
		if v, ok := h.PopTop(); !ok || v != 1 {
			t.Fatalf("PopTop() = (%d, %v), want (1, true)", v, ok)
		}
		h.Clear()
		if h.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", h.Len())
		}
	})

	t.Run("poptop_returns_monotonic_order", func(t *testing.T) {
		h := New[int](4, cmpInt)
		values := []int{7, 1, 4, 4, 9, 2, 6}
		for _, v := range values {
			h.Push(v)
		}
		sort.Ints(values)
		for i, want := range values {
			got, ok := h.PopTop()
			if !ok || got != want {
				t.Fatalf("PopTop #%d = (%d, %v), want (%d, true)", i, got, ok, want)
			}
		}
		if _, ok := h.PopTop(); ok {
			t.Fatalf("PopTop() ok = true after draining heap")
		}
	})

	t.Run("randomized_against_sorted_model", func(t *testing.T) {
		h := New[int](4, cmpInt)
		model := make([]int, 0)
		rng := rand.New(rand.NewSource(3))
		for step := 0; step < 400; step++ {
			switch rng.Intn(3) {
			case 0:
				value := rng.Intn(1000)
				h.Push(value)
				model = append(model, value)
				sort.Ints(model)
			case 1:
				got, ok := h.PopTop()
				if len(model) == 0 {
					if ok {
						t.Fatalf("PopTop() ok = true on empty model")
					}
					continue
				}
				want := model[0]
				model = model[1:]
				if !ok || got != want {
					t.Fatalf("PopTop() = (%d, %v), want (%d, true)", got, ok, want)
				}
			case 2:
				got, ok := h.PeekTop()
				if len(model) == 0 {
					if ok {
						t.Fatalf("PeekTop() ok = true on empty model")
					}
					continue
				}
				if !ok || got != model[0] {
					t.Fatalf("PeekTop() = (%d, %v), want (%d, true)", got, ok, model[0])
				}
			}
			if h.Len() != len(model) {
				t.Fatalf("Len() = %d, want %d", h.Len(), len(model))
			}
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		h := New[int](4, cmpInt)
		h.Push(3)
		h.Push(1)
		h.Push(2)
		got := collectSeq(h.Values())
		if len(got) != 3 {
			t.Fatalf("Values count = %d, want 3", len(got))
		}

		count := 0
		for range h.Values() {
			count++
			if count == 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("early-stop count = %d, want 2", count)
		}

		h.Clear()
		emptyCount := 0
		for range h.Values() {
			emptyCount++
		}
		if emptyCount != 0 {
			t.Fatalf("empty Values count = %d, want 0", emptyCount)
		}
		t.Log("Values order is internal array order; mutation during iteration is not safe by contract")
	})
}

func TestHeapCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		type item struct {
			priority int
			ref      *int
		}

		cmpItem := func(a, b item) int {
			if a.priority < b.priority {
				return -1
			}
			if a.priority > b.priority {
				return 1
			}
			return 0
		}

		a := 1
		b := 2
		c := 3

		h := New[item](4, cmpItem)
		h.Push(item{priority: 3, ref: &c})
		h.Push(item{priority: 1, ref: &a})
		h.Push(item{priority: 2, ref: &b})

		originalOrder := collectSeq(h.Values())
		cloned := h.Clone()
		if cloned == h {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != h.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), h.Len())
		}
		if cloned.Cap() != h.Cap() {
			t.Fatalf("clone Cap() = %d, want %d", cloned.Cap(), h.Cap())
		}

		clonedOrder := collectSeq(cloned.Values())
		if len(clonedOrder) != len(originalOrder) {
			t.Fatalf("clone Values() len = %d, want %d", len(clonedOrder), len(originalOrder))
		}
		for i := range originalOrder {
			if clonedOrder[i].priority != originalOrder[i].priority || clonedOrder[i].ref != originalOrder[i].ref {
				t.Fatalf("clone Values()[%d] = %+v, want %+v", i, clonedOrder[i], originalOrder[i])
			}
		}

		if _, ok := h.PopTop(); !ok {
			t.Fatalf("PopTop on original failed")
		}
		if cloned.Len() != 3 {
			t.Fatalf("clone Len after original mutation = %d, want 3", cloned.Len())
		}
	})

	t.Run("clonewith_nil_and_custom_hook", func(t *testing.T) {
		h := New[int](4, cmpInt)
		h.Push(3)
		h.Push(1)
		h.Push(2)

		originalOrder := collectSeq(h.Values())
		nilClone := h.CloneWith(nil)
		nilOrder := collectSeq(nilClone.Values())
		if len(nilOrder) != len(originalOrder) {
			t.Fatalf("CloneWith(nil) len = %d, want %d", len(nilOrder), len(originalOrder))
		}
		for i := range originalOrder {
			if nilOrder[i] != originalOrder[i] {
				t.Fatalf("CloneWith(nil) Values()[%d] = %d, want %d", i, nilOrder[i], originalOrder[i])
			}
		}

		calls := 0
		cloned := h.CloneWith(func(v int) int {
			calls++
			return v + 10
		})
		if calls != len(originalOrder) {
			t.Fatalf("cloneValue calls = %d, want %d", calls, len(originalOrder))
		}

		got := collectSeq(cloned.Values())
		if len(got) != len(originalOrder) {
			t.Fatalf("CloneWith Values() len = %d, want %d", len(got), len(originalOrder))
		}
		for i := range originalOrder {
			if got[i] != originalOrder[i]+10 {
				t.Fatalf("CloneWith Values()[%d] = %d, want %d", i, got[i], originalOrder[i]+10)
			}
		}
	})

	t.Run("clone_empty_hook_and_internal_order", func(t *testing.T) {
		empty := New[int](4, cmpInt)
		calls := 0
		clonedEmpty := empty.CloneWith(func(v int) int {
			calls++
			return v
		})
		if calls != 0 {
			t.Fatalf("empty CloneWith hook calls = %d, want 0", calls)
		}
		if clonedEmpty.Len() != 0 || clonedEmpty.Cap() != empty.Cap() {
			t.Fatalf("empty clone = (len=%d cap=%d), want (0 %d)", clonedEmpty.Len(), clonedEmpty.Cap(), empty.Cap())
		}

		h := New[int](4, cmpInt)
		h.Push(4)
		h.Push(1)
		h.Push(3)
		original := collectSeq(h.Values())
		cloned := h.Clone()
		got := collectSeq(cloned.Values())
		if len(got) != len(original) {
			t.Fatalf("clone Values len = %d, want %d", len(got), len(original))
		}
		for i := range original {
			if got[i] != original[i] {
				t.Fatalf("clone Values()[%d] = %d, want %d", i, got[i], original[i])
			}
		}
	})
}
