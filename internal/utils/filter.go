package utils


func contains(value int, array []int) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}
	return false
}

func FilterBy(input []int, targetItems []int) []int {

	out := make([]int, 0)
	for _, element := range input {
		if contains(element, targetItems) {
			out = append(out, element)
		}
	}

	return out
}
