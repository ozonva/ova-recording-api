package utils

var itemsToFilter = []int{1, 2, 3}


func contains(value int, array []int) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}
	return false
}

func FilterBy(input []int) []int {

	out := make([]int, 0)
	for _, element := range input {
		if contains(element, itemsToFilter) {
			out = append(out, element)
		}
	}

	return out
}
