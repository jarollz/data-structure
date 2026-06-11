package maptreeredblack

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

func collectSeq2[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	out := make(map[K]V)
	for k, v := range seq {
		out[k] = v
	}
	return out
}
