package repo

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/pkg/recording"
)

type Repo interface {
	AddEntities(entities []recording.Appointment) error
    ListEntities(limit, offset uint64) ([]recording.Appointment, error)
    DescribeEntity(entityId uint64) (*recording.Appointment, error)
}

func NewDummyRepo() Repo {
	return &dummyRepo{}
}

type dummyRepo struct {

}

func (r *dummyRepo) AddEntities(entities []recording.Appointment) error {
	fmt.Println("Hello from dummyRepo::AddEntities")
	return nil
}

func (r *dummyRepo) ListEntities(limit, offset uint64) ([]recording.Appointment, error) {
	fmt.Println("Hello from dummyRepo::ListEntities")
	return []recording.Appointment{}, nil
}

func (r *dummyRepo) DescribeEntity(entityId uint64) (*recording.Appointment, error) {
	fmt.Println("Hello from dummyRepo::DescribeEntity")
	return nil, nil
}
