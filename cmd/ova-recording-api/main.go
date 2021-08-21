package main

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/internal/flusher"
	"github.com/ozonva/ova-recording-api/internal/repo"
	"github.com/ozonva/ova-recording-api/internal/saver"
	"github.com/ozonva/ova-recording-api/internal/utils"
	"github.com/ozonva/ova-recording-api/pkg/recording"
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
	fl := flusher.NewFlusher(3, dummyRepo)

	s := saver.NewSaver(5, fl, 5)
	s.Save(recording.Appointment{AppointmentID: 1})
	s.Save(recording.Appointment{AppointmentID: 2})
	time.Sleep(time.Second * 3)
	s.Save(recording.Appointment{AppointmentID: 3})
	s.Save(recording.Appointment{AppointmentID: 4})
	time.Sleep(time.Second * 3)
	s.Save(recording.Appointment{AppointmentID: 5})
	s.Save(recording.Appointment{AppointmentID: 6})
	s.Close()
}
