package utils

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	"reflect"
	"testing"
	"time"
)

type testCaseBatchAppointment struct {
	name string
	src []recording.Appointment
	batchSize int
	expected [][]recording.Appointment
	expectingError bool
}

func doTestBatchAppointment(t *testing.T, currTestCase *testCaseBatchAppointment) {

	batches, err := SplitAppointmentsToBatches(currTestCase.src, currTestCase.batchSize)

	if err != nil {
		if  !currTestCase.expectingError {
			t.Fatalf("Got unexpected error %s", err)
		}
		return
	} else if currTestCase.expectingError {
		t.Fatalf("Expected error, but got nil")
	}

	if len(currTestCase.expected) != len(batches) {
		t.Fatalf("Result size %d != %d", len(batches), len(currTestCase.expected))
	}

	if !reflect.DeepEqual(currTestCase.expected, batches) {
		t.Fatal("Result batches are not as expected")
	}
}


func TestSplitAppointmentsToBatchesToBatches(t *testing.T) {

	appointments := make([]recording.Appointment, 6)
	for i := range appointments {
		now := time.Now()
		appointments[i] = recording.Appointment{
			UserID: 1,
			AppointmentID: uint64(i + 1),
			Name: fmt.Sprintf("Appointment №%d", i),
			Description: fmt.Sprintf("Description for appointment %d", i),
			StartTime: now,
			EndTime: now.Add(time.Duration(3600 + i)),
		}
	}

	testCases := []testCaseBatchAppointment{
		{name: "exact division", src: appointments, batchSize: 3,  expected: [][]recording.Appointment{appointments[0:3], appointments[3:]}},
		{name: "with remainder", src: appointments[:4],     batchSize: 3,  expected: [][]recording.Appointment{appointments[0:3], appointments[3:4]}},
		{name: "single batch", src: appointments,     batchSize: 10, expected: [][]recording.Appointment{appointments}},
		{name: "nil input", src: nil,                batchSize: 10, expected: [][]recording.Appointment{}},
		{name: "error if zero batch size", src: appointments, batchSize: 0, expectingError: true},
	}

	for _, currTest := range testCases {
		t.Run(currTest.name, func(t *testing.T){ doTestBatchAppointment(t, &currTest) })
	}
}
