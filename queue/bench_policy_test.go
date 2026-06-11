package queue

import (
	"math"
	"testing"
	"unsafe"
)

type benchComplexity int

const (
	benchO1 benchComplexity = iota
	benchOLogN
	benchON
)

type benchLargePayload struct {
	Data [32]uint64
}

func payloadBytes[T any]() int {
	var zero T
	return int(unsafe.Sizeof(zero))
}

func reportBenchmarkBudget(b *testing.B, complexity benchComplexity, payloadSize, n int) {
	b.Helper()
	actual := float64(b.Elapsed().Nanoseconds()) / float64(b.N)
	target := benchmarkBudgetNS(complexity, payloadSize, n)
	b.ReportMetric(target, "target-ns/op")
	b.ReportMetric(actual/target, "budget-ratio")
	switch complexity {
	case benchOLogN:
		b.ReportMetric(actual/math.Log2(float64(maxInt(2, n))), "ns/op/logn")
	case benchON:
		b.ReportMetric(actual/float64(maxInt(1, n)), "ns/op/n")
	}
}

func benchmarkBudgetNS(complexity benchComplexity, payloadSize, n int) float64 {
	copyTerm := float64(payloadSize) * 1.5
	switch complexity {
	case benchO1:
		return 120.0 + copyTerm
	case benchOLogN:
		return 160.0 + copyTerm + 18.0*math.Log2(float64(maxInt(2, n)))
	case benchON:
		return (6.0 + copyTerm/64.0) * float64(maxInt(1, n))
	default:
		return 0
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
