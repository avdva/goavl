//go:build go1.23

package goavl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeIteratorGo123(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i*2)
		a.Equal(i*2, *ptr)
		a.True(inserted)
	}
	i := 0
	for k, v := range tree.All() {
		a.Equal(k, i)
		a.Equal(v, i*2)
		i++
	}
	a.Equal(128, i)
}

func TestTreeMutIteratorGo123(t *testing.T) {
	a := assert.New(t)
	tree := NewComparable[int, int]()
	for i := 0; i < 128; i++ {
		ptr, inserted := tree.Insert(i, i*2)
		a.Equal(i*2, *ptr)
		a.True(inserted)
	}

	i := 0
	for m := range tree.AllMut() {
		a.Equal(m.E.Key, i)
		a.Equal(*m.E.Value, i*2)
		*m.E.Value = i * 4
		i++
	}
	a.Equal(128, i)

	i = 0
	for m := range tree.AllMut() {
		a.Equal(m.E.Key, i)
		a.Equal(*m.E.Value, i*4)
		m.Delete()
		m.Delete()
		a.Equal(m.E.Key, i)
		a.Equal(*m.E.Value, i*4)
		i++
	}
	a.Equal(128, i)
	a.Zero(tree.Len())
}
