package queue

import "testing"

var benchQueueSinkInt int
var benchQueueSinkLarge benchLargePayload

func BenchmarkQueueEnqueue(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("tiny_size="+itoa(n), func(b *testing.B) {
			q := New[int](n)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				q.Enqueue(next)
				next++
				if q.Len() > n*2 {
					b.StopTimer()
					for q.Len() > n {
						q.Dequeue()
					}
					b.StartTimer()
				}
			}
			benchQueueSinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			q := New[benchLargePayload](n)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(next)
				q.Enqueue(value)
				next++
				if q.Len() > n*2 {
					b.StopTimer()
					for q.Len() > n {
						q.Dequeue()
					}
					b.StartTimer()
				}
			}
			benchQueueSinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[benchLargePayload](), n)
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
				benchQueueSinkInt = v
				if q.Len() == 0 {
					b.StopTimer()
					for j := 0; j < n; j++ {
						q.Enqueue(j)
					}
					b.StartTimer()
				}
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkQueuePeekFront(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			q := New[int](n)
			for i := 0; i < n; i++ {
				q.Enqueue(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := q.PeekFront()
				benchQueueSinkInt = v
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkQueueClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				q := New[int](n)
				for j := 0; j < n; j++ {
					q.Enqueue(j)
				}
				q.Clear()
				benchQueueSinkInt = q.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkQueueClone(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			q := New[benchLargePayload](n)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				q.Enqueue(value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := q.Clone()
				v, _ := cloned.PeekFront()
				benchQueueSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkQueueCloneWith(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			q := New[benchLargePayload](n)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				q.Enqueue(value)
			}
			cloneValue := func(v benchLargePayload) benchLargePayload {
				v.Data[1]++
				return v
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := q.CloneWith(cloneValue)
				v, _ := cloned.PeekFront()
				benchQueueSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkQueueValues(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			q := New[benchLargePayload](n)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				q.Enqueue(value)
			}
			var sum uint64
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for value := range q.Values() {
					sum += value.Data[0]
				}
			}
			benchQueueSinkLarge.Data[0] = sum
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
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
			next := n
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				q.Enqueue(next)
				next++
				v, _ := q.Dequeue()
				benchQueueSinkInt = v
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
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
