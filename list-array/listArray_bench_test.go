package listarray

import "testing"

var benchListArraySinkInt int
var benchListArraySinkLarge benchLargePayload

func BenchmarkListArrayAppend(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("tiny_size="+itoa(n), func(b *testing.B) {
			list := New[int](n)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.Append(next)
				next++
				if list.Len() > n*2 {
					b.StopTimer()
					for list.Len() > n {
						list.Delete(list.Len() - 1)
					}
					b.StartTimer()
				}
			}
			benchListArraySinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload](n)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(next)
				list.Append(value)
				next++
				if list.Len() > n*2 {
					b.StopTimer()
					for list.Len() > n {
						list.Delete(list.Len() - 1)
					}
					b.StartTimer()
				}
			}
			benchListArraySinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListArrayGet(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int](n)
			for i := 0; i < n; i++ {
				list.Append(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := list.Get(i % n)
				benchListArraySinkInt = v
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListArraySet(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int](n)
			for i := 0; i < n; i++ {
				list.Append(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.Set(i%n, i)
			}
			v, _ := list.Get((b.N - 1 + n) % n)
			benchListArraySinkInt = v
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListArrayInsert(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int](n)
			for i := 0; i < n; i++ {
				list.Append(i)
			}
			mid := n / 2
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.Insert(mid, i)
				b.StopTimer()
				list.Delete(mid)
				b.StartTimer()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListArrayDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int](n)
			for i := 0; i < n; i++ {
				list.Append(i)
			}
			mid := n / 2
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := list.Delete(mid)
				benchListArraySinkInt = v
				b.StopTimer()
				list.Insert(mid, v)
				b.StartTimer()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListArrayClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list := New[int](n)
				for j := 0; j < n; j++ {
					list.Append(j)
				}
				list.Clear()
				benchListArraySinkInt = list.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListArrayClone(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload](n)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.Append(value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := list.Clone()
				v, _ := cloned.Get(0)
				benchListArraySinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListArrayCloneWith(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload](n)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.Append(value)
			}
			cloneValue := func(v benchLargePayload) benchLargePayload {
				v.Data[1]++
				return v
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := list.CloneWith(cloneValue)
				v, _ := cloned.Get(0)
				benchListArraySinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListArrayValues(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload](n)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.Append(value)
			}
			var sum uint64
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for value := range list.Values() {
					sum += value.Data[0]
				}
			}
			benchListArraySinkLarge.Data[0] = sum
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
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
