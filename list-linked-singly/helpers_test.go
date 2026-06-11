package listlinkedsingly

import "iter"

func collectSeq[T any](seq iter.Seq[T]) []T {
	out := make([]T, 0)
	for v := range seq {
		out = append(out, v)
	}
	return out
}
