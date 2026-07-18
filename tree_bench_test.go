package goavl

import (
	"math/rand"
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

func BenchmarkTreeAllocsSimpleCache(b *testing.B) {
	benchmarkTreeAllocs(b, 10000)
}

func BenchmarkTreeAllocsSyncPool(b *testing.B) {
	benchmarkTreeAllocs(b, 10000, WithSyncPoolAllocator(true))
}

func benchmarkTreeAllocs(b *testing.B, n int, opts ...Option) {
	tree := NewComparable[int, int](opts...)
	var sum int
	b.StopTimer()
	for i := 0; i < n; i++ {
		tree.Insert(i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < n; i++ {
			tree.Delete(i)
		}
		for i := 0; i < n; i++ {
			ptr, _ := tree.Insert(i, i)
			sum += *ptr
		}
	}
	b.Logf("the sum is: %d", sum)
}

func BenchmarkTree_UpdateKeyRandom(b *testing.B) {
	benchmarkTreeUpdateKeyRandom(b, 10000)
}

func BenchmarkTree_DeleteInsertRandom(b *testing.B) {
	benchmarkTreeDeleteInsertRandom(b, 10000)
}

func BenchmarkTree_UpdateKeyMixed(b *testing.B) {
	benchmarkTreeUpdateKeyMixed(b, 10000)
}

func BenchmarkTree_DeleteInsertMixed(b *testing.B) {
	benchmarkTreeDeleteInsertMixed(b, 10000)
}

func benchmarkTreeUpdateKeyRandom(b *testing.B, n int) {
	keys := shuffledInts(n)
	var sum int
	b.ReportAllocs()
	for outer := 0; outer < b.N; outer++ {
		b.StopTimer()
		tree := benchmarkTreeWithKeys(n)
		b.StartTimer()
		for _, oldKey := range keys {
			ptr, updated := tree.UpdateKey(oldKey, oldKey+n)
			if !updated {
				b.Fatalf("key was not updated: %d", oldKey)
			}
			sum += *ptr
		}
	}
	b.Logf("the sum is: %d", sum)
}

func benchmarkTreeDeleteInsertRandom(b *testing.B, n int) {
	keys := shuffledInts(n)
	var sum int
	b.ReportAllocs()
	for outer := 0; outer < b.N; outer++ {
		b.StopTimer()
		tree := benchmarkTreeWithKeys(n)
		b.StartTimer()
		for _, oldKey := range keys {
			ptr, found := tree.Find(oldKey)
			if !found {
				b.Fatalf("key was not found: %d", oldKey)
			}
			value := *ptr
			tree.Delete(oldKey)
			ptr, _ = tree.Insert(oldKey+n, value)
			sum += *ptr
		}
	}
	b.Logf("the sum is: %d", sum)
}

func benchmarkTreeUpdateKeyMixed(b *testing.B, n int) {
	groups := n / 4
	ops := shuffledInts(groups * 3)
	var sum int
	b.ReportAllocs()
	for outer := 0; outer < b.N; outer++ {
		b.StopTimer()
		tree := benchmarkMixedTree(groups)
		b.StartTimer()
		for _, op := range ops {
			oldKey, newKey := benchmarkMixedUpdateKeys(groups, op)
			ptr, updated := tree.UpdateKey(oldKey, newKey)
			if !updated {
				b.Fatalf("key was not updated: %d", oldKey)
			}
			sum += *ptr
		}
	}
	b.Logf("the sum is: %d", sum)
}

func benchmarkTreeDeleteInsertMixed(b *testing.B, n int) {
	groups := n / 4
	ops := shuffledInts(groups * 3)
	var sum int
	b.ReportAllocs()
	for outer := 0; outer < b.N; outer++ {
		b.StopTimer()
		tree := benchmarkMixedTree(groups)
		b.StartTimer()
		for _, op := range ops {
			oldKey, newKey := benchmarkMixedUpdateKeys(groups, op)
			ptr, found := tree.Find(oldKey)
			if !found {
				b.Fatalf("key was not found: %d", oldKey)
			}
			value := *ptr
			tree.Delete(oldKey)
			ptr, _ = tree.Insert(newKey, value)
			sum += *ptr
		}
	}
	b.Logf("the sum is: %d", sum)
}

func benchmarkTreeWithKeys(n int) *Tree[int, int, func(a int, b int) int] {
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < n; i++ {
		tree.Insert(i, i)
	}
	return tree
}

func benchmarkMixedTree(groups int) *Tree[int, int, func(a int, b int) int] {
	tree := NewComparable[int, int](WithCountChildren(true))
	for group := 0; group < groups; group++ {
		base := group * 8
		tree.Insert(base, base)
		tree.Insert(base+2, base+2)
		tree.Insert(base+4, base+4)
		tree.Insert(base+6, base+6)
	}
	return tree
}

func benchmarkMixedUpdateKeys(groups int, op int) (oldKey int, newKey int) {
	group := op / 3
	base := group * 8
	switch op % 3 {
	case 0:
		return base + 2, base + 3
	case 1:
		return base, base + 4
	default:
		return base + 6, groups*8 + 1024 + group
	}
}

func shuffledInts(n int) []int {
	result := make([]int, n)
	for i := range result {
		result[i] = i
	}
	rand.New(rand.NewSource(1)).Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}
