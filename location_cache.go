package goavl

import "sync"

var (
	_ locationCache[int, int] = (*basicLocationCache[int, int])(nil)
	_ locationCache[int, int] = (*pooledLocationCache[int, int])(nil)
)

type locationCache[K, V any] interface {
	new(k K, v V) ptrLocation[K, V]
	release(loc ptrLocation[K, V])
}

type basicLocationCache[K, V any] struct{}

func newBasicLocationCache[K, V any]() *basicLocationCache[K, V] {
	return &basicLocationCache[K, V]{}
}

func (lc *basicLocationCache[K, V]) new(k K, v V) ptrLocation[K, V] { //nolint:unused // used in locationCache iface
	pn := &ptrNode[K, V]{}
	pn.init(k, v)
	return ptrLocation[K, V]{
		ptrNode: pn,
	}
}

func (lc *basicLocationCache[K, V]) release(ptrLocation[K, V]) {} //nolint:unused // used in locationCache iface

type pooledLocationCache[K, V any] struct {
	p *sync.Pool
}

func newPooledLocationCache[K, V any](p *sync.Pool) *pooledLocationCache[K, V] {
	result := &pooledLocationCache[K, V]{}
	if p != nil && p.New == nil {
		result.p = p
	} else {
		result.p = &sync.Pool{
			New: func() any {
				return new(ptrNode[K, V])
			},
		}
	}
	return result
}

func (lc *pooledLocationCache[K, V]) new(k K, v V) ptrLocation[K, V] { //nolint:unused // used in locationCache iface
	pn, _ := lc.p.Get().(*ptrNode[K, V])
	if pn == nil {
		pn = &ptrNode[K, V]{}
	}
	pn.init(k, v)
	return ptrLocation[K, V]{
		ptrNode: pn,
	}
}

func (lc *pooledLocationCache[K, V]) release(loc ptrLocation[K, V]) { //nolint:unused // used in locationCache iface
	lc.p.Put(loc.ptrNode)
}
