package utils

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/pkg/recording"
)

func AppointmentsSliceToMap(appointments []recording.Appointment) (map[uint64]recording.Appointment, error) {
	out := make(map[uint64]recording.Appointment)

	for _, app := range appointments {
		if _, ok := out[app.AppointmentID]; ok {
			return nil, fmt.Errorf("duplicate id %d", app.AppointmentID)
		}
		out[app.AppointmentID] = app
	}

	return out, nil
}
