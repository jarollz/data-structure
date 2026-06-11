package queue

import (
	"math/rand"
	"testing"
)

func TestQueueSpec(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		q := New[int](0)
		if q == nil {
			t.Fatalf("New(0) = nil, want non-nil")
		}
		if q.Cap() < 16 {
			t.Fatalf("Cap() = %d, want >= 16", q.Cap())
		}
		if q.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", q.Len())
		}
	})

	t.Run("enqueue_dequeue_peek_clear", func(t *testing.T) {
		q := New[int](2)
		if _, ok := q.Dequeue(); ok {
			t.Fatalf("Dequeue() ok = true on empty queue")
		}
		if _, ok := q.PeekFront(); ok {
			t.Fatalf("PeekFront() ok = true on empty queue")
		}
		if !q.Enqueue(1) || !q.Enqueue(2) || !q.Enqueue(3) {
			t.Fatalf("Enqueue should return true")
		}
		if v, ok := q.PeekFront(); !ok || v != 1 {
			t.Fatalf("PeekFront() = (%d, %v), want (1, true)", v, ok)
		}
		if v, ok := q.Dequeue(); !ok || v != 1 {
			t.Fatalf("Dequeue() = (%d, %v), want (1, true)", v, ok)
		}
		q.Clear()
		if q.Len() != 0 {
			t.Fatalf("Len after Clear = %d, want 0", q.Len())
		}
	})

	t.Run("wraparound_growth_and_fifo_model", func(t *testing.T) {
		q := New[int](2)
		for i := 0; i < 32; i++ {
			q.Enqueue(i)
			if i%3 == 0 {
				v, ok := q.Dequeue()
				if !ok || v != i/3*3 {
					t.Fatalf("Dequeue() = (%d, %v), want FIFO value %d", v, ok, i/3*3)
				}
			}
		}
		prev := -1
		for q.Len() > 0 {
			v, ok := q.Dequeue()
			if !ok {
				t.Fatalf("Dequeue() failed while queue non-empty")
			}
			if prev >= v {
				t.Fatalf("dequeue order not FIFO ascending tail snapshot: prev=%d now=%d", prev, v)
			}
			prev = v
		}
	})

	t.Run("randomized_against_slice_model", func(t *testing.T) {
		q := New[int](4)
		model := make([]int, 0)
		rng := rand.New(rand.NewSource(2))
		for step := 0; step < 500; step++ {
			switch rng.Intn(3) {
			case 0:
				value := step + 1
				q.Enqueue(value)
				model = append(model, value)
			case 1:
				got, ok := q.Dequeue()
				if len(model) == 0 {
					if ok {
						t.Fatalf("Dequeue() ok = true on empty model")
					}
					continue
				}
				want := model[0]
				model = model[1:]
				if !ok || got != want {
					t.Fatalf("Dequeue() = (%d, %v), want (%d, true)", got, ok, want)
				}
			case 2:
				got, ok := q.PeekFront()
				if len(model) == 0 {
					if ok {
						t.Fatalf("PeekFront() ok = true on empty model")
					}
					continue
				}
				if !ok || got != model[0] {
					t.Fatalf("PeekFront() = (%d, %v), want (%d, true)", got, ok, model[0])
				}
			}
			if q.Len() != len(model) {
				t.Fatalf("Len() = %d, want %d", q.Len(), len(model))
			}
		}
	})

	t.Run("iterator_contract", func(t *testing.T) {
		q := New[int](4)
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		got := collectSeq(q.Values())
		if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
			t.Fatalf("Values order/count = %v, want [1 2 3]", got)
		}

		count := 0
		for range q.Values() {
			count++
			if count == 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("early-stop count = %d, want 2", count)
		}

		q.Clear()
		emptyCount := 0
		for range q.Values() {
			emptyCount++
		}
		if emptyCount != 0 {
			t.Fatalf("empty Values count = %d, want 0", emptyCount)
		}
		t.Log("mutation during iteration is not safe by contract")
	})
}

func TestQueueCloneSpec(t *testing.T) {
	t.Run("clone_assignment_copy_and_independence", func(t *testing.T) {
		type node struct{ value int }

		a := &node{value: 1}
		b := &node{value: 2}
		c := &node{value: 3}

		q := New[*node](4)
		q.Enqueue(a)
		q.Enqueue(b)
		q.Enqueue(c)

		cloned := q.Clone()
		if cloned == q {
			t.Fatalf("Clone() returned original pointer")
		}
		if cloned.Len() != q.Len() {
			t.Fatalf("clone Len() = %d, want %d", cloned.Len(), q.Len())
		}
		if cloned.Cap() != q.Cap() {
			t.Fatalf("clone Cap() = %d, want %d", cloned.Cap(), q.Cap())
		}

		got := collectSeq(cloned.Values())
		if len(got) != 3 || got[0] != a || got[1] != b || got[2] != c {
			t.Fatalf("clone Values() = %v, want [%p %p %p]", got, a, b, c)
		}

		if v, ok := q.Dequeue(); !ok || v != a {
			t.Fatalf("original Dequeue() = (%p, %v), want (%p, true)", v, ok, a)
		}
		if cloned.Len() != 3 {
			t.Fatalf("clone Len after original mutation = %d, want 3", cloned.Len())
		}

		d := &node{value: 4}
		if !cloned.Enqueue(d) {
			t.Fatalf("Enqueue on clone failed")
		}
		if q.Len() != 2 {
			t.Fatalf("original Len after clone mutation = %d, want 2", q.Len())
		}
	})

	t.Run("clonewith_nil_and_custom_hook", func(t *testing.T) {
		type item struct{ value int }

		q := New[item](4)
		q.Enqueue(item{value: 1})
		q.Enqueue(item{value: 2})
		q.Enqueue(item{value: 3})

		nilClone := q.CloneWith(nil)
		nilValues := collectSeq(nilClone.Values())
		if len(nilValues) != 3 || nilValues[0].value != 1 || nilValues[1].value != 2 || nilValues[2].value != 3 {
			t.Fatalf("CloneWith(nil) Values() = %v, want [1 2 3]", nilValues)
		}

		calls := 0
		cloned := q.CloneWith(func(v item) item {
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

	t.Run("clone_wraparound_order_and_empty_hook", func(t *testing.T) {
		q := New[int](2)
		q.Enqueue(1)
		q.Enqueue(2)
		q.Dequeue()
		q.Enqueue(3)
		cloned := q.Clone()
		got := collectSeq(cloned.Values())
		if len(got) != 2 || got[0] != 2 || got[1] != 3 {
			t.Fatalf("wraparound clone Values() = %v, want [2 3]", got)
		}

		calls := 0
		empty := New[int](2)
		empty.CloneWith(func(v int) int {
			calls++
			return v
		})
		if calls != 0 {
			t.Fatalf("empty CloneWith hook calls = %d, want 0", calls)
		}
	})
}
