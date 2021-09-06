package repo

import (
	"context"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Repo interface {
	AddEntities(ctx context.Context, entities []recording.Appointment) ([]uint64, error)
	UpdateEntity(ctx context.Context,
				entityId uint64,
				userId uint64,
				name string,
				description string,
				startTime time.Time,
				endTime time.Time) error
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

func (r *repo) AddEntities(ctx context.Context, entities []recording.Appointment) (out []uint64, err error) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("appointments")
	ib.Cols("user_id", "name", "description", "start_time", "end_time")
	for _, ent := range entities {
		ib.Values(ent.UserID, ent.Name, ent.Description, ent.StartTime, ent.EndTime)
	}
	sql, args := ib.Build()

	sql = sql + "RETURNING appointment_id"

	res, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var currId uint64
		err = res.Scan(&currId)
		if err != nil {
			log.Warnf("cannot get returned id: %s", err)
			return nil, nil
		}
		out = append(out, currId)
	}

	err = res.Close()
	if err != nil {
		log.Warnf("cannot close result set: %s", err)
		return nil, nil
	}

	r.m.Lock()
	r.addedCount += len(out)
	r.m.Unlock()

	return out, nil
}

func (r *repo) UpdateEntity(ctx context.Context,
				entityId uint64,
				userId uint64,
				name string,
				description string,
				startTime time.Time,
				endTime time.Time) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update("appointments")

	if userId > 0 {
		ub.Set(ub.Assign("user_id", userId))
	}
	if len(name) > 0 {
		ub.SetMore(ub.Assign("name", name))
	}
	if len(description) > 0 {
		ub.SetMore(ub.Assign("description", description))
	}
	if !startTime.IsZero() {
		ub.SetMore(ub.Assign("start_time", startTime))
	}
	if !endTime.IsZero() {
		ub.SetMore(ub.Assign("end_time", endTime))
	}

	ub.Where(ub.Equal("appointment_id", entityId))

	sql, args := ub.Build()

	_, err := r.db.ExecContext(ctx, sql, args...)
	return err
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

	for result.Next() {
		var a recording.Appointment
		err := result.StructScan(&a)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}

	if result.Err() != nil {
		return out, result.Err()
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

func (r *dummyRepo) AddEntities(ctx context.Context, entities []recording.Appointment) ([]uint64, error) {
	r.m.Lock()
	defer r.m.Unlock()
	r.addedCount += len(entities)
	for _, entity := range entities {
		log.Infof("dummyRepo: Add entity %s\n", entity)

	}
	return nil, nil
}

func (r *dummyRepo) UpdateEntity(ctx context.Context,
				entityId uint64,
				userId uint64,
				name string,
				description string,
				startTime time.Time,
				endTime time.Time) error {
	log.Infof("dummyRepo: UpdateEntity(%d, %d, %s, %s, %v, %v",
		entityId, userId, name, description, startTime, endTime)

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
