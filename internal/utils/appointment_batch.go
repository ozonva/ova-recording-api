package utils

import "github.com/ozonva/ova-recording-api/pkg/recording"
import "errors"

func SplitAppointmentsToBatches(input []recording.Appointment, batchSize int) ([][]recording.Appointment, error) {
	if batchSize <= 0 {
		return nil, errors.New("batchSize must be greater than zero")
	}

	numBatches := (len(input) + batchSize - 1) / batchSize

	out := make([][]recording.Appointment, numBatches)

	batchStartIdx := 0
	batchIdx := 0
	for batchStartIdx < len(input) {
		numLeft := min(len(input)-batchStartIdx, batchSize)
		currBatch := make([]recording.Appointment, numLeft)
		copy(currBatch, input[batchStartIdx:batchStartIdx+numLeft])
		batchStartIdx += numLeft
		out[batchIdx] = currBatch
		batchIdx++
	}

	return out, nil
}
