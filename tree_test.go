package goavl

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestEmptyTree(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	it := tree.IteratorAtFirst()
	e, ok := it.Next()
	a.Equal(0, e.Key)
	a.Equal((*int)(nil), e.Value)
	a.Equal(false, ok)
	v, ok := tree.Delete(0)
	a.Equal(0, v)
	a.Equal(false, ok)
	it = tree.IteratorAtLast()
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
		a.NoErrorf(checkHeightAndBalance(tree.root, tree.options.countChildren), "iter = %d", i)
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
		a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	}
	for i := 127; i >= 0; i-- {
		val, found := tree.Find(i)
		a.True(found)
		a.Equal(i*2, *val)
	}
}

func TestTreeInsertWithComparatorReturningMagnitude(t *testing.T) {
	a := assert.New(t)
	tree := New[int, int](func(a, b int) int {
		return a - b
	}, WithCountChildren(true))

	for _, i := range []int{10, 20, 5, 30, 1} {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
		a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	}

	a.Equal(5, tree.Len())
	for i, want := range []int{1, 5, 10, 20, 30} {
		e := tree.At(i)
		a.Equal(want, e.Key)
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
	a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))

	ptr, inserted = tree.Insert(0, 0)
	a.Equal(0, *ptr)
	a.True(inserted)
	ptr, inserted = tree.Insert(-1, -1)
	a.Equal(-1, *ptr)
	a.True(inserted)
	a.Equal(2, tree.Len())
	a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
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
	a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	v, deleted = tree.Delete(1)
	a.True(deleted)
	a.Equal(1, v)
	a.Equal(1, tree.Len())
	_, deleted = tree.Delete(-1)
	a.False(deleted)
	a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	v, deleted = tree.Delete(0)
	a.True(deleted)
	a.Equal(0, v)
	a.Equal(0, tree.Len())
	a.True(tree.root.isNil())
	a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))

	ptr, inserted = tree.Insert(0, 0)
	a.Equal(0, *ptr)
	a.True(inserted)
	ptr, inserted = tree.Insert(1, 1)
	a.Equal(1, *ptr)
	a.True(inserted)
	a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	v, deleted = tree.Delete(0)
	a.True(deleted)
	a.Equal(0, v)
	a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	a.Equal(1, tree.Len())
	v, deleted = tree.Delete(1)
	a.True(deleted)
	a.Equal(1, v)
	a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	a.True(tree.root.isNil())
	a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))

	for i := 128; i >= 0; i-- {
		ptr, inserted = tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
		a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	}
	for i := 128; i >= 0; i-- {
		v, deleted = tree.Delete(i)
		a.True(deleted)
		a.Equal(i, v)
		a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	}
	a.Equal(0, tree.Len())
}

func TestTreeUpdateKey(t *testing.T) {
	a := assert.New(t)

	t.Run("missing old key", func(t *testing.T) {
		tree := NewComparable[int, int](WithCountChildren(true))
		tree.Insert(1, 10)

		ptr, updated := tree.UpdateKey(2, 3)

		a.Nil(ptr)
		a.False(updated)
		a.Equal(1, tree.Len())
		assertTreeKeys(t, tree, []int{1})
	})

	t.Run("equivalent key", func(t *testing.T) {
		tree := New[int, int](func(a, b int) int {
			return (a / 10) - (b / 10)
		}, WithCountChildren(true))
		tree.Insert(12, 100)

		ptr, updated := tree.UpdateKey(12, 18)

		a.True(updated)
		a.Equal(100, *ptr)
		a.Equal(1, tree.Len())
		entry, found := tree.Min()
		a.True(found)
		a.Equal(18, entry.Key)
		a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren))
	})

	t.Run("existing new key", func(t *testing.T) {
		tree := NewComparable[int, int](WithCountChildren(true))
		for _, k := range []int{1, 3, 5, 7} {
			tree.Insert(k, k*10)
		}

		ptr, updated := tree.UpdateKey(3, 5)

		a.True(updated)
		a.Equal(30, *ptr)
		a.Equal(3, tree.Len())
		_, found := tree.Find(3)
		a.False(found)
		got, found := tree.Find(5)
		a.True(found)
		a.Equal(30, *got)
		assertTreeKeys(t, tree, []int{1, 5, 7})
	})

	t.Run("in place", func(t *testing.T) {
		tree := NewComparable[int, int](WithCountChildren(true))
		for _, k := range []int{1, 3, 5} {
			tree.Insert(k, k*10)
		}

		ptr, updated := tree.UpdateKey(3, 4)

		a.True(updated)
		a.Equal(30, *ptr)
		a.Equal(3, tree.Len())
		_, found := tree.Find(3)
		a.False(found)
		got, found := tree.Find(4)
		a.True(found)
		a.Equal(30, *got)
		assertTreeKeys(t, tree, []int{1, 4, 5})
	})

	t.Run("generic path", func(t *testing.T) {
		tree := NewComparable[int, int](WithCountChildren(true))
		for i := 1; i <= 8; i++ {
			tree.Insert(i, i*10)
		}

		ptr, updated := tree.UpdateKey(2, 9)

		a.True(updated)
		a.Equal(20, *ptr)
		a.Equal(8, tree.Len())
		_, found := tree.Find(2)
		a.False(found)
		got, found := tree.Find(9)
		a.True(found)
		a.Equal(20, *got)
		assertTreeKeys(t, tree, []int{1, 3, 4, 5, 6, 7, 8, 9})
	})
}

