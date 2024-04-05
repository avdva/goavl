package goavl

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyTree(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	it := tree.AscendFromStart()
	e, ok := it.Next()
	a.Equal(0, e.Key)
	a.Equal((*int)(nil), e.Value)
	a.Equal(false, ok)
	v, ok := tree.Delete(0)
	a.Equal(0, v)
	a.Equal(false, ok)
	it = tree.DescendFromEnd()
	e, ok = it.Prev()
	a.Equal(0, e.Key)
	a.Equal((*int)(nil), e.Value)
	a.Equal(false, ok)
	val, ok := tree.Find(0)
	a.Equal((*int)(nil), val)
	a.Equal(false, ok)
	e, ok = tree.Max()
	a.Equal(0, e.Key)
	a.Equal((*int)(nil), e.Value)
	a.Equal(false, ok)
	e, ok = tree.Min()
	a.Equal(0, e.Key)
	a.Equal((*int)(nil), e.Value)
	a.Equal(false, ok)
	a.Zero(tree.Len())
	tree.Clear()
	a.Zero(tree.Len())
}

func TestTreeInsert(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.Truef(inserted, "k: %v", i)
		e, found := tree.Min()
		a.Equal(0, e.Key)
		a.Equal(0, *e.Value)
		a.True(found)

		e, found = tree.Max()
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
		a.True(found)
		a.NoErrorf(checkHeightAndBalance(tree.root), "iter = %d", i)
	}
	for i := 0; i < 128; i++ {
		val, found := tree.Find(i)
		a.True(found)
		a.Equal(i, *val)
	}

	for i := 127; i >= 0; i-- {
		ptr, inserted := tree.Insert(i, i*2)
		a.Equal(i*2, *ptr)
		a.Falsef(inserted, "k: %v", i)
		a.NoError(checkHeightAndBalance(tree.root))
	}
	for i := 127; i >= 0; i-- {
		val, found := tree.Find(i)
		a.True(found)
		a.Equal(i*2, *val)
	}
}

func TestTreeDelete(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	a.Equal(0, tree.Len())

	ptr, inserted := tree.Insert(0, 0)
	a.Equal(0, *ptr)
	a.True(inserted)
	v, deleted := tree.Delete(0)
	a.True(deleted)
	a.Equal(0, v)
	a.Equal(0, tree.Len())
	a.True(tree.root.isNil())
	a.NoError(checkHeightAndBalance(tree.root))

	ptr, inserted = tree.Insert(0, 0)
	a.Equal(0, *ptr)
	a.True(inserted)
	ptr, inserted = tree.Insert(-1, -1)
	a.Equal(-1, *ptr)
	a.True(inserted)
	a.Equal(2, tree.Len())
	a.NoError(checkHeightAndBalance(tree.root))
	v, deleted = tree.Delete(0)
	a.True(deleted)
	a.Equal(0, v)
	v, deleted = tree.Delete(-1)
	a.Equal(-1, v)
	a.True(deleted)

	ptr, inserted = tree.Insert(0, 0)
	a.Equal(0, *ptr)
	a.True(inserted)
	ptr, inserted = tree.Insert(1, 1)
	a.Equal(1, *ptr)
	a.True(inserted)
	a.Equal(2, tree.Len())
	a.NoError(checkHeightAndBalance(tree.root))
	v, deleted = tree.Delete(1)
	a.True(deleted)
	a.Equal(1, v)
	a.Equal(1, tree.Len())
	_, deleted = tree.Delete(-1)
	a.False(deleted)
	a.NoError(checkHeightAndBalance(tree.root))
	v, deleted = tree.Delete(0)
	a.True(deleted)
	a.Equal(0, v)
	a.Equal(0, tree.Len())
	a.True(tree.root.isNil())
	a.NoError(checkHeightAndBalance(tree.root))

	ptr, inserted = tree.Insert(0, 0)
	a.Equal(0, *ptr)
	a.True(inserted)
	ptr, inserted = tree.Insert(1, 1)
	a.Equal(1, *ptr)
	a.True(inserted)
	a.NoError(checkHeightAndBalance(tree.root))
	v, deleted = tree.Delete(0)
	a.True(deleted)
	a.Equal(0, v)
	a.NoError(checkHeightAndBalance(tree.root))
	a.Equal(1, tree.Len())
	v, deleted = tree.Delete(1)
	a.True(deleted)
	a.Equal(1, v)
	a.NoError(checkHeightAndBalance(tree.root))
	a.True(tree.root.isNil())
	a.NoError(checkHeightAndBalance(tree.root))

	for i := 128; i <= 0; i-- {
		ptr, inserted = tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
		a.NoError(checkHeightAndBalance(tree.root))
	}
	for i := 128; i <= 0; i-- {
		v, deleted = tree.Delete(i)
		a.True(deleted)
		a.Equal(i, v)
		a.NoError(checkHeightAndBalance(tree.root))
	}
	a.Equal(0, tree.Len())
}

