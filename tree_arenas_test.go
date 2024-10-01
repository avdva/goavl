//go:build goexperiment.arenas

package goavl

import (
	"arena"
	"testing"
)

func TestTreeRandomArenas(t *testing.T) {
	a := arena.NewArena()
	defer a.Free()
	doTestTreeRandom(t, WithCountChildren(true), WithArena(a))
}