func TestTreeUpdateKeyRandom(t *testing.T) {
	const count = 128
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	keys := make([]int, count)
	for i := 0; i < count; i++ {
		keys[i] = i
		tree.Insert(i, i*10)
	}
	rand.New(rand.NewSource(1)).Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for i, k := range keys {
		newKey := k + count
		ptr, updated := tree.UpdateKey(k, newKey)
		a.Truef(updated, "key: %d, iter = %d", k, i)
		a.Equal(k*10, *ptr)
		_, found := tree.Find(k)
		a.False(found)
		got, found := tree.Find(newKey)
		a.True(found)
		a.Equal(k*10, *got)
		a.NoErrorf(checkHeightAndBalance(tree.root, tree.options.countChildren), "key: %d, iter = %d", k, i)
	}

	a.Equal(count, tree.Len())
	for i := 0; i < count; i++ {
		entry := tree.At(i)
		a.Equal(i+count, entry.Key)
		a.Equal(i*10, *entry.Value)
	}
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
	doTestTreeRandom(t, WithCountChildren(true))
}

func TestTreeRandomSyncPool(t *testing.T) {
	doTestTreeRandom(t, WithSyncPoolAllocator(true))
}

func TestTreeRandomSyncPoolCustom(t *testing.T) {
	var sp sync.Pool
	doTestTreeRandom(t, WithSyncPool(&sp))
}

