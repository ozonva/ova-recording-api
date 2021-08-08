package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContains(t *testing.T) {
	require.True(t, contains(1, []int{1,2,3}))
	require.True(t, contains(2, []int{1,2,3}))
	require.False(t, contains(5, []int{1,2,3}))
	require.False(t, contains(5, []int{1,2,3,6,10}))
}

func TestFilter(t *testing.T) {
	require.True(t, slicesEqual(Filter([]int{1,2,3,4,5,6}), []int{1,2,3}))
	require.True(t, slicesEqual(Filter([]int{3,1,5,2,4,3,4,5,6}), []int{3,1,2,3}))
	require.True(t, slicesEqual(Filter([]int{4,5,6}), []int{}))
	require.True(t, slicesEqual(Filter([]int{}), []int{}))
}
