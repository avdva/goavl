package goavl

type allocator[K, V any] interface {
	new(K, V) location[K, V]
	free(l location[K, V])
}

type simpleAllocator[K, V any] struct{}

func (a *simpleAllocator[K, V]) new(k K, v V) location[K, V] {
	return makeLocation(k, v)
}

func (a *simpleAllocator[K, V]) free(l location[K, V]) {}

type preAllocator[K, V any] struct {
	allocated []node[K, V]
}

func (a *preAllocator[K, V]) new(k K, v V) location[K, V] {
	if len(a.allocated) == 0 {
		a.allocated = make([]node[K, V], 128)
	}
	result := location[K, V]{
		ptr: &a.allocated[0],
	}
	result.ptr.k = k
	result.ptr.v = v
	a.allocated = a.allocated[1:]
	return result
}

func (a *preAllocator[K, V]) free(l location[K, V]) {}
