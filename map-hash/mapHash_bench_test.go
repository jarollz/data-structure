package maphash

import "testing"

var benchMapHashSinkInt int
var benchMapHashSinkFloat float64
var benchMapHashSinkLarge benchLargePayload

func BenchmarkMapHashPut(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](n, hashInt, eqInt)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Put(next, next)
				next++
			}
			benchMapHashSinkInt = next
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapHashGetHit(b *testing.B) {
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
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapHashGetMiss(b *testing.B) {
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
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapHashDelete(b *testing.B) {
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
				b.StopTimer()
				m.Put(k, i)
				b.StartTimer()
			}
			benchMapHashSinkInt = deleted
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapHashHas(b *testing.B) {
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
				if m.Has(i % n) {
					hits++
				}
			}
			benchMapHashSinkInt = hits
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapHashClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m := New[int, int](n, hashInt, eqInt)
				for j := 0; j < n; j++ {
					m.Put(j, j)
				}
				m.Clear()
				benchMapHashSinkInt = m.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapHashClone(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			m := New[int, benchLargePayload](n, hashInt, eqInt)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				m.Put(i, value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := m.Clone()
				v, _ := cloned.Get(0)
				benchMapHashSinkLarge = v
			}
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkMapHashCloneWith(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			m := New[int, benchLargePayload](n, hashInt, eqInt)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				m.Put(i, value)
			}
			cloneValue := func(v benchLargePayload) benchLargePayload {
				v.Data[1]++
				return v
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := m.CloneWith(func(k int) int { return k }, cloneValue)
				v, _ := cloned.Get(0)
				benchMapHashSinkLarge = v
			}
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkMapHashAll(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](n, hashInt, eqInt)
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
			benchMapHashSinkInt = sum
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
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
				b.StopTimer()
				m.Put(i%n, i)
				b.StartTimer()
			}
			benchMapHashSinkInt = acc
			lf := m.LoadFactor()
			benchMapHashSinkFloat = lf
			b.ReportMetric(lf, "loadfactor")
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
