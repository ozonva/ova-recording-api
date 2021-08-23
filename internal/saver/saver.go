package saver

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	"sync"
	"time"
)


type Flusher interface {
	Flush(entities []recording.Appointment) ([]recording.Appointment, error)
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
		fmt.Printf("Cannot Close saver: %s\n", err)
	}
}

func (s* saver) doFlush() error {
	fmt.Println("Going to flush", len(s.entities), "entities")

	notFlushed, err := s.fl.Flush(s.entities)

	s.entities = s.entities[:0]
	if notFlushed != nil {
		copy(s.entities, notFlushed)
	}

	return err
}

func (s *saver) flush() error {
	s.m.Lock()
	defer s.m.Unlock()

	fmt.Println("Flushing...")

	return s.doFlush()
}

func (s *saver) init() {
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(s.saveInterval))
		defer ticker.Stop()
		for {
			select {
				case <-s.doneCh:
					fmt.Println("Closing ticking goroutine")
					return
				case <-ticker.C:
					fmt.Println("Tick, saving")
					err := s.flush()
					if err != nil {
						fmt.Printf("Cannot flush: %s\n", err)
					}
			}
		}

	}()
}
