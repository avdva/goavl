//go:build goexperiment.arenas

package extbench

import (
	"arena"
	"testing"

	"github.com/avdva/goavl"
)

func BenchmarkAVLInsertArenas(b *testing.B) {
	a := arena.NewArena()
	defer a.Free()
	doBenchmarkAVLInsert[int](b, goavl.WithArena(a))
}
