package queue

import "testing"

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
