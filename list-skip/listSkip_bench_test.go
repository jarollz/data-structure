package listskip

import "testing"

var benchListSkipSink int

func BenchmarkListSkipSearch(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](16, cmpInt)
			for i := 0; i < n; i++ {
				s.Insert(i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if s.Has(i % n) {
					hits++
				}
			}
			benchListSkipSink = hits
		})
	}
}

func BenchmarkListSkipInsert(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](16, cmpInt)
			for i := 0; i < n; i++ {
				s.Insert(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Insert(n + i)
			}
		})
	}
}

func BenchmarkListSkipDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](16, cmpInt)
			for i := 0; i < n; i++ {
				s.Insert(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := i % n
				s.Delete(k)
				b.StopTimer()
				s.Insert(k)
				b.StartTimer()
			}
		})
	}
}

func BenchmarkListSkipMixedOrdered(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](16, cmpInt)
			for i := 0; i < n; i++ {
				s.Insert(i)
			}
			next := n
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Has(i % n)
				s.Insert(next)
				next++
				s.Delete(i % n)
			}
			benchListSkipSink = next
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
