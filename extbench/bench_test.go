package extbench

import (
	"math/rand"
	"testing"

	"github.com/avdva/goavl"
	gavl "github.com/karask/go-avltree"
	"github.com/tidwall/btree"
)

func intCmp(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func BenchmarkTidwallBTreeInsert(b *testing.B) {
	type item struct {
		k, v int
	}
	r := rand.New(rand.NewSource(0))
	t := btree.NewBTreeGOptions(func(a, b item) bool {
		return a.k < b.k
	}, btree.Options{
		NoLocks: true,
	})
	for i := 0; i < b.N; i++ {
		k := r.Int()
		it := item{
			k: k,
			v: k,
		}
		t.Set(it)
	}
	if t.Len() == 0 {
		panic("empty")
	}
}

func BenchmarkKaraskAVLInsert(b *testing.B) {
	tree := gavl.AVLTree{}
	r := rand.New(rand.NewSource(0))
	for i := 0; i < b.N; i++ {
		k := r.Int()
		tree.Add(k, k)
	}
}

func BenchmarkAVLInsert(b *testing.B) {
	tree := goavl.New[int, int](intCmp)
	r := rand.New(rand.NewSource(0))
	for i := 0; i < b.N; i++ {
		k := r.Int()
		tree.Insert(k, k)
	}
	if tree.Len() == 0 {
		panic("empty")
	}
}

func BenchmarkGoMapInsert(b *testing.B) {
	m := make(map[int]int)
	r := rand.New(rand.NewSource(0))
	for i := 0; i < b.N; i++ {
		k := r.Int()
		m[k] = k
	}
	if len(m) == 0 {
		panic("empty")
	}
}

func BenchmarkTidwallBTreeFind(b *testing.B) {
	type item struct {
		k, v int
	}
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	t := btree.NewBTreeGOptions(func(a, b item) bool {
		return a.k < b.k
	}, btree.Options{
		NoLocks: true,
	})
	count := 250000
	keys := make([]int, count)
	for i := 0; i < count; i++ {
		k := r.Int()
		keys[i] = k
		it := item{
			k: k,
			v: k,
		}
		t.Set(it)
	}
	b.StartTimer()
	var sum int
	for i := 0; i < b.N; i++ {
		toFind := keys[i%len(keys)]
		it, _ := t.Get(item{
			k: toFind,
		})
		sum += it.v
	}
	b.Logf("sum = %d", sum)
}

func BenchmarkAVLFind(b *testing.B) {
	b.StopTimer()
	tree := goavl.New[int, int](intCmp)
	r := rand.New(rand.NewSource(0))
	count := 250000
	keys := make([]int, count)
	for i := 0; i < count; i++ {
		k := r.Int()
		keys[i] = k
		tree.Insert(k, k)
	}
	b.StartTimer()
	var sum int
	for i := 0; i < b.N; i++ {
		toFind := keys[i%len(keys)]
		v, _ := tree.Find(toFind)
		sum += v
	}
	b.Logf("sum = %d", sum)
}

func BenchmarkTidwallBTreeAt(b *testing.B) {
	type item struct {
		k, v int
	}
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	t := btree.NewBTreeGOptions(func(a, b item) bool {
		return a.k < b.k
	}, btree.Options{
		NoLocks: true,
	})
	for i := 0; i < b.N; i++ {
		k := r.Intn(50000)
		it := item{
			k: k,
			v: k,
		}
		t.Set(it)
	}
	b.StartTimer()
	var sum int
	for i := 0; i < b.N; i++ {
		it, _ := t.GetAt(r.Int() % t.Len())
		sum += it.k
	}
	b.Logf("sum = %d", sum)
}

func BenchmarkAVLAt(b *testing.B) {
	tree := goavl.New[int, int](intCmp, goavl.WithCountChildren(true))
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	for i := 0; i < b.N; i++ {
		k := r.Intn(50000)
		tree.Insert(k, k)
	}
	var sum int
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		k, _ := tree.At(r.Int() % tree.Len())
		sum += k
	}
	b.Logf("sum = %d", sum)
}
