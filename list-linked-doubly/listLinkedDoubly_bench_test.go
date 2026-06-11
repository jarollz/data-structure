package listlinkeddoubly

import "testing"

var benchListLinkedDoublySinkInt int
var benchListLinkedDoublySinkLarge benchLargePayload

func BenchmarkListLinkedDoublyPushFront(b *testing.B) {
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
						list.PopBack()
					}
					b.StartTimer()
				}
			}
			benchListLinkedDoublySinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedDoublyPushBack(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int]()
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list.PushBack(next)
				next++
				if list.Len() > n*2 {
					b.StopTimer()
					for list.Len() > n {
						list.PopFront()
					}
					b.StartTimer()
				}
			}
			benchListLinkedDoublySinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedDoublyPopFront(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int]()
			for i := 0; i < n; i++ {
				list.PushBack(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := list.PopFront()
				benchListLinkedDoublySinkInt = v
				if list.Len() == 0 {
					b.StopTimer()
					for j := 0; j < n; j++ {
						list.PushBack(j)
					}
					b.StartTimer()
				}
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedDoublyPopBack(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			list := New[int]()
			for i := 0; i < n; i++ {
				list.PushBack(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := list.PopBack()
				benchListLinkedDoublySinkInt = v
				if list.Len() == 0 {
					b.StopTimer()
					for j := 0; j < n; j++ {
						list.PushBack(j)
					}
					b.StartTimer()
				}
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedDoublyClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				list := New[int]()
				for j := 0; j < n; j++ {
					list.PushBack(j)
				}
				list.Clear()
				benchListLinkedDoublySinkInt = list.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkListLinkedDoublyClone(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload]()
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.PushBack(value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := list.Clone()
				for value := range cloned.Values() {
					benchListLinkedDoublySinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListLinkedDoublyCloneWith(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload]()
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.PushBack(value)
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
					benchListLinkedDoublySinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkListLinkedDoublyValues(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			list := New[benchLargePayload]()
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				list.PushBack(value)
			}
			var sum uint64
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for value := range list.Values() {
					sum += value.Data[0]
				}
			}
			benchListLinkedDoublySinkLarge.Data[0] = sum
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
