package maphash

import "testing"

var benchMapHashSinkInt int
var benchMapHashSinkFloat float64

func BenchmarkMapHashPutUnique(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](n, hashInt, eqInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Put(n+i, i)
			}
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
		})
	}
}

func BenchmarkMapHashGetHitHeavy(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](n, hashInt, eqInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, ok := m.Get(i % n); ok {
					hits++
				}
			}
			benchMapHashSinkInt = hits
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
		})
	}
}

func BenchmarkMapHashGetMissHeavy(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](n, hashInt, eqInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			miss := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, ok := m.Get(n + i); !ok {
					miss++
				}
			}
			benchMapHashSinkInt = miss
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
		})
	}
}

func BenchmarkMapHashDeleteMixed(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](n, hashInt, eqInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			deleted := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := i % n
				if m.Delete(k) {
					deleted++
				}
				m.Put(k, i)
			}
			benchMapHashSinkInt = deleted
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
		})
	}
}

func BenchmarkMapHashMixedPutGetDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](n, hashInt, eqInt)
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
			benchMapHashSinkInt = acc
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
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
