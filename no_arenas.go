//go:build !goexperiment.arenas

package goavl

type arenaOptions struct{}

func newArenaLocationCache[K, V any](ao arenaOptions) locationCache[K, V] {
	panic("unreachable")
}
