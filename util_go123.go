//go:build go1.23

package goavl

import (
	"golang.org/x/exp/constraints"
)

func max2[T constraints.Ordered](a, b T) T {
	return max(a, b)
}

func min2[T constraints.Ordered](a, b T) T {
	return min(a, b)
}
