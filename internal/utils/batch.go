package utils

import (
	"errors"
)

// This sucks
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func SplitToBatches(input []int, batchSize int) ([][]int, error) {
	if batchSize <= 0 {
		return nil, errors.New("batchSize must be greater than zero")
	}

	numBatches := (len(input) + batchSize - 1) / batchSize

	out := make([][]int, numBatches)

	batchStartIdx := 0
	batchIdx := 0
	for batchStartIdx < len(input) {
		numLeft := min(len(input)-batchStartIdx, batchSize)
		currBatch := make([]int, numLeft)
		copy(currBatch, input[batchStartIdx:batchStartIdx+numLeft])
		batchStartIdx += numLeft
		out[batchIdx] = currBatch
		batchIdx++
	}

	return out, nil
}
