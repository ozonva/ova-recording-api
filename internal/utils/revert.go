package utils

import (
	"errors"
	"fmt"
)

func Revert(input map[string]int) (map[int]string, error){

	out := make(map[int]string)
	for key, value := range input {
		if v, ok := out[value]; ok {
			format := "Duplicate value in map. Key = %s, Value = %d. Duplicate has key: %s"
			return nil, errors.New(fmt.Sprintf(format, key, value, v))
		}
		out[value] = key
	}

	return out, nil
}
