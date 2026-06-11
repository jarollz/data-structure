package maptreeavl

import "testing"

var benchMapTreeAvlSink int

func BenchmarkMapTreeAvlPut(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			next := n
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Put(next, next)
				next++
			}
			benchMapTreeAvlSink = next
		})
	}
}

func BenchmarkMapTreeAvlGet(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			sum := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if v, ok := m.Get(i % n); ok {
					sum += v
				}
			}
			benchMapTreeAvlSink = sum
		})
	}
}

func BenchmarkMapTreeAvlDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := i % n
				m.Delete(k)
				b.StopTimer()
				m.Put(k, k)
				b.StartTimer()
			}
		})
	}
}

func BenchmarkMapTreeAvlAll(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			sum := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for k, v := range m.All() {
					sum += k + v
				}
			}
			benchMapTreeAvlSink = sum
		})
	}
}

func BenchmarkMapTreeAvlMixedReadWrite(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			next := n
			acc := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Put(next, next)
				next++
				if v, ok := m.Get(i % n); ok {
					acc += v
				}
				m.Delete(i % n)
			}
			benchMapTreeAvlSink = acc
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
