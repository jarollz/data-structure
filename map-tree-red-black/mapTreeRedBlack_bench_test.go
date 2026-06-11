package maptreeredblack

import "testing"

var benchMapTreeRedBlackSinkInt int
var benchMapTreeRedBlackSinkLarge benchLargePayload

func BenchmarkMapTreeRedBlackPut(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Put(next, next)
				next++
				if m.Len() > n*2 {
					b.StopTimer()
					for j := next - n; j < next; j++ {
						m.Delete(j)
					}
					next = n
					b.StartTimer()
				}
			}
			benchMapTreeRedBlackSinkInt = next
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackGet(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if v, ok := m.Get(i % n); ok {
					hits += v
				}
			}
			benchMapTreeRedBlackSinkInt = hits
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackDelete(b *testing.B) {
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
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackHas(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
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
			benchMapTreeRedBlackSinkInt = hits
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackMin(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k, _, _ := m.Min()
				benchMapTreeRedBlackSinkInt = k
			}
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackMax(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			m := New[int, int](cmpInt)
			for i := 0; i < n; i++ {
				m.Put(i, i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k, _, _ := m.Max()
				benchMapTreeRedBlackSinkInt = k
			}
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m := New[int, int](cmpInt)
				for j := 0; j < n; j++ {
					m.Put(j, j)
				}
				m.Clear()
				benchMapTreeRedBlackSinkInt = m.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackClone(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			m := New[int, benchLargePayload](cmpInt)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				m.Put(i, value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := m.Clone()
				_, v, _ := cloned.Min()
				benchMapTreeRedBlackSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackCloneWith(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			m := New[int, benchLargePayload](cmpInt)
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
				_, v, _ := cloned.Min()
				benchMapTreeRedBlackSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackAll(b *testing.B) {
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
			benchMapTreeRedBlackSinkInt = sum
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTreeRedBlackMixedReadWrite(b *testing.B) {
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
				b.StopTimer()
				m.Put(i%n, i)
				b.StartTimer()
			}
			benchMapTreeRedBlackSinkInt = acc
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
