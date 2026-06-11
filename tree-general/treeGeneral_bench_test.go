package treegeneral

import "testing"

var benchTreeGeneralSinkInt int
var benchTreeGeneralSinkLarge benchLargePayload

func benchmarkGeneralTreeOfSize(n int) *TreeGeneral[benchLargePayload] {
	tree := New[benchLargePayload](benchLargePayload{})
	parents := make([]int, 0, n)
	parents = append(parents, 0)
	for i := 1; i < n; i++ {
		parentID := parents[(i-1)/2]
		childID, _ := tree.AddChild(parentID, benchLargePayload{Data: [32]uint64{uint64(i)}})
		parents = append(parents, childID)
	}
	return tree
}

func BenchmarkTreeGeneralAddChild(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := New[int](0)
			parents := make([]int, 1, n*2)
			parents[0] = 0
			nextParent := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				childID, _ := tree.AddChild(parents[nextParent], i)
				parents = append(parents, childID)
				nextParent++
				if nextParent == len(parents) {
					nextParent = 0
				}
			}
			benchTreeGeneralSinkInt = tree.Len()
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkTreeGeneralRemoveSubtree(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree := benchmarkGeneralTreeOfSize(n)
				tree.RemoveSubtree(1)
				benchTreeGeneralSinkInt = tree.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkTreeGeneralGet(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := benchmarkGeneralTreeOfSize(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v, _ := tree.Get(i % n)
				benchTreeGeneralSinkLarge = v
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkTreeGeneralParent(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := benchmarkGeneralTreeOfSize(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p, _ := tree.Parent((i % (n - 1)) + 1)
				benchTreeGeneralSinkInt = p
			}
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkTreeGeneralChildCount(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			tree := benchmarkGeneralTreeOfSize(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				count := tree.ChildCount(i % n)
				benchTreeGeneralSinkInt = count
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkTreeGeneralClone(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			tree := benchmarkGeneralTreeOfSize(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := tree.Clone()
				for value := range cloned.PreOrder() {
					benchTreeGeneralSinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkTreeGeneralCloneWith(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			tree := benchmarkGeneralTreeOfSize(n)
			cloneValue := func(v benchLargePayload) benchLargePayload {
				v.Data[1]++
				return v
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := tree.CloneWith(cloneValue)
				for value := range cloned.PreOrder() {
					benchTreeGeneralSinkLarge = value
					break
				}
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkTreeGeneralPreOrder(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			tree := benchmarkGeneralTreeOfSize(n)
			var sum uint64
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for value := range tree.PreOrder() {
					sum += value.Data[0]
				}
			}
			benchTreeGeneralSinkLarge.Data[0] = sum
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
