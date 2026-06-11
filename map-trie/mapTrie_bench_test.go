package maptrie

import (
	"strconv"
	"testing"
)

var benchMapTrieSinkInt int
var benchMapTrieSinkString string
var benchMapTrieSinkLarge benchLargePayload

func BenchmarkMapTriePut(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(maxInt(n, b.N))
			m := New[int]()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Put(keys[i], i)
			}
			benchMapTrieSinkInt = m.Len()
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTrieGetHit(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(n)
			m := New[int]()
			for i := 0; i < n; i++ {
				m.Put(keys[i], i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, ok := m.Get(keys[i%n]); ok {
					hits++
				}
			}
			benchMapTrieSinkInt = hits
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTrieGetMiss(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(n)
			misses := benchKeysFrom(n, n)
			m := New[int]()
			for i := 0; i < n; i++ {
				m.Put(keys[i], i)
			}
			count := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, ok := m.Get(misses[i%n]); !ok {
					count++
				}
			}
			benchMapTrieSinkInt = count
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTrieDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(n)
			m := New[int]()
			for i := 0; i < n; i++ {
				m.Put(keys[i], i)
			}
			deleted := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := keys[i%n]
				if m.Delete(key) {
					deleted++
				}
				b.StopTimer()
				m.Put(key, i)
				b.StartTimer()
			}
			benchMapTrieSinkInt = deleted
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTrieHas(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(n)
			m := New[int]()
			for i := 0; i < n; i++ {
				m.Put(keys[i], i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if m.Has(keys[i%n]) {
					hits++
				}
			}
			benchMapTrieSinkInt = hits
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTrieHasPrefix(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchGroupedKeys(n)
			prefixes := benchPrefixes()
			m := New[int]()
			for i := 0; i < n; i++ {
				m.Put(keys[i], i)
			}
			hits := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if m.HasPrefix(prefixes[i%len(prefixes)]) {
					hits++
				}
			}
			benchMapTrieSinkInt = hits
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTrieClear(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m := New[int]()
				for j := 0; j < n; j++ {
					m.Put(keys[j], j)
				}
				m.Clear()
				benchMapTrieSinkInt = m.Len()
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTrieClone(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(n)
			m := New[benchLargePayload]()
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				m.Put(keys[i], value)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := m.Clone()
				v, _ := cloned.Get(keys[0])
				benchMapTrieSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkMapTrieCloneWith(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("large_size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(n)
			m := New[benchLargePayload]()
			for i := 0; i < n; i++ {
				var value benchLargePayload
				value.Data[0] = uint64(i)
				m.Put(keys[i], value)
			}
			cloneValue := func(v benchLargePayload) benchLargePayload {
				v.Data[1]++
				return v
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cloned := m.CloneWith(cloneValue)
				v, _ := cloned.Get(keys[0])
				benchMapTrieSinkLarge = v
			}
			reportBenchmarkBudget(b, benchON, payloadBytes[benchLargePayload](), n)
		})
	}
}

func BenchmarkMapTrieAll(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(n)
			m := New[int]()
			for i := 0; i < n; i++ {
				m.Put(keys[i], i)
			}
			sum := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for key, value := range m.All() {
					sum += len(key) + value
				}
			}
			benchMapTrieSinkInt = sum
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTrieWithPrefix(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchGroupedKeys(n)
			prefixes := benchPrefixes()
			m := New[int]()
			for i := 0; i < n; i++ {
				m.Put(keys[i], i)
			}
			sum := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for key, value := range m.WithPrefix(prefixes[i%len(prefixes)]) {
					sum += len(key) + value
				}
			}
			benchMapTrieSinkInt = sum
			reportBenchmarkBudget(b, benchON, payloadBytes[int](), n)
		})
	}
}

func BenchmarkMapTrieMixedPutGetDelete(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run("size="+itoa(n), func(b *testing.B) {
			keys := benchKeys(n)
			extra := benchKeysFrom(n, maxInt(n, b.N))
			m := New[int]()
			for i := 0; i < n; i++ {
				m.Put(keys[i], i)
			}
			acc := 0
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Put(extra[i], i)
				if v, ok := m.Get(keys[i%n]); ok {
					acc += v
				}
				m.Delete(keys[i%n])
				b.StopTimer()
				m.Put(keys[i%n], i)
				b.StartTimer()
			}
			benchMapTrieSinkInt = acc
			reportBenchmarkBudget(b, benchO1, payloadBytes[int](), n)
		})
	}
}

func benchKeys(n int) []string {
	return benchKeysFrom(0, n)
}

func benchKeysFrom(start, n int) []string {
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = "key/" + strconv.Itoa(start+i)
	}
	return keys
}

func benchGroupedKeys(n int) []string {
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = "group/" + strconv.Itoa(i%10) + "/key/" + strconv.Itoa(i)
	}
	return keys
}

func benchPrefixes() []string {
	return []string{
		"group/0",
		"group/1",
		"group/2",
		"group/3",
		"group/4",
	}
}

func itoa(v int) string {
	if v == 1_000 {
		return "1e3"
	}
	if v == 10_000 {
		return "1e4"
	}
	if v == 100_000 {
		return "1e5"
	}
	return strconv.Itoa(v)
}
