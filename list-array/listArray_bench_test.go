package listarray

import "testing"

var benchListArraySink int

func BenchmarkListArrayAppend(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				s := New[int](n)
				s.Append(i)
			}
		})
	}
}

func BenchmarkListArrayMiddleInsertDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](n)
			for i := 0; i < n; i++ {
				s.Append(i)
			}
			mid := n / 2
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Insert(mid, i)
				s.Delete(mid)
			}
		})
	}
}

func BenchmarkListArrayValues(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](n)
			for i := 0; i < n; i++ {
				s.Append(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			sum := 0
			for i := 0; i < b.N; i++ {
				for v := range s.Values() {
					sum += v
				}
			}
			benchListArraySink = sum
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
