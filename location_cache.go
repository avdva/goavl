package goavl

import "sync"

var (
	_ locationCache[int, int] = (*simpleLocationCache[int, int])(nil)
	_ locationCache[int, int] = (*pooledLocationCache[int, int])(nil)
)

type locationCache[K, V any] interface {
	new(k K, v V) ptrLocation[K, V]
	release(loc ptrLocation[K, V])
}

type simpleLocationCache[K, V any] struct{}

func newSimpleLocationCache[K, V any]() *simpleLocationCache[K, V] {
	return &simpleLocationCache[K, V]{}
}

func (lc *simpleLocationCache[K, V]) new(k K, v V) ptrLocation[K, V] { //nolint:unused // used in locationCache iface
	pn := &ptrNode[K, V]{}
	pn.init(k, v)
	return ptrLocation[K, V]{
		ptrNode: pn,
	}
}

type pooledLocationCache[K, V any] struct {
	p sync.Pool
}

func (lc *simpleLocationCache[K, V]) release(ptrLocation[K, V]) {} //nolint:unused // used in locationCache iface

func newPooledLocationCache[K, V any]() *pooledLocationCache[K, V] {
	return &pooledLocationCache[K, V]{
		p: sync.Pool{
			New: func() any {
				return new(ptrNode[K, V])
			},
		},
	}
}

func (lc *pooledLocationCache[K, V]) new(k K, v V) ptrLocation[K, V] { //nolint:unused // used in locationCache iface
	pn, _ := lc.p.Get().(*ptrNode[K, V])
	pn.init(k, v)
	return ptrLocation[K, V]{
		ptrNode: pn,
	}
}

func (lc *pooledLocationCache[K, V]) release(loc ptrLocation[K, V]) { //nolint:unused // used in locationCache iface
	lc.p.Put(loc.ptrNode)
}
