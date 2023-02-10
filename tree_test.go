package goavl

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeInsert(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		a.Truef(tree.Insert(i, i), "k: %v", i)
		mink, minv, found := tree.Min()
		a.Equal(0, mink)
		a.Equal(0, minv)
		a.True(found)

		maxk, maxv, found := tree.Max()
		a.Equal(i, maxk)
		a.Equal(i, maxv)
		a.True(found)
		a.NoErrorf(checkHeightAndBalance(tree.root), "iter = %d", i)
	}
	for i := 0; i < 128; i++ {
		val, found := tree.Find(i)
		a.True(found)
		a.Equal(i, val)
	}

	for i := 127; i >= 0; i-- {
		a.Falsef(tree.Insert(i, i*2), "k: %v", i)
		a.NoError(checkHeightAndBalance(tree.root))
	}
	for i := 127; i >= 0; i-- {
		val, found := tree.Find(i)
		a.True(found)
		a.Equal(i*2, val)
	}
}

func TestTreeDelete(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	a.Equal(0, tree.Len())

	a.True(tree.Insert(0, 0))
	v, deleted := tree.Delete(0)
	a.True(deleted)
	a.Equal(0, v)
	a.Equal(0, tree.Len())
	a.True(tree.root.isNil())
	a.NoError(checkHeightAndBalance(tree.root))

	a.True(tree.Insert(0, 0))
	a.True(tree.Insert(-1, -1))
	a.Equal(2, tree.Len())
	a.NoError(checkHeightAndBalance(tree.root))
	v, deleted = tree.Delete(0)
	a.True(deleted)
	a.Equal(0, v)
	v, deleted = tree.Delete(-1)
	a.Equal(-1, v)
	a.True(deleted)

	a.True(tree.Insert(0, 0))
	a.True(tree.Insert(1, 1))
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

	a.True(tree.Insert(0, 0))
	a.True(tree.Insert(1, 1))
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
		a.True(tree.Insert(i, i))
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

func TestTreeAt(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		a.Truef(tree.Insert(i, i), "k: %v", i)
	}
	for i := 0; i < 128; i++ {
		k, v := tree.At(i)
		a.Equal(i, k)
		a.Equal(i, v)
	}
	a.Panics(func() {
		tree.At(128)
	})
	a.Panics(func() {
		tree := NewComparable[int, int](WithCountChildren(false))
		tree.Insert(0, 0)
		tree.At(0)
	})
}

func TestTreeDeleteAt(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		a.Truef(tree.Insert(i, i), "k: %v", i)
	}
	for i := 64; i < 128; i++ {
		a.Equal(i, tree.DeleteAt(64))
	}
	for i := 0; i < 64; i++ {
		a.Equal(i, tree.DeleteAt(0))
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
			a.True(tree.Insert(v, v))
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
		a.Truef(tree.Insert(i, i), "k: %v", i)
	}
	it := tree.ForwardIterator()
	for i := 0; ; i++ {
		k, v, ok := it.Next()
		if i == 128 {
			a.False(ok)
			break
		}
		a.True(ok)
		a.Equal(i, k)
		a.Equal(i, v)
	}
	rev := tree.ReverseIterator()
	for i := 127; ; i-- {
		k, v, ok := rev.Next()
		if i == -1 {
			a.False(ok)
			break
		}
		a.True(ok)
		a.Equal(i, k)
		a.Equal(i, v)
	}
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
		height = uint8Max(height, 1+rHeight)
		rCount = rlCount + rrCount + 1
	}
	if height != l.height() {
		return 0, 0, 0, fmt.Errorf("invalid height for k=%v, v=%v, curr=%d, actual=%d", l.key(), l.value(), l.height(), height)
	}
	if l.balance() < -1 || l.balance() > 1 {
		return 0, 0, 0, fmt.Errorf("invalid balance %d for k=%v, v=%v", l.balance(), l.key(), l.value())
	}
	if lCount != l.leftCount() {
		return 0, 0, 0, fmt.Errorf("invalid left node count for k=%v, v=%v, curr=%d, actual=%d", l.key(), l.value(), l.leftCount(), lCount)
	}
	if rCount != l.rightCount() {
		return 0, 0, 0, fmt.Errorf("invalid right node count for k=%v, v=%v, curr=%d, actual=%d", l.key(), l.value(), l.rightCount(), rCount)
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
