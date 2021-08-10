package utils


func SplitToBatches(input []int, batchSize int) [][]int {
	if batchSize <= 0 {
		panic("batchSize must be greater than zero")
	}

	numBatches := len(input) / batchSize
	if numBatches * batchSize < len(input) {
		numBatches++
	}
	out := make([][]int, numBatches)
	var currSlice []int

	batchIdx := 0
	for _, element := range input {

		if len(currSlice) < batchSize {
			currSlice = append(currSlice, element)
		} else {
			out[batchIdx] = currSlice
			batchIdx++
			currSlice = []int{element}
		}
	}

	if len(currSlice) > 0 {
		out[batchIdx] = currSlice
	}

	return out
}
