package maptrie

import (
	"iter"
	"strconv"
	"strings"
)

type pair[V any] struct {
	key   string
	value V
}

type cloneNode struct {
	value int
}

func collectPairs[V any](seq iter.Seq2[string, V]) []pair[V] {
	out := make([]pair[V], 0)
	for key, value := range seq {
		out = append(out, pair[V]{key: key, value: value})
	}
	return out
}

func collectKeys[V any](seq iter.Seq2[string, V]) []string {
	out := make([]string, 0)
	for key := range seq {
		out = append(out, key)
	}
	return out
}

func oracleHasPrefix[V any](model map[string]V, prefix string) bool {
	for key := range model {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

func trieKey(group, index int) string {
	return "g" + strconv.Itoa(group) + "/key/" + strconv.Itoa(index)
}

func randomKey(seed int) string {
	return string('a'+byte(seed%5)) + string('a'+byte((seed/5)%5)) + string('a'+byte((seed/25)%5))
}
