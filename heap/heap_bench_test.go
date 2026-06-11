package heap

import "testing"

var benchHeapSink int

func BenchmarkHeapPushOnly(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			h := New[int](n, cmpInt)
			for i := 0; i < n; i++ {
				h.Push(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Push(i)
			}
		})
	}
}

func BenchmarkHeapPopOnlyAfterPrefill(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := n
			if b.N > m {
				m = b.N
			}
			h := New[int](m, cmpInt)
			for i := 0; i < m; i++ {
				h.Push(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := h.PopTop()
				benchHeapSink = v
			}
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
				benchHeapSink = v
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
