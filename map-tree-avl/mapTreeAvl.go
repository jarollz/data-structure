package maptreeavl

// New creates empty ordered map backed by AVL tree.
//
// cmp defines key ordering used by Put, Get, Delete, Min, and Max.
//
// Example: m := New[string, int](cmpString)
func New[K any, V any](cmp func(a, b K) int) *MapTreeAvl[K, V] {
	return nil
}