func TestTreeDeleteMin(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	for i := 0; i < 128; i++ {
		e, found := tree.Min()
		a.True(found)
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
		v, found := tree.Delete(e.Key)
		a.True(found)
		a.Equal(i, v)
	}
	a.Equal(0, tree.Len())
}

func TestTreeDeleteMax(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	for i := 0; i < 128; i++ {
		e, found := tree.Max()
		a.True(found)
		a.Equal(127-i, e.Key)
		a.Equal(127-i, *e.Value)
		v, found := tree.Delete(e.Key)
		a.True(found)
		a.Equal(127-i, v)
	}
	a.Equal(0, tree.Len())
}

func TestTreeAt_WithCountChildren(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	for i := 0; i < 128; i++ {
		e := tree.At(i)
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
	}
	a.Panics(func() {
		tree.At(128)
	})
}

func TestTreeAt_WithoutCountChildren(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(false))
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	for i := 0; i < 128; i++ {
		e := tree.At(i)
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
	}
	a.Panics(func() {
		tree.At(128)
	})
}

func TestTreeDeleteAt(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i*2, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	for i := 64; i < 128; i++ {
		k, v := tree.DeleteAt(64)
		a.Equal(i*2, k)
		a.Equal(i, v)
	}
	for i := 0; i < 64; i++ {
		k, v := tree.DeleteAt(0)
		a.Equal(i*2, k)
		a.Equal(i, v)
	}
	a.Equal(0, tree.Len())
	a.Panics(func() {
		tree.DeleteAt(128)
	})
}

func TestTreeRandom(t *testing.T) {
	const count = 1024
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	data := make([]int, count)
	for i := 0; i < count; i++ {
		data[i] = i
	}
	for i := 0; i < 10; i++ {
		rand.Shuffle(len(data), func(i, j int) {
			data[i], data[j] = data[j], data[i]
		})
		for _, v := range data {
			ptr, inserted := tree.Insert(v, v)
			a.Equal(v, *ptr)
			a.True(inserted)
			if !a.NoError(checkHeightAndBalance(tree.root)) {
				tree.locate(v)
				fmt.Println(tree.Len())
				printTree(tree, os.Stdout)
				t.FailNow()
			}
		}
		for i, v := range data {
			val, deleted := tree.Delete(v)
			a.Equal(v, val)
			a.Truef(deleted, "key: %d, iter = %d", v, i)
			a.NoErrorf(checkHeightAndBalance(tree.root), "%d", i)
		}
		a.Equal(0, tree.Len())
	}
}

func TestTreeIterator(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	it := tree.AscendFromStart()
	for i := 0; ; i++ {
		e, ok := it.Next()
		if i == 128 {
			a.False(ok)
			break
		}
		a.True(ok)
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
	}
	for i := 127; ; i-- {
		e, ok := it.Prev()
		if i == -1 {
			a.False(ok)
			break
		}
		a.True(ok)
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
	}
}

func TestTreeIteratorValue(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	it := tree.AscendFromStart()
	for i := 0; ; i++ {
		e, ok := it.Value()
		if i == 128 {
			a.False(ok)
			break
		}
		a.True(ok)
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
		it.Next()
	}
	it = tree.DescendFromEnd()
	for i := 127; ; i-- {
		e, ok := it.Value()
		if i == -1 {
			a.False(ok)
			break
		}
		a.True(ok)
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
		it.Prev()
	}
}

