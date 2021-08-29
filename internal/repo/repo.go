package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Repo interface {
	AddEntities(ctx context.Context, entities []recording.Appointment) error
    ListEntities(ctx context.Context, limit, offset uint64) ([]recording.Appointment, error)
    DescribeEntity(ctx context.Context, entityId uint64) (*recording.Appointment, error)
	RemoveEntity(ctx context.Context, entityId uint64) error
	GetAddedCount(ctx context.Context) int
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

func (r *repo) AddEntities(ctx context.Context, entities []recording.Appointment) error {
	query := "INSERT INTO appointments(user_id, name, description, start_time, end_time) VALUES (:user_id, :name, :description, :start_time, :end_time)"
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	for _, ent := range entities {
		_, err := tx.NamedExecContext(ctx, query, ent)
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

func (r *repo) ListEntities(ctx context.Context, limit, offset uint64) ([]recording.Appointment, error) {
	query := "SELECT appointment_id, user_id, name, description, start_time, end_time FROM appointments LIMIT $1 OFFSET $2"

	result, err := r.db.QueryxContext(ctx, query, limit, offset)
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

func (r *repo) DescribeEntity(ctx context.Context, entityId uint64) (*recording.Appointment, error) {
	query := "SELECT appointment_id, user_id, name, description, start_time, end_time FROM appointments WHERE appointment_id = $1"

	result := r.db.QueryRowxContext(ctx, query, entityId)
	var a recording.Appointment
	err := result.StructScan(&a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *repo) RemoveEntity(ctx context.Context, entityId uint64) error {
	query := "DELETE FROM appointments WHERE appointment_id = $1"

	_, err := r.db.ExecContext(ctx, query, entityId)

	return err
}

func (r *repo) GetAddedCount(ctx context.Context) int {
	r.m.Lock()
	out := r.addedCount
	r.m.Unlock()
	return out
}

type dummyRepo struct {
	addedCount int
	m sync.Mutex
}

func (r *dummyRepo) AddEntities(ctx context.Context, entities []recording.Appointment) error {
	r.m.Lock()
	defer r.m.Unlock()
	r.addedCount += len(entities)
	for _, entity := range entities {
		log.Infof("dummyRepo: Add entity %s\n", entity)

	}
	return nil
}

func (r *dummyRepo) ListEntities(ctx context.Context, limit, offset uint64) ([]recording.Appointment, error) {
	log.Infof("dummyRepo: list entities, offset: %d, limit: %d\n", offset, limit)
	return []recording.Appointment{}, nil
}

func (r *dummyRepo) DescribeEntity(ctx context.Context, entityId uint64) (*recording.Appointment, error) {
	log.Infof("dummyRepo: describe entity with id: %d\n", entityId)
	return nil, nil
}

func (r *dummyRepo) RemoveEntity(ctx context.Context, entityId uint64) error {
	log.Infof("dummyRepo: remove entity with id: %d\n", entityId)
	return nil
}

func (r *dummyRepo) GetAddedCount(ctx context.Context) int {
	return r.addedCount
}
