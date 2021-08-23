package utils

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	"reflect"
	"testing"
)

type testCaseSplitAppointments struct {
	name string
	src []recording.Appointment
	expected map[uint64]recording.Appointment
	expectingError bool
}

func doTestSplitAppointments(t *testing.T, currTestCase *testCaseSplitAppointments) {

	result, err := AppointmentsSliceToMap(currTestCase.src)

	if err != nil {
		if  !currTestCase.expectingError {
			t.Fatalf("Got unexpected error %s", err)
		}
		return
	} else if currTestCase.expectingError {
		t.Fatalf("Expected error, but got nil")
	}

	if len(currTestCase.expected) != len(result) {
		t.Fatalf("Result size %d != %d", len(result), len(currTestCase.expected))
	}

	if !reflect.DeepEqual(currTestCase.expected, result) {
		t.Fatal("Result batches are not as expected")
	}
}


func TestAppointmentsSliceToMap(t *testing.T) {

	appointments := make([]recording.Appointment, 2)
	for i := range appointments {
		appointments[i] = recording.Appointment{
			UserID: 1,
			AppointmentID: uint64(i + 1),
			Name: fmt.Sprintf("Appointment №%d", i),
		}
	}

	testCases := []testCaseSplitAppointments{
		{name: "basic", src: appointments, expected: map[uint64]recording.Appointment{1: appointments[0], 2: appointments[1]}},
		{name: "nil input", src: nil,                expected: map[uint64]recording.Appointment{}},
		{name: "error if has duplicates", src: []recording.Appointment{appointments[0], appointments[0]}, expectingError: true},
	}

	for _, currTest := range testCases {
		t.Run(currTest.name, func(t *testing.T){ doTestSplitAppointments(t, &currTest) })
	}
}
