package utils

import "fmt"

func Revert(input map[string]int) map[int]string{

	out := make(map[int]string)
	for key, value := range input {
		if v, ok := out[value]; ok {
			panic(fmt.Sprintf("Duplicate value in map. Key = %s, Value = %d. Duplicate has key: %s", key, value, v))
		}
		out[value] = key
	}

	return out
}
