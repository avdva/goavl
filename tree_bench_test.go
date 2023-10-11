package goavl

import (
	"testing"
)

func BenchmarkTree_At_WithCountChildren(b *testing.B) {
	benchmarkTreeAtFirstN(b, 16392, 16392, WithCountChildren(true))
}

func BenchmarkTree_At_WithoutCountChildren(b *testing.B) {
	benchmarkTreeAtFirstN(b, 16392, 16392, WithCountChildren(false))
}

func benchmarkTreeAtFirstN(b *testing.B, total, n int, opts ...Option) {
	tree := NewComparable[int, int](opts...)
	b.StopTimer()
	for i := 0; i <= total; i++ {
		tree.Insert(i, i)
	}
	b.StartTimer()
	var sum int
	for outer := 0; outer < b.N; outer++ {
		for i := 0; i < n; i++ {
			e := tree.At(i)
			sum += e.Key
			e = tree.At(tree.Len() - i - 1)
			sum += e.Key
		}
	}
	b.Logf("the sum is: %d", sum)
}
