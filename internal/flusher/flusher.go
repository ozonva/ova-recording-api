package flusher

import (
	"github.com/ozonva/ova-recording-api/internal/repo"
	"github.com/ozonva/ova-recording-api/internal/utils"
	"github.com/ozonva/ova-recording-api/pkg/recording"
)

type Flusher interface {
	Flush(entities []recording.Appointment) ([]recording.Appointment, error)
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

func (f *flusher) Flush (entities []recording.Appointment) ([]recording.Appointment, error) {

	batches, err := utils.SplitAppointmentsToBatches(entities, f.chunkSize)
	if err != nil {
		return entities, err
	}

	currIndex := 0

	for _, batch := range batches {
		err = f.entityRepo.AddEntities(batch)
		if err != nil {
			return entities[currIndex:], err
		}
		currIndex += len(batch)
	}

	return nil, nil
}
