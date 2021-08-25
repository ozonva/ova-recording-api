package main

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/internal/flusher"
	"github.com/ozonva/ova-recording-api/internal/repo"
	"github.com/ozonva/ova-recording-api/internal/saver"
	"github.com/ozonva/ova-recording-api/internal/utils"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	"sync"
	"time"
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

	//utils.OpenFileInCycle("/tmp/test.txt")

	dummyRepo := repo.NewDummyRepo()
	fl := flusher.NewFlusher(10, dummyRepo)
	s := saver.NewSaver(20, fl, time.Second*5)
	wg := sync.WaitGroup{}
	numGoroutines := 10
	numEntitiesPerGoroutine := 100
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func (myNum int) {
			for j := 0; j < numEntitiesPerGoroutine; j++ {
				err := s.Save(recording.Appointment{AppointmentID: uint64(j + numEntitiesPerGoroutine*myNum)})
				if err != nil {
					fmt.Printf("Cannot Save entitiy. what: %s\n", err)
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	//time.Sleep(time.Second * 10)

	s.Close()

	//time.Sleep(time.Second * 3)

	fmt.Printf("Added entities: %d\n", dummyRepo.GetAddedCount())
}
