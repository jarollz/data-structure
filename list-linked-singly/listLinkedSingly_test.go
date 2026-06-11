package listlinkedsingly

import (
	"math/rand"
	"testing"
)

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

	t.Run("head_tail_transitions_and_append_contract", func(t *testing.T) {
		list := New[int]()
		list.Append(1)
		list.Append(2)
		list.PushFront(0)
		if got := collectSeq(list.Values()); len(got) != 3 || got[0] != 0 || got[1] != 1 || got[2] != 2 {
			t.Fatalf("Values() = %v, want [0 1 2]", got)
		}
		if v, ok := list.PopFront(); !ok || v != 0 {
			t.Fatalf("PopFront() = (%d, %v), want (0, true)", v, ok)
		}
		if v, ok := list.PopFront(); !ok || v != 1 {
			t.Fatalf("PopFront() = (%d, %v), want (1, true)", v, ok)
		}
		if v, ok := list.PopFront(); !ok || v != 2 {
			t.Fatalf("PopFront() = (%d, %v), want (2, true)", v, ok)
		}
		if _, ok := list.PopFront(); ok {
			t.Fatalf("PopFront() ok = true on empty list after full drain")
		}
	})

	t.Run("randomized_against_slice_model", func(t *testing.T) {
		list := New[int]()
		model := make([]int, 0)
		rng := rand.New(rand.NewSource(5))
		for step := 0; step < 400; step++ {
			switch rng.Intn(4) {
			case 0:
				value := step + 1
				list.PushFront(value)
				model = append([]int{value}, model...)
			case 1:
				value := step + 1000
				list.Append(value)
				model = append(model, value)
			case 2:
				got, ok := list.PopFront()
				if len(model) == 0 {
					if ok {
						t.Fatalf("PopFront() ok = true on empty model")
					}
					continue
				}
				want := model[0]
				model = model[1:]
				if !ok || got != want {
					t.Fatalf("PopFront() = (%d, %v), want (%d, true)", got, ok, want)
				}
			case 3:
				if len(model) == 0 {
					if list.DeleteFirst(func(v int) bool { return true }) {
						t.Fatalf("DeleteFirst on empty model = true, want false")
					}
					continue
				}
				idx := rng.Intn(len(model))
				want := model[idx]
				if !list.DeleteFirst(func(v int) bool { return v == want }) {
					t.Fatalf("DeleteFirst(%d) = false, want true", want)
				}
				model = append(model[:idx], model[idx+1:]...)
			}
			if list.Len() != len(model) {
				t.Fatalf("Len() = %d, want %d", list.Len(), len(model))
			}
			got := collectSeq(list.Values())
			if len(got) != len(model) {
				t.Fatalf("Values len = %d, want %d", len(got), len(model))
			}
			for i := range model {
				if got[i] != model[i] {
					t.Fatalf("Values()[%d] = %d, want %d", i, got[i], model[i])
				}
			}
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

func TestListLinkedSinglyCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		type node struct{ value int }

		a := &node{value: 1}
		b := &node{value: 2}
		c := &node{value: 3}

		list := New[*node]()
		list.Append(a)
		list.Append(b)
		list.Append(c)

		cloned := list.Clone()
		if cloned == list {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != list.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), list.Len())
		}

		got := collectSeq(cloned.Values())
		if len(got) != 3 || got[0] != a || got[1] != b || got[2] != c {
			t.Fatalf("clone Values() = %v, want [%p %p %p]", got, a, b, c)
		}

		if !list.DeleteFirst(func(v *node) bool { return v == b }) {
			t.Fatalf("DeleteFirst on original failed")
		}
		if cloned.Len() != 3 {
			t.Fatalf("clone Len after original mutation = %d, want 3", cloned.Len())
		}

		if v, ok := cloned.PopFront(); !ok || v != a {
			t.Fatalf("clone PopFront() = (%p, %v), want (%p, true)", v, ok, a)
		}
		if original := collectSeq(list.Values()); len(original) != 2 || original[0] != a || original[1] != c {
			t.Fatalf("original Values() after clone mutation = %v, want [%p %p]", original, a, c)
		}
	})

	t.Run("clonewith_nil_and_custom_hook", func(t *testing.T) {
		type item struct{ value int }

		list := New[item]()
		list.Append(item{value: 1})
		list.Append(item{value: 2})
		list.Append(item{value: 3})

		nilClone := list.CloneWith(nil)
		nilValues := collectSeq(nilClone.Values())
		if len(nilValues) != 3 || nilValues[0].value != 1 || nilValues[1].value != 2 || nilValues[2].value != 3 {
			t.Fatalf("CloneWith(nil) Values() = %v, want [1 2 3]", nilValues)
		}

		calls := 0
		cloned := list.CloneWith(func(v item) item {
			calls++
			return item{value: v.value * 10}
		})
		if calls != 3 {
			t.Fatalf("cloneValue calls = %d, want 3", calls)
		}

		got := collectSeq(cloned.Values())
		if len(got) != 3 || got[0].value != 10 || got[1].value != 20 || got[2].value != 30 {
			t.Fatalf("CloneWith Values() = %v, want [10 20 30]", got)
		}
	})

	t.Run("clone_empty_and_hook_order", func(t *testing.T) {
		empty := New[int]()
		calls := 0
		empty.CloneWith(func(v int) int {
			calls++
			return v
		})
		if calls != 0 {
			t.Fatalf("empty CloneWith hook calls = %d, want 0", calls)
		}

		list := New[int]()
		list.Append(1)
		list.Append(2)
		list.Append(3)
		order := make([]int, 0, 3)
		list.CloneWith(func(v int) int {
			order = append(order, v)
			return v
		})
		if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
			t.Fatalf("CloneWith hook order = %v, want [1 2 3]", order)
		}
	})
}
