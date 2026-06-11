package listlinkedsingly

import "testing"

var benchListLinkedSinglySinkInt int
var benchListLinkedSinglySinkLarge benchLargePayload

func BenchmarkListLinkedSinglyPushFront(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int]()
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.PushFront(next)
				next++
				if list.Len() > n*2 {
					b.StopTimer()
					for list.Len() > n {
						list.PopFront()
					}
					b.StartTimer()
				}
			}
			benchListLinkedSinglySinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedSinglyPopFront(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int]()
			for i := 0; i < n; i++ {
				list.Append(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := list.PopFront()
				benchListLinkedSinglySinkInt = v
				if list.Len() == 0 {
					b.StopTimer()
					for j := 0; j < n; j++ {
						list.Append(j)
					}
					b.StartTimer()
				}
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedSinglyAppend(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int]()
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.Append(next)
				next++
				if list.Len() > n*2 {
					b.StopTimer()
					for list.Len() > n {
						list.PopFront()
					}
					b.StartTimer()
				}
			}
			benchListLinkedSinglySinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedSinglyDeleteFirst(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int]()
			for i := 0; i < n; i++ {
				list.Append(i)
			}
			match := func(v int) bool { return v == n/2 }
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.DeleteFirst(match)
				b.StopTimer()
				list.Append(n / 2)
				b.StartTimer()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedSinglyClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list := New[int]()
				for j := 0; j < n; j++ {
					list.Append(j)
				}
				list.Clear()
				benchListLinkedSinglySinkInt = list.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedSinglyClone(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload]()
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.Append(value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := list.Clone()
				for value := range cloned.Values() {
					benchListLinkedSinglySinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListLinkedSinglyCloneWith(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload]()
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
				for value := range cloned.Values() {
					benchListLinkedSinglySinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListLinkedSinglyValues(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload]()
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
			benchListLinkedSinglySinkLarge.Data[0] = sum
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