func TestTreeAscend(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	for i := 0; i <= 100; i += 5 {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	it := tree.Ascend(-1)
	e, ok := it.Next()
	a.True(ok)
	a.Equal(0, e.Key)
	a.Equal(0, *e.Value)
	for i := 0; i <= 100; i++ {
		it = tree.Ascend(i)
		e, ok := it.Next()
		a.True(ok)
		if rem := i % 5; rem == 0 {
			a.Equal(i, e.Key)
			a.Equal(i, *e.Value)
		} else {
			a.Equal(i-rem+5, e.Key)
			a.Equal(i-rem+5, *e.Value)
		}
	}
	it = tree.Ascend(101)
	_, ok = it.Next()
	a.False(ok)
}

func TestTreeAscendAt(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	a.Panics(func() {
		tree.AscendAt(0)
	})
	for i := 0; i <= 100; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	for i := 0; i <= 100; i++ {
		it := tree.AscendAt(i)
		e, ok := it.Value()
		a.True(ok)
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
		for j := i - 1; j >= 0; j-- {
			it.Prev()
			e, ok = it.Value()
			a.True(ok)
			a.Equal(j, e.Key)
			a.Equal(j, *e.Value)
		}

		it = tree.AscendAt(i)
		for j := i + 1; j < tree.Len(); j++ {
			it.Next()
			e, ok = it.Value()
			a.True(ok)
			a.Equal(j, e.Key)
			a.Equal(j, *e.Value)
		}
	}
}

func TestTreeDescend(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	for i := 0; i <= 100; i += 5 {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	it := tree.Descend(101)
	e, ok := it.Next()
	a.True(ok)
	a.Equal(100, e.Key)
	a.Equal(100, *e.Value)
	for i := 0; i <= 100; i++ {
		it = tree.Descend(i)
		e, ok := it.Next()
		a.True(ok)
		if rem := i % 5; rem == 0 {
			a.Equal(i, e.Key)
			a.Equal(i, *e.Value)
		} else {
			a.Equal(i-rem, e.Key)
			a.Equal(i-rem, *e.Value)
		}
	}
	it = tree.Descend(-1)
	_, ok = it.Next()
	a.False(ok)
}

func checkHeightAndBalance[K, V any](l ptrLocation[K, V]) error {
	_, _, _, err := recalcHeightAndBalance(l)
	return err
}

func recalcHeightAndBalance[K, V any](l ptrLocation[K, V]) (height uint8, lCount, rCount uint32, err error) {
	if l.isNil() {
		return 0, 0, 0, nil
	}
	if !l.left().isNil() {
		lHeight, llCount, rrCount, err := recalcHeightAndBalance(l.left())
		if err != nil {
			return 0, 0, 0, err
		}
		height = 1 + lHeight
		lCount = llCount + rrCount + 1
	}
	if !l.right().isNil() {
		rHeight, rlCount, rrCount, err := recalcHeightAndBalance(l.right())
		if err != nil {
			return 0, 0, 0, err
		}
		height = max(height, 1+rHeight)
		rCount = rlCount + rrCount + 1
	}
	if height != l.height() {
		return 0, 0, 0, fmt.Errorf("invalid height for k=%v, v=%v, curr=%d, actual=%d", l.key(), *l.valuePtr(), l.height(), height)
	}
	if l.balance() < -1 || l.balance() > 1 {
		return 0, 0, 0, fmt.Errorf("invalid balance %d for k=%v, v=%v", l.balance(), l.key(), *l.valuePtr())
	}
	return height, lCount, rCount, nil
}

func printTree[K, V any, Cmp func(a, b K) int](t *Tree[K, V, Cmp], w io.Writer) {
	traverseTree(t, func(loc ptrLocation[K, V]) bool {
		_, _ = w.Write([]byte(loc.String()))
		_, _ = w.Write([]byte{'\n'})
		return true
	})
}

func traverseTree[K, V any, Cmp func(a, b K) int](t *Tree[K, V, Cmp], f func(loc ptrLocation[K, V]) bool) {
	if t.root.isNil() {
		return
	}
	traverseLocation(t.root, f)
}

func traverseLocation[K, V any](loc ptrLocation[K, V], f func(loc ptrLocation[K, V]) bool) {
	if !loc.left().isNil() {
		traverseLocation(loc.left(), f)
	}
	f(loc)
	if !loc.right().isNil() {
		traverseLocation(loc.right(), f)
	}
}
