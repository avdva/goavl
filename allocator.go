package goavl

type allocator[K, V any] interface {
	new(K, V) ptrLocation[K, V]
	free(l ptrLocation[K, V])
}

type simpleAllocator[K, V any] struct{}

func (a *simpleAllocator[K, V]) new(k K, v V) ptrLocation[K, V] {
	return makeLocation(k, v)
}

func (a *simpleAllocator[K, V]) free(l ptrLocation[K, V]) {}

type preAllocator[K, V any] struct {
	allocated []node[K, V]
}

func (a *preAllocator[K, V]) new(k K, v V) ptrLocation[K, V] {
	if len(a.allocated) == 0 {
		a.allocated = make([]node[K, V], 128)
	}
	result := ptrLocation[K, V]{
		ptr: &a.allocated[0],
	}
	result.ptr.k = k
	result.ptr.v = v
	a.allocated = a.allocated[1:]
	return result
}

func (a *preAllocator[K, V]) free(l ptrLocation[K, V]) {}

type idxNode[K, V any] struct {
	node[K, V]
	l, r, parent int
}

type arrayAllocator[K, V any] struct {
	data []idxNode[K, V]
}

func (aa *arrayAllocator[K, V]) new(k K, v V) idxLocation[K, V] {
	return idxLocation[K, V]{
		data: aa.data,
		idx:  0,
	}
}

type idxLocation[K, V any] struct {
	data []idxNode[K, V]
	idx  int
}
