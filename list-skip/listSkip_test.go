package listskip

import (
	"math/rand"
	"sort"
	"testing"
)

func TestListSkipSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		s := New[int](0, cmpInt)
		if s == nil {
			t.Fatalf("New(0, cmp) = nil, want non-nil")
		}
		if s.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", s.Len())
		}
	})

	t.Run("insert_delete_has_clear", func(t *testing.T) {
		s := New[int](4, cmpInt)
		if !s.Insert(2) || !s.Insert(1) || !s.Insert(3) {
			t.Fatalf("Insert new values should return true")
		}
		if s.Insert(2) {
			t.Fatalf("duplicate Insert should return false")
		}
		if !s.Has(1) {
			t.Fatalf("Has(1) = false, want true")
		}
		if !s.Delete(2) {
			t.Fatalf("Delete(2) = false, want true")
		}
		if s.Delete(9) {
			t.Fatalf("Delete(9) = true, want false")
		}
		s.Clear()
		if s.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", s.Len())
		}
	})

	t.Run("duplicate_missing_and_sorted_populated_state", func(t *testing.T) {
		list := New[int](4, cmpInt)
		for _, value := range []int{5, 1, 3, 2, 4} {
			if !list.Insert(value) {
				t.Fatalf("Insert(%d) = false, want true", value)
			}
		}
		if list.Insert(3) {
			t.Fatalf("Insert(3) = true on duplicate, want false")
		}
		if list.Delete(99) {
			t.Fatalf("Delete(99) = true, want false")
		}
		got := collectSeq(list.Values())
		for i := 0; i < len(got); i++ {
			if got[i] != i+1 {
				t.Fatalf("Values()[%d] = %d, want %d", i, got[i], i+1)
			}
		}
	})

	t.Run("randomized_against_sorted_model", func(t *testing.T) {
		list := New[int](8, cmpInt)
		model := make([]int, 0)
		rng := rand.New(rand.NewSource(7))
		for step := 0; step < 400; step++ {
			value := rng.Intn(200)
			switch rng.Intn(3) {
			case 0:
				inserted := list.Insert(value)
				found := false
				for _, existing := range model {
					if existing == value {
						found = true
						break
					}
				}
				if inserted == found {
					t.Fatalf("Insert(%d) returned %v, duplicate found=%v", value, inserted, found)
				}
				if !found {
					model = append(model, value)
					sort.Ints(model)
				}
			case 1:
				deleted := list.Delete(value)
				idx := -1
				for i, existing := range model {
					if existing == value {
						idx = i
						break
					}
				}
				if deleted != (idx >= 0) {
					t.Fatalf("Delete(%d) returned %v, model index=%d", value, deleted, idx)
				}
				if idx >= 0 {
					model = append(model[:idx], model[idx+1:]...)
				}
			case 2:
				has := list.Has(value)
				found := false
				for _, existing := range model {
					if existing == value {
						found = true
						break
					}
				}
				if has != found {
					t.Fatalf("Has(%d) = %v, want %v", value, has, found)
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
		s := New[int](6, cmpInt)
		s.Insert(3)
		s.Insert(1)
		s.Insert(2)
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

func TestListSkipCloneSpec(t *testing.T) {
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
		list := New[item](4, cmpItem)
		list.Insert(item{key: 1, ref: &a})
		list.Insert(item{key: 2, ref: &b})

		cloned := list.Clone()
		if cloned == list {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != list.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), list.Len())
		}

		got := collectSeq(cloned.Values())
		if len(got) != 2 || got[0].ref != &a || got[1].ref != &b {
			t.Fatalf("clone Values() = %v, want shared refs", got)
		}

		if !list.Delete(item{key: 1}) {
			t.Fatalf("Delete on original failed")
		}
		if cloned.Len() != 2 {
			t.Fatalf("clone Len after original mutation = %d, want 2", cloned.Len())
		}

		if !cloned.Has(item{key: 1}) || !cloned.Has(item{key: 2}) {
			t.Fatalf("clone lost values after original mutation")
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

		list := New[item](4, cmpItem)
		list.Insert(item{key: 1})
		list.Insert(item{key: 2})
		list.Insert(item{key: 3})

		nilClone := list.CloneWith(nil)
		nilValues := collectSeq(nilClone.Values())
		if len(nilValues) != 3 || nilValues[0].key != 1 || nilValues[1].key != 2 || nilValues[2].key != 3 {
			t.Fatalf("CloneWith(nil) Values() = %v, want [1 2 3]", nilValues)
		}

		calls := 0
		cloned := list.CloneWith(func(v item) item {
			calls++
			return item{key: v.key + 10}
		})
		if calls != 3 {
			t.Fatalf("cloneValue calls = %d, want 3", calls)
		}

		got := collectSeq(cloned.Values())
		if len(got) != 3 || got[0].key != 11 || got[1].key != 12 || got[2].key != 13 {
			t.Fatalf("CloneWith Values() = %v, want [11 12 13]", got)
		}
	})

	t.Run("clone_empty_and_hook_order", func(t *testing.T) {
		empty := New[int](4, cmpInt)
		calls := 0
		empty.CloneWith(func(v int) int {
			calls++
			return v
		})
		if calls != 0 {
			t.Fatalf("empty CloneWith hook calls = %d, want 0", calls)
		}

		list := New[int](4, cmpInt)
		list.Insert(3)
		list.Insert(1)
		list.Insert(2)
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
