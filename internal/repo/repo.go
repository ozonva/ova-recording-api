package repo

import (
	"context"
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

func NewDummyRepo() Repo {
	return &dummyRepo{}
}

type repo struct {
	db *sqlx.DB
	addedCount int
	m sync.Mutex
}

func (r *repo) AddEntities(entities []recording.Appointment) error {
	query := "INSERT INTO appointments(user_id, name, description, start_time, end_time) VALUES (:user_id, :name, :description, :start_time, :end_time)"
	tx, err := r.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, ent := range entities {
		_, err := tx.NamedExec(query, ent)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	r.m.Lock()
	r.addedCount += len(entities)
	r.m.Unlock()

	return nil
}

func (r *repo) ListEntities(limit, offset uint64) ([]recording.Appointment, error) {
	query := "SELECT appointment_id, user_id, name, description, start_time, end_time FROM appointments LIMIT $1 OFFSET $2"

	result, err := r.db.Queryx(query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer func(result *sqlx.Rows) {
		err := result.Close()
		if err != nil {
		}
	}(result)

	out := make([]recording.Appointment, 0)

	for ok := result.Next(); ok; ok = result.Next() {
		var a recording.Appointment
		err := result.StructScan(&a)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (r *repo) DescribeEntity(entityId uint64) (*recording.Appointment, error) {
	query := "SELECT appointment_id, user_id, name, description, start_time, end_time FROM appointments WHERE appointment_id = $1"

	result := r.db.QueryRowx(query, entityId)
	var a recording.Appointment
	err := result.StructScan(&a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *repo) GetAddedCount() int {
	r.m.Lock()
	out := r.addedCount
	r.m.Unlock()
	return out
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
