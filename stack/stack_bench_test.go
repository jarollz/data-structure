package stack

import "testing"

var benchStackSink int

func BenchmarkStackPush(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int](n)
			for i := 0; i < n; i++ {
				s.Push(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Push(i)
			}
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
				benchStackSink = v
				b.StopTimer()
				s.Push(i)
				b.StartTimer()
			}
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
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Push(i)
				v, _ := s.Pop()
				benchStackSink = v
			}
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
