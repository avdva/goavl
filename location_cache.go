package goavl

type locationCache[K, V any] interface {
	new(k K, v V) ptrLocation[K, V]
	release(loc ptrLocation[K, V])
}
