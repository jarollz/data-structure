package treeredblack

import "iter"

func cmpInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func collectSeq[T any](seq iter.Seq[T]) []T {
	out := make([]T, 0)
	for v := range seq {
		out = append(out, v)
	}
	return out
}
