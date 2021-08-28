package repo

import (
	"github.com/jmoiron/sqlx"
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

func NewRepo(db *sqlx.DB) Repo {
	return &repo{db: db}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) AddEntities(entities []recording.Appointment) error {
	query := "INSERT INTO appointments(user_id, name, description, start_time, end_time) VALUES (:user_id, :name, :description, :start_time, :end_time)"
	for _, ent := range entities {
		_, err := r.db.NamedExec(query, &ent)
		if err != nil {
			log.Errorf("cannot insert entity: `%v`, error: %s", ent, err)
			return err
		}
	}

	return nil
}

func (r *repo) ListEntities(limit, offset uint64) ([]recording.Appointment, error) {
	return nil, nil
}

func (r *repo) DescribeEntity(entityId uint64) (*recording.Appointment, error) {
	return nil, nil
}

func (r *repo) GetAddedCount() int {
	return 0
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
