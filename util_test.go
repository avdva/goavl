package goavl

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog2(t *testing.T) {
	tests := []struct {
		val uint64
		exp uint64
	}{
		{
			val: 2,
			exp: 1,
		},
		{
			val: 32,
			exp: 5,
		},
		{
			val: 33,
			exp: 5,
		},
		{
			val: 1 << 63,
			exp: 63,
		},
	}
	for i := range tests {
		i := i
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assert.New(t)
			a.Equal(tests[i].exp, log2(tests[i].val))
		})
	}
}