func doTestTreeRandom(t *testing.T, opts ...Option) {
	const count = 1024
	a := assert.New(t)
	tree := NewComparable[int, int](opts...)
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
			if !a.NoError(checkHeightAndBalance(tree.root, tree.options.countChildren)) {
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
			a.NoErrorf(checkHeightAndBalance(tree.root, tree.options.countChildren), "%d", i)
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
	it := tree.IteratorAtFirst()
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

func TestTreeDeleteIterator(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	it := tree.IteratorAtFirst()
	_, ok := it.Value()
	a.False(ok)

	it = tree.DeleteIterator(it)
	_, ok = it.Value()
	a.False(ok)

	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	it = tree.IteratorAtFirst()
	for i := 0; i < 128; i++ {
		e, ok := it.Value()
		a.True(ok)
		a.Equal(i, e.Key)
		a.Equal(i, *e.Value)
		it = tree.DeleteIterator(it)
	}
	_, ok = it.Value()
	a.False(ok)
	a.Zero(tree.Len())
}

func TestTreeDeleteIterator2(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	it := tree.IteratorAtFirst()
	// delete all even keys
	for i := 0; i < 64; i++ {
		e, ok := it.Value()
		a.True(ok)
		a.Equal(i*2, e.Key)
		a.Equal(i*2, *e.Value)
		it = tree.DeleteIterator(it)

		e, ok = it.Next()
		a.True(ok)
		a.Equal(i*2+1, e.Key)
		a.Equal(i*2+1, *e.Value)
	}
	// delete all odd keys
	it = tree.IteratorAtFirst()
	for i := 0; i < 64; i++ {
		e, ok := it.Value()
		a.True(ok)
		a.Equal(i*2+1, e.Key)
		a.Equal(i*2+1, *e.Value)
		it = tree.DeleteIterator(it)
	}
	_, ok := it.Value()
	a.False(ok)
	a.Zero(tree.Len())
}

func TestTreeDeleteIteratorRejectsStaleIterator(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	tree.Insert(1, 1)
	tree.Insert(2, 2)

	it := tree.IteratorAtFirst()
	next := tree.DeleteIterator(it)
	a.Equal(1, tree.Len())

	invalid := tree.DeleteIterator(it)
	_, ok := invalid.Value()
	a.False(ok)
	a.Equal(1, tree.Len())

	e, ok := next.Value()
	a.True(ok)
	a.Equal(2, e.Key)
}

func TestTreeDeleteIteratorRejectsReusedStaleIterator(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithSyncPoolAllocator(true))
	tree.Insert(1, 1)

	it := tree.IteratorAtFirst()
	tree.DeleteIterator(it)
	tree.Insert(2, 2)

	tree.DeleteIterator(it)
	a.Equal(1, tree.Len())
	e, ok := tree.Min()
	a.True(ok)
	a.Equal(2, e.Key)
}

func TestTreeDeleteIteratorRejectsOtherTreeIterator(t *testing.T) {
	a := assert.New(t)
	tree1 := NewComparable[int, int]()
	tree2 := NewComparable[int, int]()
	tree1.Insert(1, 1)
	tree2.Insert(2, 2)

	it := tree1.IteratorAtFirst()
	invalid := tree2.DeleteIterator(it)

	_, ok := invalid.Value()
	a.False(ok)
	a.Equal(1, tree1.Len())
	a.Equal(1, tree2.Len())
	e, ok := tree1.Min()
	a.True(ok)
	a.Equal(1, e.Key)
	e, ok = tree2.Min()
	a.True(ok)
	a.Equal(2, e.Key)
}

func TestTreeIteratorValue(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	it := tree.IteratorAtFirst()
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
	it = tree.IteratorAtLast()
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

func TestTreeLowerBound(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	for i := 0; i <= 100; i += 5 {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	it := tree.LowerBound(-1)
	e, ok := it.Next()
	a.True(ok)
	a.Equal(0, e.Key)
	a.Equal(0, *e.Value)
	for i := 0; i <= 100; i++ {
		it = tree.LowerBound(i)
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
	it = tree.LowerBound(101)
	_, ok = it.Next()
	a.False(ok)
}

func TestTreeIteratorAt(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	a.Panics(func() {
		tree.IteratorAt(0)
	})
	for i := 0; i <= 100; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	for i := 0; i <= 100; i++ {
		it := tree.IteratorAt(i)
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

		it = tree.IteratorAt(i)
		for j := i + 1; j < tree.Len(); j++ {
			it.Next()
			e, ok = it.Value()
			a.True(ok)
			a.Equal(j, e.Key)
			a.Equal(j, *e.Value)
		}
	}
}

func TestTreeFloor(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	for i := 0; i <= 100; i += 5 {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}
	it := tree.Floor(101)
	e, ok := it.Next()
	a.True(ok)
	a.Equal(100, e.Key)
	a.Equal(100, *e.Value)
	for i := 0; i <= 100; i++ {
		it = tree.Floor(i)
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
	it = tree.Floor(-1)
	_, ok = it.Next()
	a.False(ok)
}

func TestTreeNewIteratorNames(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	a.Panics(func() {
		tree.IteratorAt(0)
	})
	for i := 0; i <= 100; i++ {
		ptr, inserted := tree.Insert(i, i)
		a.Equal(i, *ptr)
		a.True(inserted)
	}

	it := tree.IteratorAtFirst()
	e, ok := it.Value()
	a.True(ok)
	a.Equal(0, e.Key)

	it = tree.IteratorAtLast()
	e, ok = it.Value()
	a.True(ok)
	a.Equal(100, e.Key)

	for i := 0; i <= 100; i++ {
		it = tree.IteratorAt(i)
		e, ok = it.Value()
		a.True(ok)
		a.Equal(i, e.Key)
	}
}

func TestTreeBounds(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int](WithCountChildren(true))
	for _, key := range []int{0, 10, 20, 30, 40} {
		ptr, inserted := tree.Insert(key, key)
		a.Equal(key, *ptr)
		a.True(inserted)
	}

	it := tree.LowerBound(5)
	e, ok := it.Value()
	a.True(ok)
	a.Equal(10, e.Key)
	it = tree.LowerBound(10)
	e, ok = it.Value()
	a.True(ok)
	a.Equal(10, e.Key)
	it = tree.LowerBound(41)
	e, ok = it.Value()
	a.False(ok)
	a.Equal(0, e.Key)

	it = tree.UpperBound(5)
	e, ok = it.Value()
	a.True(ok)
	a.Equal(10, e.Key)
	it = tree.UpperBound(10)
	e, ok = it.Value()
	a.True(ok)
	a.Equal(20, e.Key)
	it = tree.UpperBound(40)
	e, ok = it.Value()
	a.False(ok)
	a.Equal(0, e.Key)

	it = tree.Floor(25)
	e, ok = it.Value()
	a.True(ok)
	a.Equal(20, e.Key)
	it = tree.Floor(20)
	e, ok = it.Value()
	a.True(ok)
	a.Equal(20, e.Key)
	it = tree.Floor(-1)
	e, ok = it.Value()
	a.False(ok)
	a.Equal(0, e.Key)

	it = tree.LowerBound(10)
	e, ok = it.Value()
	a.True(ok)
	a.Equal(10, e.Key)
	it = tree.Floor(25)
	e, ok = it.Value()
	a.True(ok)
	a.Equal(20, e.Key)
}

func sortedRank(keys []int, key int) (rank int, found bool) {
	for i, candidate := range keys {
		switch cmp := intCmp(key, candidate); {
		case cmp < 0:
			return 0, false
		case cmp == 0:
			return i, true
		}
	}
	return 0, false
}

func sortedLowerBoundRank(keys []int, key int) (rank int, found bool) {
	for i, candidate := range keys {
		if intCmp(key, candidate) <= 0 {
			return i, true
		}
	}
	return 0, false
}

func sortedFloorRank(keys []int, key int) (rank int, found bool) {
	for i, candidate := range keys {
		if intCmp(key, candidate) < 0 {
			if i == 0 {
				return 0, false
			}
			return i - 1, true
		}
	}
	if len(keys) == 0 {
		return 0, false
	}
	return len(keys) - 1, true
}

func sortedUpperBoundRank(keys []int, key int) (rank int, found bool) {
	for i, candidate := range keys {
		if intCmp(key, candidate) < 0 {
			return i, true
		}
	}
	return 0, false
}

func sortedRankDistance(keys []int, k1 int, k2 int) (distance int, found bool) {
	r1, found := sortedRank(keys, k1)
	if !found {
		return 0, false
	}
	r2, found := sortedRank(keys, k2)
	if !found {
		return 0, false
	}
	if r2 >= r1 {
		return r2 - r1, true
	}
	return r1 - r2, true
}

func sortedCountInRange(keys []int, k1 int, k2 int) int {
	r1, found := sortedLowerBoundRank(keys, k1)
	if !found {
		return 0
	}
	r2, found := sortedFloorRank(keys, k2)
	if !found {
		return 0
	}
	if r2 >= r1 {
		return r2 - r1 + 1
	}
	return 0
}

func assertOptionalEntryKey(t *testing.T, expected int, expectedFound bool, entry Entry[int, int], found bool) {
	t.Helper()
	a := assert.New(t)
	a.Equal(expectedFound, found)
	if expectedFound {
		a.Equal(expected, entry.Key)
		a.Equal(expected, *entry.Value)
	} else {
		a.Equal(0, entry.Key)
		a.Nil(entry.Value)
	}
}

func testTreeRankRangeAndBoundsAgainstSortedSlice(t *testing.T, opts ...Option) {
	t.Helper()
	a := assert.New(t)
	tree := NewComparable[int, int](opts...)

	sortedKeys := []int{-50, -10, 0, 3, 4, 10, 17, 31, 32, 99}
	insertKeys := append([]int(nil), sortedKeys...)
	rand.New(rand.NewSource(0x5eed)).Shuffle(len(insertKeys), func(i, j int) {
		insertKeys[i], insertKeys[j] = insertKeys[j], insertKeys[i]
	})

	for _, key := range insertKeys {
		ptr, inserted := tree.Insert(key, key)
		a.Equal(key, *ptr)
		a.True(inserted)
	}

	for rank, key := range sortedKeys {
		actualRank, found := tree.Rank(key)
		a.True(found)
		a.Equal(rank, actualRank)
		a.Equal(key, tree.At(rank).Key)
		it := tree.IteratorAt(rank)
		entry, ok := it.Value()
		a.True(ok)
		a.Equal(key, entry.Key)
	}

	queryKeys := []int{-60, -50, -49, -11, -10, -9, -1, 0, 1, 3, 4, 5, 10, 16, 17, 18, 30, 31, 32, 33, 98, 99, 100}
	for _, key := range queryKeys {
		expectedRank, expectedFound := sortedRank(sortedKeys, key)
		actualRank, actualFound := tree.Rank(key)
		a.Equal(expectedFound, actualFound)
		if expectedFound {
			a.Equal(expectedRank, actualRank)
		}

		lowerRank, lowerFound := sortedLowerBoundRank(sortedKeys, key)
		lowerIt := tree.LowerBound(key)
		lowerEntry, lowerOK := lowerIt.Value()
		if lowerFound {
			assertOptionalEntryKey(t, sortedKeys[lowerRank], true, lowerEntry, lowerOK)
		} else {
			assertOptionalEntryKey(t, 0, false, lowerEntry, lowerOK)
		}

		upperRank, upperFound := sortedUpperBoundRank(sortedKeys, key)
		upperIt := tree.UpperBound(key)
		upperEntry, upperOK := upperIt.Value()
		if upperFound {
			assertOptionalEntryKey(t, sortedKeys[upperRank], true, upperEntry, upperOK)
		} else {
			assertOptionalEntryKey(t, 0, false, upperEntry, upperOK)
		}
	}

	for _, k1 := range queryKeys {
		for _, k2 := range queryKeys {
			expectedDistance, expectedFound := sortedRankDistance(sortedKeys, k1, k2)
			actualDistance, actualFound := tree.RankDistance(k1, k2)
			a.Equal(expectedFound, actualFound)
			if expectedFound {
				a.Equal(expectedDistance, actualDistance)
			}
			a.Equal(sortedCountInRange(sortedKeys, k1, k2), tree.CountInRange(k1, k2))
		}
	}
}

func TestTreeRankRangeAndBoundsAgainstSortedSlice(t *testing.T) {
	t.Run("basic without counts", func(t *testing.T) {
		testTreeRankRangeAndBoundsAgainstSortedSlice(t, WithCountChildren(false))
	})
	t.Run("basic with counts", func(t *testing.T) {
		testTreeRankRangeAndBoundsAgainstSortedSlice(t, WithCountChildren(true))
	})
	t.Run("sync pool without counts", func(t *testing.T) {
		testTreeRankRangeAndBoundsAgainstSortedSlice(t, WithCountChildren(false), WithSyncPool(nil))
	})
	t.Run("sync pool with counts", func(t *testing.T) {
		testTreeRankRangeAndBoundsAgainstSortedSlice(t, WithCountChildren(true), WithSyncPool(nil))
	})
}

func checkHeightAndBalance[K, V any](l location[K, V], checkCounts bool) error {
	_, _, _, err := recalcHeightAndBalance(l, checkCounts)
	return err
}

func recalcHeightAndBalance[K, V any](l location[K, V], checkCounts bool) (height uint8, lCount, rCount uint32, err error) {
	if l.isNil() {
		return 0, 0, 0, nil
	}
	if !l.left().isNil() {
		lHeight, llCount, rrCount, err := recalcHeightAndBalance(l.left(), checkCounts)
		if err != nil {
			return 0, 0, 0, err
		}
		height = 1 + lHeight
		lCount = llCount + rrCount + 1
	}
	if !l.right().isNil() {
		rHeight, rlCount, rrCount, err := recalcHeightAndBalance(l.right(), checkCounts)
		if err != nil {
			return 0, 0, 0, err
		}
		height = max2(height, 1+rHeight)
		rCount = rlCount + rrCount + 1
	}
	if height != l.height() {
		return 0, 0, 0, fmt.Errorf("invalid height for k=%v, v=%v, curr=%d, actual=%d", l.key(), *l.valuePtr(), l.height(), height)
	}
	if l.balance() < -1 || l.balance() > 1 {
		return 0, 0, 0, fmt.Errorf("invalid balance %d for k=%v, v=%v", l.balance(), l.key(), *l.valuePtr())
	}
	if count := lCount + rCount; checkCounts && count != l.childrenCount() {
		return 0, 0, 0, fmt.Errorf("invalid children count for k=%v, v=%v, curr=%d, actual=%d", l.key(), *l.valuePtr(), l.childrenCount(), count)
	}
	return height, lCount, rCount, nil
}

func printTree[K, V any, Cmp func(a, b K) int](t *Tree[K, V, Cmp], w io.Writer) {
	traverseTree(t, func(loc location[K, V]) bool {
		_, _ = w.Write([]byte(loc.String()))
		_, _ = w.Write([]byte{'\n'})
		return true
	})
}

func traverseTree[K, V any, Cmp func(a, b K) int](t *Tree[K, V, Cmp], f func(loc location[K, V]) bool) {
	if t.root.isNil() {
		return
	}
	traverseLocation(t.root, f)
}

func traverseLocation[K, V any](loc location[K, V], f func(loc location[K, V]) bool) {
	if !loc.left().isNil() {
		traverseLocation(loc.left(), f)
	}
	f(loc)
	if !loc.right().isNil() {
		traverseLocation(loc.right(), f)
	}
}

func assertTreeKeys[V any, Cmp func(a, b int) int](t *testing.T, tree *Tree[int, V, Cmp], want []int) {
	t.Helper()
	assert.NoError(t, checkHeightAndBalance(tree.root, tree.options.countChildren))
	assert.Equal(t, len(want), tree.Len())
	for i, key := range want {
		entry := tree.At(i)
		assert.Equal(t, key, entry.Key)
	}
}
