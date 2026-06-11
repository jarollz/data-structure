package listlinkedsingly

import "testing"

var benchListLinkedSinglySink int

func BenchmarkListLinkedSinglyPushFrontPopFront(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int]()
			for i := 0; i < n; i++ {
				s.PushFront(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.PushFront(i)
				v, _ := s.PopFront()
				benchListLinkedSinglySink = v
			}
		})
	}
}

func BenchmarkListLinkedSinglyAppend(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int]()
			for i := 0; i < n; i++ {
				s.Append(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Append(i)
			}
		})
	}
}

func BenchmarkListLinkedSinglyValues(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			s := New[int]()
			for i := 0; i < n; i++ {
				s.Append(i)
			}
			sum := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for v := range s.Values() {
					sum += v
				}
			}
			benchListLinkedSinglySink = sum
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
