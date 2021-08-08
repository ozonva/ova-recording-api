package main

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/internal/utils"
)


func main() {
	src := []int{1,2,3,4,5}
	batchSize:= 3
	fmt.Println("src:", src, "batch size:", batchSize,"batches:", utils.Batch(src, batchSize))

	srcMap := map[string]int{"one": 1, "two": 2, "three": 3}

	fmt.Println(utils.Revert(&srcMap))
}
