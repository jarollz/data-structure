package maphash

import "iter"

func hashInt(v int) uint64 {
	return uint64(uint32(v) * 2654435761)
}

func eqInt(a, b int) bool {
	return a == b
}

func collectSeq2[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	out := make(map[K]V)
	for k, v := range seq {
		out[k] = v
	}
	return out
}
