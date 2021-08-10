package utils

import "github.com/ozonva/ova-recording-api/pkg/recording"
import "errors"

func SplitAppointmentsToBatches(input []recording.Appointment, batchSize uint) ([][]recording.Appointment, error) {
	if batchSize == 0 {
		return nil, errors.New("batchSize must be greater than zero")
	}

	numBatches := uint(len(input)) / batchSize
	if numBatches * batchSize < uint(len(input)) {
		numBatches++
	}
	out := make([][]recording.Appointment, numBatches)
	var currSlice []recording.Appointment

	batchIdx := 0
	for _, element := range input {

		if uint(len(currSlice)) < batchSize {
			currSlice = append(currSlice, element)
		} else {
			out[batchIdx] = currSlice
			batchIdx++
			currSlice = []recording.Appointment{element}
		}
	}

	if len(currSlice) > 0 {
		out[batchIdx] = currSlice
	}

	return out, nil
}
