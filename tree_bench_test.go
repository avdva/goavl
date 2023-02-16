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
			k, _ := tree.At(i)
			sum += k
			k, _ = tree.At(tree.Len() - i - 1)
			sum += k
		}
	}
	b.Logf("the sum is: %d", sum)
}
