package goavl

import (
	"math/bits"

	"golang.org/x/exp/constraints"
)

func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func log2(a uint64) uint64 {
	return uint64(63 - bits.LeadingZeros64(a))
}
