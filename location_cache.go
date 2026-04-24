package goavl

type locationCache[K, V any] interface {
	new(k K, v V) location[K, V]
	release(loc location[K, V])
}
