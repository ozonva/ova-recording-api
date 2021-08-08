package utils

func Batch(input []int, batchSize int) [][]int {
	var out [][]int
	var currSlice []int

	if batchSize <= 0 {
		panic("batchSize must be greater than zero")
	}

	for _, element := range input {

		if len(currSlice) < batchSize {
			currSlice = append(currSlice, element)
		} else {
			out = append(out, currSlice)
			currSlice = []int{element}
		}
	}

	if len(currSlice) > 0 {
		out = append(out, currSlice)
	}

	return out
}
