package stack

import "testing"

var benchStackSinkInt int
var benchStackSinkLarge benchLargePayload

func BenchmarkStackPush(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("tiny_size="+itoa(n), func(b *testing.B) {
			s := New[int](n)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Push(next)
				next++
				if s.Len() > n*2 {
					b.StopTimer()
					for s.Len() > n {
						s.Pop()
					}
					b.StartTimer()
				}
			}
			benchStackSinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			s := New[benchLargePayload](n)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(next)
				s.Push(value)
				next++
				if s.Len() > n*2 {
					b.StopTimer()
					for s.Len() > n {
						s.Pop()
					}
					b.StartTimer()
				}
			}
			benchStackSinkInt = next
			reportBenchmarkBudget(b, benchO1, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkStackPop(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](n)
			for i := 0; i < n; i++ {
				s.Push(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := s.Pop()
				benchStackSinkInt = v
				if s.Len() == 0 {
					b.StopTimer()
					for j := 0; j < n; j++ {
						s.Push(j)
					}
					b.StartTimer()
				}
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkStackPeekTop(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](n)
			for i := 0; i < n; i++ {
				s.Push(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := s.PeekTop()
				benchStackSinkInt = v
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkStackClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s := New[int](n)
				for j := 0; j < n; j++ {
					s.Push(j)
				}
				s.Clear()
				benchStackSinkInt = s.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkStackClone(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			s := New[benchLargePayload](n)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				s.Push(value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := s.Clone()
				v, _ := cloned.PeekTop()
				benchStackSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkStackCloneWith(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			s := New[benchLargePayload](n)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				s.Push(value)
			}
			cloneValue := func(v benchLargePayload) benchLargePayload {
				v.Data[1]++
				return v
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := s.CloneWith(cloneValue)
				v, _ := cloned.PeekTop()
				benchStackSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkStackValues(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			s := New[benchLargePayload](n)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				s.Push(value)
			}
			var sum uint64
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for value := range s.Values() {
					sum += value.Data[0]
				}
			}
			benchStackSinkLarge.Data[0] = sum
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkStackMixedPushPop(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](n)
			for i := 0; i < n; i++ {
				s.Push(i)
			}
			next := n
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Push(next)
				next++
				v, _ := s.Pop()
				benchStackSinkInt = v
			}
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
