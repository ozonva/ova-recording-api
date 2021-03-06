package utils

import (
	"reflect"
	"testing"
)

type testCaseRevert struct {
	name string
	src map[string]int
	expected map[int]string
	expectingError bool
}

func doTestRevert(t *testing.T, currTestCase* testCaseRevert) {

	reverse, err := Revert(currTestCase.src)

	if err != nil {
		if  !currTestCase.expectingError {
			t.Fatalf("Got unexpected error %s", err)
		}
		return
	} else if currTestCase.expectingError {
		t.Fatalf("Expected error, but got nil")
	}

	t.Logf("Expected: %v, actual: %v", currTestCase.expected, reverse)

	if len(currTestCase.expected) != len(reverse) {
		t.Fatalf("Result size %d != %d", len(reverse), len(currTestCase.expected))
	}

	if !reflect.DeepEqual(currTestCase.expected, reverse) {
		t.Fatal("Result of the reversion is not as expected")
	}
}


func TestRevert(t *testing.T) {

	testCases := []testCaseRevert{
		{name: "simple",      src: map[string]int{"one": 1, "two": 2}, expected: map[int]string{1:"one", 2:"two"}},
		{name: "empty input", src: map[string]int{},                   expected: map[int]string{}},
		{name: "error if map has duplicates", src: map[string]int{"one": 1, "two": 2, "duplicate":1}, expectingError: true},
	}

	for _, currTest := range testCases {
		t.Run(currTest.name, func(t *testing.T){ doTestRevert(t, &currTest) })
	}
}
