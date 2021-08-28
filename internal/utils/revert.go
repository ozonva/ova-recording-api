package utils

import (
	"fmt"
)

func Revert(input map[string]int) (map[int]string, error){

	out := make(map[int]string)
	for key, value := range input {
		if v, ok := out[value]; ok {
			format := "duplicate value in map. Key = %s, Value = %d. Duplicate has key: %s"
			return nil, fmt.Errorf(format, key, value, v)
		}
		out[value] = key
	}

	return out, nil
}
