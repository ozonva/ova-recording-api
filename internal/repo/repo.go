package repo

import (
	"github.com/ozonva/ova-recording-api/pkg/recording"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Repo interface {
	AddEntities(entities []recording.Appointment) error
    ListEntities(limit, offset uint64) ([]recording.Appointment, error)
    DescribeEntity(entityId uint64) (*recording.Appointment, error)
	GetAddedCount() int
}

func NewDummyRepo() Repo {
	return &dummyRepo{}
}

type dummyRepo struct {
	addedCount int
	m sync.Mutex
}

func (r *dummyRepo) AddEntities(entities []recording.Appointment) error {
	r.m.Lock()
	defer r.m.Unlock()
	r.addedCount += len(entities)
	for _, entity := range entities {
		log.Infof("dummyRepo: Add entity %s\n", entity)

	}
	return nil
}

func (r *dummyRepo) ListEntities(limit, offset uint64) ([]recording.Appointment, error) {
	log.Infof("dummyRepo: list entities, offset: %d, limit: %d\n", offset, limit)
	return []recording.Appointment{}, nil
}

func (r *dummyRepo) DescribeEntity(entityId uint64) (*recording.Appointment, error) {
	log.Infof("dummyRepo: describe entity with id: %d\n", entityId)
	return nil, nil
}

func (r *dummyRepo) GetAddedCount() int {
	return r.addedCount
}
