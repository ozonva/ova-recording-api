package utils

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	require.True(t, contains(1, []int{1,2,3}))
	require.True(t, contains(2, []int{1,2,3}))
	require.False(t, contains(5, []int{1,2,3}))
	require.False(t, contains(5, []int{1,2,3,6,10}))
	require.False(t, contains(1, nil))
}

type testCaseFilter struct {
	src      []int
	target   []int
	expected []int
}

func doTestFilter(t *testing.T, currTestCase* testCaseFilter) {

	result := FilterBy(currTestCase.src, currTestCase.target)
	t.Logf("Expected: %v, actual: %v", currTestCase.expected, result)

	if len(currTestCase.expected) != len(result) {
		t.Fatalf("Result size %d != %d", len(result), len(currTestCase.expected))
	}

	if !reflect.DeepEqual(currTestCase.expected, result) {
		t.Fatal("Result of the reversion is not as expected")
	}
}

func TestFilter(t *testing.T) {

	simpleFilterTarget := []int{1,2,3}

	testCases := map[string]testCaseFilter{
		"simple":             {src:[]int{1,2,3,4,5,6},       target: simpleFilterTarget, expected: []int{1,2,3}},
		"unsorted and mixed": {src:[]int{3,1,5,2,4,3,4,5,6}, target: simpleFilterTarget, expected: []int{3,1,2,3}},
		"none matched":       {src:[]int{4,5,6},             target: simpleFilterTarget, expected: []int{}},
		"empty input":        {src:[]int{},                  target: simpleFilterTarget, expected: []int{}},
		"empty target":       {src:[]int{1,2,3},             target: []int{},            expected: []int{}},
		"nil arguments":      {src:nil,                      target: nil,                expected: []int{}},
	}

	for name, currTest := range testCases {
		t.Run(name, func(t *testing.T){ doTestFilter(t, &currTest) })
	}
}
