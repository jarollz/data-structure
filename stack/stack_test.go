package stack

import (
	"math/rand"
	"testing"
)

func TestStackSpec(t *testing.T) {
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

	t.Run("push_pop_peek_clear", func(t *testing.T) {
		s := New[int](2)
		if _, ok := s.Pop(); ok {
			t.Fatalf("Pop() ok = true on empty stack")
		}
		if _, ok := s.PeekTop(); ok {
			t.Fatalf("PeekTop() ok = true on empty stack")
		}
		if !s.Push(1) || !s.Push(2) {
			t.Fatalf("Push should return true")
		}
		if v, ok := s.PeekTop(); !ok || v != 2 {
			t.Fatalf("PeekTop() = (%d, %v), want (2, true)", v, ok)
		}
		if v, ok := s.Pop(); !ok || v != 2 {
			t.Fatalf("Pop() = (%d, %v), want (2, true)", v, ok)
		}
		s.Clear()
		if s.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", s.Len())
		}
	})

	t.Run("growth_shrink_and_lifo_model", func(t *testing.T) {
		s := New[int](1)
		for i := 0; i < 64; i++ {
			if !s.Push(i) {
				t.Fatalf("Push(%d) = false, want true", i)
			}
		}
		if s.Len() != 64 {
			t.Fatalf("Len() after pushes = %d, want 64", s.Len())
		}
		if s.Cap() < 64 {
			t.Fatalf("Cap() after pushes = %d, want >= 64", s.Cap())
		}
		for i := 63; i >= 0; i-- {
			v, ok := s.Pop()
			if !ok || v != i {
				t.Fatalf("Pop() = (%d, %v), want (%d, true)", v, ok, i)
			}
		}
		if s.Len() != 0 {
			t.Fatalf("Len() after draining = %d, want 0", s.Len())
		}
		if s.Cap() < 16 {
			t.Fatalf("Cap() after draining = %d, want >= 16", s.Cap())
		}
	})

	t.Run("randomized_against_slice_model", func(t *testing.T) {
		stack := New[int](4)
		model := make([]int, 0)
		rng := rand.New(rand.NewSource(1))
		for step := 0; step < 500; step++ {
			switch rng.Intn(3) {
			case 0:
				value := step + 1
				stack.Push(value)
				model = append(model, value)
			case 1:
				got, ok := stack.Pop()
				if len(model) == 0 {
					if ok {
						t.Fatalf("Pop() ok = true on empty model")
					}
					continue
				}
				want := model[len(model)-1]
				model = model[:len(model)-1]
				if !ok || got != want {
					t.Fatalf("Pop() = (%d, %v), want (%d, true)", got, ok, want)
				}
			case 2:
				got, ok := stack.PeekTop()
				if len(model) == 0 {
					if ok {
						t.Fatalf("PeekTop() ok = true on empty model")
					}
					continue
				}
				want := model[len(model)-1]
				if !ok || got != want {
					t.Fatalf("PeekTop() = (%d, %v), want (%d, true)", got, ok, want)
				}
			}
			if stack.Len() != len(model) {
				t.Fatalf("Len() = %d, want %d", stack.Len(), len(model))
			}
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		s := New[int](4)
		s.Push(1)
		s.Push(2)
		s.Push(3)
		got := collectSeq(s.Values())
		if len(got) != 3 || got[0] != 3 || got[1] != 2 || got[2] != 1 {
			t.Fatalf("Values order/count = %v, want [3 2 1]", got)
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

func TestStackCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		type node struct{ value int }

		a := &node{value: 1}
		b := &node{value: 2}
		c := &node{value: 3}

		s := New[*node](4)
		s.Push(a)
		s.Push(b)
		s.Push(c)

		cloned := s.Clone()
		if cloned == s {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != s.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), s.Len())
		}
		if cloned.Cap() != s.Cap() {
			t.Fatalf("clone Cap() = %d, want %d", cloned.Cap(), s.Cap())
		}

		got := collectSeq(cloned.Values())
		if len(got) != 3 || got[0] != c || got[1] != b || got[2] != a {
			t.Fatalf("clone Values() = %v, want [%p %p %p]", got, c, b, a)
		}

		if v, ok := s.Pop(); !ok || v != c {
			t.Fatalf("original Pop() = (%p, %v), want (%p, true)", v, ok, c)
		}
		if cloned.Len() != 3 {
			t.Fatalf("clone Len after original mutation = %d, want 3", cloned.Len())
		}

		d := &node{value: 4}
		if !cloned.Push(d) {
			t.Fatalf("Push on clone failed")
		}
		if s.Len() != 2 {
			t.Fatalf("original Len after clone mutation = %d, want 2", s.Len())
		}
	})

		t.Run("clonewith_nil_and_custom_hook", func(t *testing.T) {
		type item struct{ value int }

		s := New[item](4)
		s.Push(item{value: 1})
		s.Push(item{value: 2})
		s.Push(item{value: 3})

		nilClone := s.CloneWith(nil)
		nilValues := collectSeq(nilClone.Values())
		if len(nilValues) != 3 || nilValues[0].value != 3 || nilValues[1].value != 2 || nilValues[2].value != 1 {
			t.Fatalf("CloneWith(nil) Values() = %v, want [3 2 1]", nilValues)
		}

		calls := 0
		cloned := s.CloneWith(func(v item) item {
			calls++
			return item{value: v.value * 10}
		})
		if calls != 3 {
			t.Fatalf("cloneValue calls = %d, want 3", calls)
		}

		got := collectSeq(cloned.Values())
			if len(got) != 3 || got[0].value != 30 || got[1].value != 20 || got[2].value != 10 {
				t.Fatalf("CloneWith Values() = %v, want [30 20 10]", got)
			}
		})

		t.Run("clone_empty_and_hook_order", func(t *testing.T) {
			empty := New[int](2)
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

			stack := New[int](4)
			stack.Push(1)
			stack.Push(2)
			stack.Push(3)
			order := make([]int, 0, 3)
			stack.CloneWith(func(v int) int {
				order = append(order, v)
				return v
			})
			if len(order) != 3 || order[0] != 3 || order[1] != 2 || order[2] != 1 {
				t.Fatalf("CloneWith hook order = %v, want [3 2 1]", order)
			}
		})
}
