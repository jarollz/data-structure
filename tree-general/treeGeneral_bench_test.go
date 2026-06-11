package treegeneral

import "testing"

var benchTreeGeneralSink int

func BenchmarkTreeGeneralAddChild(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](0)
			for i := 0; i < n; i++ {
				t.AddChild(0, i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				id, _ := t.AddChild(0, n+i)
				benchTreeGeneralSink = id
			}
		})
	}
}

func BenchmarkTreeGeneralRemoveSubtree(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](0)
			ids := make([]int, n)
			for i := 0; i < n; i++ {
				id, _ := t.AddChild(0, i)
				ids[i] = id
				_, _ = t.AddChild(id, i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := i % n
				t.RemoveSubtree(ids[idx])
				b.StopTimer()
				id, _ := t.AddChild(0, n+i)
				ids[idx] = id
				_, _ = t.AddChild(id, n+i)
				b.StartTimer()
			}
		})
	}
}

func BenchmarkTreeGeneralPreOrder(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			t := New[int](0)
			for i := 0; i < n; i++ {
				t.AddChild(0, i)
			}
			sum := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for v := range t.PreOrder() {
					sum += v
				}
			}
			benchTreeGeneralSink = sum
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
