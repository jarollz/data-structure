package treeredblack

import "testing"

var benchTreeRedBlackSink int

func BenchmarkTreeRedBlackInsert(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](cmpInt)
			for i := 0; i < n; i++ {
				t.Insert(i)
			}
			next := n
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				t.Insert(next)
				next++
			}
			benchTreeRedBlackSink = next
		})
	}
}

func BenchmarkTreeRedBlackHas(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](cmpInt)
			for i := 0; i < n; i++ {
				t.Insert(i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if t.Has(i % n) {
					hits++
				}
			}
			benchTreeRedBlackSink = hits
		})
	}
}

func BenchmarkTreeRedBlackDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](cmpInt)
			for i := 0; i < n; i++ {
				t.Insert(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := i % n
				t.Delete(k)
				b.StopTimer()
				t.Insert(k)
				b.StartTimer()
			}
		})
	}
}

func BenchmarkTreeRedBlackInOrder(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](cmpInt)
			for i := 0; i < n; i++ {
				t.Insert(i)
			}
			sum := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for v := range t.InOrder() {
					sum += v
				}
			}
			benchTreeRedBlackSink = sum
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
