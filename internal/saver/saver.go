package saver

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)


type Flusher interface {
	Flush(entities []recording.Appointment) []recording.Appointment
}

type Saver interface {
	Save(entity recording.Appointment) error
	Close()
}


type saver struct {
	entities []recording.Appointment
	fl Flusher
	m sync.Mutex
	saveInterval time.Duration
	doneCh chan bool
}

func NewSaver(capacity int, fl Flusher, saveInterval time.Duration) Saver {
	s := &saver{
		entities: make([]recording.Appointment, 0, capacity),
		fl: fl,
		saveInterval: saveInterval,
		doneCh: make(chan bool),
	}

	s.init()

	return s
}

func (s* saver) Save(entity recording.Appointment) error {
	s.m.Lock()
	defer s.m.Unlock()

	var err error
	if len(s.entities) == cap(s.entities) {
		err = s.doFlush()
	}

	s.entities = append(s.entities, entity)

	return err
}

func (s* saver) Close() {
	close(s.doneCh)

	err := s.flush()
	if err != nil {
		log.Errorf("Cannot Close saver: %s\n", err)
	}
}

func (s* saver) doFlush() error {
	log.Tracef("Going to flush %d entities\n", len(s.entities))

	notFlushed := s.fl.Flush(s.entities)

	s.entities = s.entities[:0]
	if notFlushed != nil {
		copy(s.entities, notFlushed)
		return fmt.Errorf("cannot flush %d entities", len(notFlushed))
	}
	return nil
}

func (s *saver) flush() error {
	s.m.Lock()
	defer s.m.Unlock()

	log.Tracef("Flushing...")

	return s.doFlush()
}

func (s *saver) init() {
	go func() {
		ticker := time.NewTicker(s.saveInterval)
		defer ticker.Stop()
		for {
			select {
				case <-s.doneCh:
					log.Traceln("Closing ticking goroutine")
					return
				case <-ticker.C:
					log.Traceln("Tick, saving")
					err := s.flush()
					if err != nil {
						log.Errorf("Cannot flush: %s\n", err)
					}
			}
		}

	}()
}
