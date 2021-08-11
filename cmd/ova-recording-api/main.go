package main

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/internal/utils"
)


func main() {
	src := []int{1,2,3,4,5}
	batchSize:= 3
	batches, err := utils.SplitToBatches(src, batchSize)
	if err == nil {
		fmt.Println("src:", src, "batch size:", batchSize,"batches:", batches)
	}

	srcMap := map[string]int{"one": 1, "two": 2, "three": 3}

	fmt.Println(utils.Revert(srcMap))

	fmt.Println(utils.FilterBy([]int{1,2,3,4,5,6}, []int{1,6,4}))

	utils.OpenFileInCycle("/tmp/test.txt")
}
