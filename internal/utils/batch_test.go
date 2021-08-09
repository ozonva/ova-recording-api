package utils

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type testCaseBatch struct {
	src []int
	batchSize int
	expected [][]int
	panics bool
}

func doTestBatch(t *testing.T, currTestCase*testCaseBatch) {

	if currTestCase.panics {
		require.Panics(t, func() {SplitToBatches(currTestCase.src, currTestCase.batchSize)})
		return
	}
	batches := SplitToBatches(currTestCase.src, currTestCase.batchSize)
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
		"panic if negative batch size": {src: []int{1,2}, batchSize: -2, panics: true},
	}

	for name, currTest := range testCases {
		t.Run(name, func(t *testing.T){ doTestBatch(t, &currTest) })
	}
}
