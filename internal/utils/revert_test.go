package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRevert(t *testing.T) {
	src := map[string]int{"one": 1, "two": 2}
	dst := Revert(&src)

	for value, key := range dst {
		srcValue, ok := src[key]
		require.True(t, ok)
		require.Equal(t, value, srcValue)
	}

	require.Equal(t, 0, len(Revert(&map[string]int{})))

	require.Panics(t, func() { Revert(nil)})
	require.Panics(t, func() { Revert(&map[string]int{"one": 1, "two": 2, "duplicate":1})})
}
