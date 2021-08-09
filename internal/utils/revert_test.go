package utils

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type testCaseRevert struct {
	src map[string]int
	expected map[int]string
	panics bool
}

func doTestRevert(t *testing.T, currTestCase* testCaseRevert) {
	if currTestCase.panics {
		require.Panics(t, func() {Revert(currTestCase.src)})
		return
	}

	reverse := Revert(currTestCase.src)
	t.Logf("Expected: %v, actual: %v", currTestCase.expected, reverse)

	if len(currTestCase.expected) != len(reverse) {
		t.Fatalf("Map size %d != %d", len(reverse), len(currTestCase.expected))
	}

	if !reflect.DeepEqual(currTestCase.expected, reverse) {
		t.Fatal("Result of the reversion is not as expected")
	}
}


func TestRevert(t *testing.T) {

	testCases := map[string]testCaseRevert{
		"simple":      {src: map[string]int{"one": 1, "two": 2}, expected: map[int]string{1:"one", 2:"two"}},
		"empty input": {src: map[string]int{},                   expected: map[int]string{}},
		"panic if map has duplicates": {src: map[string]int{"one": 1, "two": 2, "duplicate":1}, panics: true},
	}

	for name, currTest := range testCases {
		t.Run(name, func(t *testing.T){ doTestRevert(t, &currTest) })
	}
}
