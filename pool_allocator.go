package goavl

import "sync"

var _ locationCache[int, int] = (*pooledLocationCache[int, int])(nil)

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

func (lc *pooledLocationCache[K, V]) new(k K, v V) location[K, V] { //nolint:unused // used in locationCache iface
	pn, _ := lc.p.Get().(*ptrNode[K, V])
	if pn == nil {
		pn = &ptrNode[K, V]{}
	}
	pn.init(k, v)
	return location[K, V]{
		ptrNode: pn,
	}
}

func (lc *pooledLocationCache[K, V]) release(loc location[K, V]) { //nolint:unused // used in locationCache iface
	pn := loc.ptrNode
	*pn = ptrNode[K, V]{}
	lc.p.Put(pn)
}
