package treeredblack

import "testing"

var benchTreeRedBlackSink int
var benchTreeRedBlackSinkLarge benchLargePayload

func BenchmarkTreeRedBlackInsert(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := New[int](cmpInt)
			next := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.Insert(next)
				next++
				if tree.Len() > n*2 {
					b.StopTimer()
					for j := next - n; j < next; j++ {
						tree.Delete(j)
					}
					next = n
					b.StartTimer()
				}
			}
			benchTreeRedBlackSink = next
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkTreeRedBlackDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := New[int](cmpInt)
			for i := 0; i < n; i++ {
				tree.Insert(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := i % n
				tree.Delete(k)
				b.StopTimer()
				tree.Insert(k)
				b.StartTimer()
			}
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkTreeRedBlackHas(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := New[int](cmpInt)
			for i := 0; i < n; i++ {
				tree.Insert(i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if tree.Has(i % n) {
					hits++
				}
			}
			benchTreeRedBlackSink = hits
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkTreeRedBlackMin(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := New[int](cmpInt)
			for i := 0; i < n; i++ {
				tree.Insert(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := tree.Min()
				benchTreeRedBlackSink = v
			}
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkTreeRedBlackMax(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := New[int](cmpInt)
			for i := 0; i < n; i++ {
				tree.Insert(i)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := tree.Max()
				benchTreeRedBlackSink = v
			}
			reportBenchmarkBudget(b, benchOLogN, payloadBytes[int](), n)
		})
	}
}

func BenchmarkTreeRedBlackClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := New[int](cmpInt)
				for j := 0; j < n; j++ {
					tree.Insert(j)
				}
				tree.Clear()
				benchTreeRedBlackSink = tree.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkTreeRedBlackClone(b *testing.B) {
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
			tree := New[benchLargePayload](cmpLarge)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				tree.Insert(value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := tree.Clone()
				for value := range cloned.InOrder() {
					benchTreeRedBlackSinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkTreeRedBlackCloneWith(b *testing.B) {
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
			tree := New[benchLargePayload](cmpLarge)
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				tree.Insert(value)
			}
			cloneValue := func(v benchLargePayload) benchLargePayload {
				v.Data[1]++
				return v
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := tree.CloneWith(cloneValue)
				for value := range cloned.InOrder() {
					benchTreeRedBlackSinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkTreeRedBlackInOrder(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := New[int](cmpInt)
			for i := 0; i < n; i++ {
				tree.Insert(i)
			}
			sum := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for value := range tree.InOrder() {
					sum += value
				}
			}
			benchTreeRedBlackSink = sum
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
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
