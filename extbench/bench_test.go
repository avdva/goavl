package extbench

import (
	"math/rand"
	"testing"

	"github.com/avdva/goavl"
	gavl "github.com/karask/go-avltree"
	"github.com/tidwall/btree"
)

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
	b.Log(t.Len())
}

func BenchmarkTidwallBTreeFind(b *testing.B) {
	type item struct {
		k, v int
	}
	r := rand.New(rand.NewSource(0))
	t := btree.NewBTreeGOptions(func(a, b item) bool {
		return a.k < b.k
	}, btree.Options{
		NoLocks: true,
	})
	b.StopTimer()
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Int()
		keys[i] = k
		it := item{
			k: k,
			v: k,
		}
		t.Set(it)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		k := keys[i]
		it := item{
			k: k,
			v: k,
		}
		if _, found := t.Get(it); !found {
			panic("not found")
		}
	}
	b.Log(t.Len())
}

func BenchmarkTidwallBTreeDelete(b *testing.B) {
	type item struct {
		k, v int
	}
	r := rand.New(rand.NewSource(0))
	t := btree.NewBTreeGOptions(func(a, b item) bool {
		return a.k < b.k
	}, btree.Options{
		NoLocks: true,
	})
	b.StopTimer()
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Int()
		keys[i] = k
		it := item{
			k: k,
			v: k,
		}
		t.Set(it)
	}
	b.StartTimer()
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

func BenchmarkKaraskAVLInsert(b *testing.B) {
	tree := gavl.AVLTree{}
	r := rand.New(rand.NewSource(0))
	for i := 0; i < b.N; i++ {
		k := r.Int()
		tree.Add(k, k)
	}
}

func BenchmarkKaraskAVLFind(b *testing.B) {
	tree := gavl.AVLTree{}
	r := rand.New(rand.NewSource(0))
	b.StopTimer()
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Int()
		keys[i] = k
		tree.Add(k, k)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n := tree.Search(keys[i]); n == nil {
			panic("not found")
		}
	}
}

func BenchmarkKaraskAVLDelete(b *testing.B) {
	tree := gavl.AVLTree{}
	r := rand.New(rand.NewSource(0))
	b.StopTimer()
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Int()
		keys[i] = k
		tree.Add(k, k)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Remove(keys[i])
	}
}

func BenchmarkAVLInsert(b *testing.B) {
	tree := goavl.New[int, int](func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	r := rand.New(rand.NewSource(0))
	for i := 0; i < b.N; i++ {
		k := r.Int()
		tree.Insert(k, k)
	}
	b.Log(tree.Len())
}

func BenchmarkAVLFind(b *testing.B) {
	tree := goavl.New[int, int](func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Int()
		keys[i] = k
		tree.Insert(k, k)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if _, found := tree.Find(keys[i]); !found {
			panic("not found")
		}
	}
	b.Log(tree.Len())
}

func BenchmarkAVLDelete(b *testing.B) {
	tree := goavl.New[int, int](func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	keys := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		k := r.Int()
		keys[i] = k
		tree.Insert(k, k)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(keys[i])
	}
	if tree.Len() != 0 {
		panic("not empty")
	}
}
