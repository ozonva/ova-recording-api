package utils

import "sort"

func contains(value int, array []int) bool {
	i := sort.SearchInts(array, value)
	return i < len(array) && array[i] == value
}

func Filter(input []int) []int {

	itemsToFilter := []int{1, 2, 3}

	out := make([]int, 0)
	for _, element := range input {
		if contains(element, itemsToFilter) {
			out = append(out, element)
		}
	}

	return out
}
