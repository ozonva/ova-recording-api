package flusher

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/internal/repo"
	"github.com/ozonva/ova-recording-api/internal/utils"
	"github.com/ozonva/ova-recording-api/pkg/recording"
)

type Flusher interface {
	Flush(entities []recording.Appointment) ([]recording.Appointment)
}

func NewFlusher(
	chunkSize int,
	entityRepo repo.Repo,
) Flusher {
	return &flusher{
		chunkSize: chunkSize,
		entityRepo:  entityRepo,
	}
}

type flusher struct {
	chunkSize int
	entityRepo  repo.Repo
}

func (f *flusher) Flush (entities []recording.Appointment) []recording.Appointment {

	if entities == nil {
		fmt.Println("Nil input")
		return entities
	}

	batches, err := utils.SplitAppointmentsToBatches(entities, f.chunkSize)
	if err != nil {
		fmt.Printf("Cannot split entities to batches: %s\n", err)
		return entities
	}

	currIndex := 0

	for _, batch := range batches {
		err = f.entityRepo.AddEntities(batch)
		if err != nil {
			fmt.Printf("Cannot save to repo: %s\n", err)
			return entities[currIndex:]
		}
		currIndex += len(batch)
	}

	return nil
}
