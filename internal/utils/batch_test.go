package utils

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMin(t *testing.T) {
	require.Equal(t, 2, min(2,3))
	require.Equal(t, 2, min(3,2))
	require.Equal(t, -1, min(-1,2))
}

type testCaseBatch struct {
	src []int
	batchSize int
	expected [][]int
	expectingError bool
}

func doTestBatch(t *testing.T, currTestCase*testCaseBatch) {

	batches, err := SplitToBatches(currTestCase.src, currTestCase.batchSize)
	if err != nil {
		if  !currTestCase.expectingError {
			t.Fatalf("Got unexpected error %s", err)
		}
		return
	} else if currTestCase.expectingError {
		t.Fatalf("Expected error, but got nil")
	}

	t.Logf("Expected: %v, actual: %v", currTestCase.expected, batches)

	if len(currTestCase.expected) != len(batches) {
		t.Fatalf("Result size %d != %d", len(batches), len(currTestCase.expected))
	}

	if !reflect.DeepEqual(currTestCase.expected, batches) {
		t.Fatal("Result batches are not as expected")
	}
}


func TestSplitToBatches(t *testing.T) {

	testCases := map[string]testCaseBatch{
		"exact division": {src: []int{1,2,3,4,5,6}, batchSize: 3,  expected: [][]int{{1, 2, 3}, {4, 5, 6}}},
		"with remainder": {src: []int{1,2,3,4},     batchSize: 3,  expected: [][]int{{1, 2, 3}, {4}}},
		"single batch":   {src: []int{1,2,3,4},     batchSize: 10, expected: [][]int{{1, 2, 3, 4}}},
		"nil input":      {src: nil,                batchSize: 10, expected: [][]int{}},
		"batch size == 1":{src: []int{1,2,3,4},     batchSize: 1,  expected: [][]int{{1}, {2}, {3}, {4}}},
		"error if negative batch size": {src: []int{1,2}, batchSize: -2, expectingError: true},
		"error if zero batch size": {src: []int{1,2}, batchSize: 0, expectingError: true},
	}

	for name, currTest := range testCases {
		t.Run(name, func(t *testing.T){ doTestBatch(t, &currTest) })
	}
}
