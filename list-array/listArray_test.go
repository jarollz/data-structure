package listarray

import (
	"math/rand"
	"testing"
)

func TestListArraySpec(t *testing.T) {
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

		t.Run("append_get_set_insert_delete", func(t *testing.T) {
		s := New[int](2)
		if !s.Append(10) || !s.Append(20) {
			t.Fatalf("Append should return true")
		}
		if v, ok := s.Get(1); !ok || v != 20 {
			t.Fatalf("Get(1) = (%d, %v), want (20, true)", v, ok)
		}
		if _, ok := s.Get(-1); ok {
			t.Fatalf("Get(-1) ok = true, want false")
		}
		if !s.Set(1, 30) {
			t.Fatalf("Set(1, 30) = false, want true")
		}
		if s.Set(5, 99) {
			t.Fatalf("Set(5, 99) = true, want false")
		}
		if !s.Insert(0, 5) {
			t.Fatalf("Insert at head should succeed")
		}
		if !s.Insert(s.Len(), 40) {
			t.Fatalf("Insert at tail should succeed")
		}
		if s.Insert(-1, 1) {
			t.Fatalf("Insert(-1, 1) = true, want false")
		}
		if v, ok := s.Delete(0); !ok || v != 5 {
			t.Fatalf("Delete(0) = (%d, %v), want (5, true)", v, ok)
		}
		if _, ok := s.Delete(99); ok {
			t.Fatalf("Delete(99) ok = true, want false")
		}
	})

	t.Run("insert_delete_shift_and_resize_paths", func(t *testing.T) {
		list := New[int](1)
		for i := 0; i < 32; i++ {
			list.Append(i)
		}
		if !list.Insert(16, 999) {
			t.Fatalf("Insert(16, 999) = false, want true")
		}
		if v, ok := list.Get(16); !ok || v != 999 {
			t.Fatalf("Get(16) = (%d, %v), want (999, true)", v, ok)
		}
		if v, ok := list.Get(17); !ok || v != 16 {
			t.Fatalf("Get(17) = (%d, %v), want (16, true)", v, ok)
		}
		deleted, ok := list.Delete(16)
		if !ok || deleted != 999 {
			t.Fatalf("Delete(16) = (%d, %v), want (999, true)", deleted, ok)
		}
		if v, ok := list.Get(16); !ok || v != 16 {
			t.Fatalf("Get(16) after delete = (%d, %v), want (16, true)", v, ok)
		}
		for list.Len() > 0 {
			list.Delete(list.Len() - 1)
		}
		if list.Cap() < 16 {
			t.Fatalf("Cap() after shrink path = %d, want >= 16", list.Cap())
		}
	})

	t.Run("randomized_against_slice_model", func(t *testing.T) {
		list := New[int](4)
		model := make([]int, 0)
		rng := rand.New(rand.NewSource(4))
		for step := 0; step < 500; step++ {
			switch rng.Intn(5) {
			case 0:
				value := step + 1
				list.Append(value)
				model = append(model, value)
			case 1:
				if len(model) == 0 {
					if _, ok := list.Delete(0); ok {
						t.Fatalf("Delete(0) ok = true on empty model")
					}
					continue
				}
				idx := rng.Intn(len(model))
				got, ok := list.Delete(idx)
				want := model[idx]
				model = append(model[:idx], model[idx+1:]...)
				if !ok || got != want {
					t.Fatalf("Delete(%d) = (%d, %v), want (%d, true)", idx, got, ok, want)
				}
			case 2:
				idx := 0
				if len(model) > 0 {
					idx = rng.Intn(len(model) + 1)
				}
				value := step + 1000
				if !list.Insert(idx, value) {
					t.Fatalf("Insert(%d, %d) = false, want true", idx, value)
				}
				model = append(model, 0)
				copy(model[idx+1:], model[idx:])
				model[idx] = value
			case 3:
				if len(model) == 0 {
					continue
				}
				idx := rng.Intn(len(model))
				value := step + 2000
				if !list.Set(idx, value) {
					t.Fatalf("Set(%d, %d) = false, want true", idx, value)
				}
				model[idx] = value
			case 4:
				if len(model) == 0 {
					if _, ok := list.Get(0); ok {
						t.Fatalf("Get(0) ok = true on empty model")
					}
					continue
				}
				idx := rng.Intn(len(model))
				got, ok := list.Get(idx)
				if !ok || got != model[idx] {
					t.Fatalf("Get(%d) = (%d, %v), want (%d, true)", idx, got, ok, model[idx])
				}
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
		s := New[int](4)
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

func TestListArrayCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		type node struct{ value int }

		a := &node{value: 1}
		b := &node{value: 2}

		list := New[*node](4)
		list.Append(a)
		list.Append(b)

		cloned := list.Clone()
		if cloned == list {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != list.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), list.Len())
		}
		if cloned.Cap() != list.Cap() {
			t.Fatalf("clone Cap() = %d, want %d", cloned.Cap(), list.Cap())
		}

		got := collectSeq(cloned.Values())
		if len(got) != 2 || got[0] != a || got[1] != b {
			t.Fatalf("clone Values() = %v, want [%p %p]", got, a, b)
		}

		if _, ok := list.Delete(0); !ok {
			t.Fatalf("Delete(0) on original failed")
		}
		if cloned.Len() != 2 {
			t.Fatalf("clone Len after original mutation = %d, want 2", cloned.Len())
		}

		replacement := &node{value: 3}
		if !cloned.Set(1, replacement) {
			t.Fatalf("Set(1, replacement) on clone failed")
		}
		if v, ok := list.Get(0); !ok || v != b {
			t.Fatalf("original Get(0) after clone mutation = (%p, %v), want (%p, true)", v, ok, b)
		}
	})

	t.Run("clonewith_nil_and_custom_hook", func(t *testing.T) {
		type item struct{ value int }

		list := New[item](4)
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

		if !list.Set(0, item{value: 99}) {
			t.Fatalf("Set on original failed")
		}
		if got := collectSeq(cloned.Values()); got[0].value != 10 {
			t.Fatalf("clone changed after original mutation: %v", got)
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

		list := New[int](4)
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
