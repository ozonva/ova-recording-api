package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)


func slicesEqual(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

func TestSlicesEqual(t *testing.T) {
	require.True(t, slicesEqual([]int{1,2,3}, []int{1,2,3}))
	require.False(t, slicesEqual([]int{1,2,3}, []int{3,2,1}))
	require.False(t, slicesEqual([]int{1,2,3}, []int{1,2,3,4}))
	require.True(t, slicesEqual([]int{}, []int{}))
}

func TestBatchSimple(t *testing.T) {

	src := []int{1,2,3,4,5,6}

	batches := Batch(src, 3)

	require.Equal(t, 2, len(batches))
	require.True(t, slicesEqual(batches[0], []int{1,2,3}))
	require.True(t, slicesEqual(batches[1], []int{4,5,6}))

	src = []int{1,2,3,4}

	batches = Batch(src, 3)

	require.Equal(t, 2, len(batches))
	require.True(t, slicesEqual(batches[0], []int{1,2,3}))
	require.True(t, slicesEqual(batches[1], []int{4}))

	src = []int{1,2,3,4}

	batches = Batch(src, 10)

	require.Equal(t, 1, len(batches))
	require.True(t, slicesEqual(batches[0], src))
}

func TestBatchEdge(t *testing.T) {
	require.Panics(t, func() { Batch([]int{1,2,3}, -2)})
	require.Equal(t, 0, len(Batch(nil, 10)))
}
