package treeavl

import "testing"

var benchTreeAvlSink int

func BenchmarkTreeAvlInsert(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](cmpInt)
			for i := 0; i < n; i++ {
				t.Insert(i)
			}
			next := n
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				t.Insert(next)
				next++
			}
			benchTreeAvlSink = next
		})
	}
}

func BenchmarkTreeAvlHas(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](cmpInt)
			for i := 0; i < n; i++ {
				t.Insert(i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if t.Has(i % n) {
					hits++
				}
			}
			benchTreeAvlSink = hits
		})
	}
}

func BenchmarkTreeAvlDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](cmpInt)
			for i := 0; i < n; i++ {
				t.Insert(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := i % n
				t.Delete(k)
				b.StopTimer()
				t.Insert(k)
				b.StartTimer()
			}
		})
	}
}

func BenchmarkTreeAvlInOrder(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](cmpInt)
			for i := 0; i < n; i++ {
				t.Insert(i)
			}
			sum := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for v := range t.InOrder() {
					sum += v
				}
			}
			benchTreeAvlSink = sum
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
