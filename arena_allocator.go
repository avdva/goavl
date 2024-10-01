//go:build goexperiment.arenas

package goavl

import (
	"arena"
)

type arenaOptions struct {
	a *arena.Arena
}

// WithArena makes Tree use arenas (currently experimental) to allocate tree nodes.
// `a` cannot be nil and `a.Free` should be called when the tree is no longer in use.
func WithArena(a *arena.Arena) Option {
	return func(o *Options) {
		o.at = allocArenas
		o.ao.a = a
	}
}

type arenaLocationCache[K, V any] struct {
	a *arena.Arena
}

func newArenaLocationCache[K, V any](ao arenaOptions) *arenaLocationCache[K, V] {
	return &arenaLocationCache[K, V]{a: ao.a}
}

func (lc *arenaLocationCache[K, V]) new(k K, v V) ptrLocation[K, V] { //nolint:unused // used in locationCache iface
	pn := arena.New[ptrNode[K, V]](lc.a)
	pn.init(k, v)
	return ptrLocation[K, V]{
		ptrNode: pn,
	}
}

func (lc *arenaLocationCache[K, V]) release(ptrLocation[K, V]) {} //nolint:unused // used in locationCache iface
