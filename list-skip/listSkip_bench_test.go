package listskip

import "testing"

var benchListSkipSinkInt int
var benchListSkipSinkLarge benchLargePayload

func BenchmarkListSkipInsert(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int](8, cmpInt)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.Insert(next)
				next++
				if list.Len() > n*2 {
					b.StopTimer()
					for j := next - n; j < next; j++ {
						list.Delete(j)
					}
					next = n
					b.StartTimer()
				}
			}
			benchListSkipSinkInt = next
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListSkipDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int](8, cmpInt)
			for i := 0; i < n; i++ {
				list.Insert(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := i % n
				list.Delete(k)
				b.StopTimer()
				list.Insert(k)
				b.StartTimer()
			}
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListSkipHas(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int](8, cmpInt)
			for i := 0; i < n; i++ {
				list.Insert(i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if list.Has(i % n) {
					hits++
				}
			}
			benchListSkipSinkInt = hits
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListSkipClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list := New[int](8, cmpInt)
				for j := 0; j < n; j++ {
					list.Insert(j)
				}
				list.Clear()
				benchListSkipSinkInt = list.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListSkipClone(b *testing.B) {
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
			list := New[benchLargePayload](8, cmpLarge)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.Insert(value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := list.Clone()
				for value := range cloned.Values() {
					benchListSkipSinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListSkipCloneWith(b *testing.B) {
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
			list := New[benchLargePayload](8, cmpLarge)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.Insert(value)
			}
			cloneValue := func(v benchLargePayload) benchLargePayload {
				v.Data[1]++
				return v
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := list.CloneWith(cloneValue)
				for value := range cloned.Values() {
					benchListSkipSinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListSkipValues(b *testing.B) {
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
			list := New[benchLargePayload](8, cmpLarge)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.Insert(value)
			}
			var sum uint64
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for value := range list.Values() {
					sum += value.Data[0]
				}
			}
			benchListSkipSinkLarge.Data[0] = sum
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListSkipMixedOrdered(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int](8, cmpInt)
			for i := 0; i < n; i++ {
				list.Insert(i)
			}
			next := n
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.Insert(next)
				next++
				list.Has(i % n)
				list.Delete(i % n)
				b.StopTimer()
				list.Insert(i % n)
				b.StartTimer()
			}
			benchListSkipSinkInt = next
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
