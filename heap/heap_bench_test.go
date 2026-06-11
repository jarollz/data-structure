package heap

import "testing"

var benchHeapSinkInt int
var benchHeapSinkLarge benchLargePayload

func BenchmarkHeapPush(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("tiny_size="+itoa(n), func(b *testing.B) {
			h := New[int](n, cmpInt)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Push(next)
				next++
				if h.Len() > n*2 {
					b.StopTimer()
					for h.Len() > n {
						h.PopTop()
					}
					b.StartTimer()
				}
			}
			benchHeapSinkInt = next
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkHeapPopTop(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			h := New[int](n, cmpInt)
			for i := 0; i < n; i++ {
				h.Push(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := h.PopTop()
				benchHeapSinkInt = v
				if h.Len() == 0 {
					b.StopTimer()
					for j := 0; j < n; j++ {
						h.Push(j)
					}
					b.StartTimer()
				}
			}
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkHeapPeekTop(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			h := New[int](n, cmpInt)
			for i := 0; i < n; i++ {
				h.Push(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := h.PeekTop()
				benchHeapSinkInt = v
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkHeapClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h := New[int](n, cmpInt)
				for j := 0; j < n; j++ {
					h.Push(j)
				}
				h.Clear()
				benchHeapSinkInt = h.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkHeapClone(b *testing.B) {
	cmpLarge := func(a, c benchLargePayload) int {
		if a.Data[0] < c.Data[0] {
			return -1
		}
		if a.Data[0] > c.Data[0] {
			return 1
		}
		return 0
	}
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			h := New[benchLargePayload](n, cmpLarge)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				h.Push(value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := h.Clone()
				v, _ := cloned.PeekTop()
				benchHeapSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkHeapCloneWith(b *testing.B) {
	cmpLarge := func(a, c benchLargePayload) int {
		if a.Data[0] < c.Data[0] {
			return -1
		}
		if a.Data[0] > c.Data[0] {
			return 1
		}
		return 0
	}
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			h := New[benchLargePayload](n, cmpLarge)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				h.Push(value)
			}
			cloneValue := func(v benchLargePayload) benchLargePayload {
				v.Data[1]++
				return v
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := h.CloneWith(cloneValue)
				v, _ := cloned.PeekTop()
				benchHeapSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkHeapValues(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			h := New[benchLargePayload](n, func(a, c benchLargePayload) int {
				if a.Data[0] < c.Data[0] {
					return -1
				}
				if a.Data[0] > c.Data[0] {
					return 1
				}
				return 0
			})
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				h.Push(value)
			}
			var sum uint64
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for value := range h.Values() {
					sum += value.Data[0]
				}
			}
			benchHeapSinkLarge.Data[0] = sum
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkHeapMixedPriorityQueue(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			h := New[int](n, cmpInt)
			for i := 0; i < n; i++ {
				h.Push(i)
			}
			next := n
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Push(next)
				next++
				v, _ := h.PopTop()
				benchHeapSinkInt = v
			}
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
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
