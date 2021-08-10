package utils

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	"reflect"
	"testing"
	"time"
)

type testCaseBatchAppointment struct {
	src []recording.Appointment
	batchSize uint
	expected [][]recording.Appointment
	hasError bool
}

func doTestBatchAppointment(t *testing.T, currTestCase *testCaseBatchAppointment) {

	batches, err := SplitAppointmentsToBatches(currTestCase.src, currTestCase.batchSize)
	//t.Logf("Expected: %v, actual: %v", currTestCase.expected, batches)

	if currTestCase.hasError && err == nil {
		t.Fatalf("Expected error, but got nil")
	}
	if !currTestCase.hasError && err != nil {
		t.Fatalf("Got unexpected error: %s", err)
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

	testCases := map[string]testCaseBatchAppointment{
		"exact division": {src: appointments, batchSize: 3,  expected: [][]recording.Appointment{appointments[0:3], appointments[3:]}},
		"with remainder": {src: appointments[:4],     batchSize: 3,  expected: [][]recording.Appointment{appointments[0:3], appointments[3:4]}},
		"single batch":   {src: appointments,     batchSize: 10, expected: [][]recording.Appointment{appointments}},
		"nil input":      {src: nil,                batchSize: 10, expected: [][]recording.Appointment{}},
		"error if zero batch size": {src: appointments, batchSize: 0, hasError: true},
	}

	for name, currTest := range testCases {
		t.Run(name, func(t *testing.T){ doTestBatchAppointment(t, &currTest) })
	}
}
