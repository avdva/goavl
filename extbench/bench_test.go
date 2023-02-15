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

func BenchmarkTidwallBTree(b *testing.B) {
	type item struct {
		k, v int
	}
	r := rand.New(rand.NewSource(0))
	t := btree.NewBTreeGOptions(func(a, b item) bool {
		return a.k < b.k
	}, btree.Options{
		NoLocks: true,
	})
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Intn(50000)
		keys[i] = k
		it := item{
			k: k,
			v: k,
		}
		t.Set(it)
		if _, found := t.Get(it); !found {
			panic("not found")
		}
	}
	for i := 0; i < b.N; i++ {
		it := item{
			k: keys[i],
		}
		t.Delete(it)
	}
	if t.Len() != 0 {
		panic("not empty")
	}
}

func BenchmarkKaraskAVL(b *testing.B) {
	tree := gavl.AVLTree{}
	r := rand.New(rand.NewSource(0))
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Intn(50000)
		keys[i] = k
		tree.Add(k, k)
		if n := tree.Search(k); n == nil {
			panic("not found")
		}
	}
	for i := 0; i < b.N; i++ {
		tree.Remove(keys[i])
	}
}

func BenchmarkAVL(b *testing.B) {
	tree := goavl.New[int, int](intCmp)
	r := rand.New(rand.NewSource(0))
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Intn(50000)
		keys[i] = k
		tree.Insert(k, k)
		if _, found := tree.Find(k); !found {
			panic("not found")
		}
	}
	for i := 0; i < b.N; i++ {
		tree.Delete(keys[i])
	}
	if tree.Len() != 0 {
		panic("not empty")
	}
}

func BenchmarkGoMap(b *testing.B) {
	m := make(map[int]int)
	r := rand.New(rand.NewSource(0))
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Intn(50000)
		keys[i] = k
		m[k] = k
		if _, found := m[k]; !found {
			panic("not found")
		}
	}
	for i := 0; i < b.N; i++ {
		delete(m, keys[i])
	}
	if len(m) != 0 {
		panic("not empty")
	}
}
