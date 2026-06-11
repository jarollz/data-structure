package queue

import "testing"

var benchQueueSink int

func BenchmarkQueueEnqueue(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			q := New[int](n)
			for i := 0; i < n; i++ {
				q.Enqueue(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				q.Enqueue(i)
			}
		})
	}
}

func BenchmarkQueueDequeue(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			q := New[int](n)
			for i := 0; i < n; i++ {
				q.Enqueue(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := q.Dequeue()
				benchQueueSink = v
				b.StopTimer()
				q.Enqueue(i)
				b.StartTimer()
			}
		})
	}
}

func BenchmarkQueueMixedEnqueueDequeue(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			q := New[int](n)
			for i := 0; i < n; i++ {
				q.Enqueue(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				q.Enqueue(i)
				v, _ := q.Dequeue()
				benchQueueSink = v
			}
		})
	}
}

func itoa(v int) string {
	if v == 1_000 {
		return "1e3"
	}
	if v == 10_000 {
		return "1e4"
	}
	return "1e5"
}
