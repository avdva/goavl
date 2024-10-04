package goavl

var _ locationCache[int, int] = (*basicLocationCache[int, int])(nil)

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
