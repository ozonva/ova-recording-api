package saver

import (
	"fmt"
	"github.com/ozonva/ova-recording-api/internal/flusher"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	"sync"
	"time"
)


type Saver interface {
	Save(entity recording.Appointment)
	Close()
}


type saver struct {
	entities []recording.Appointment
	fl flusher.Flusher
	m sync.Mutex
	saveInterval int
	doneCh chan bool
}

func NewSaver(capacity int, fl flusher.Flusher, saveInterval int) Saver {
	s := &saver{
		entities: make([]recording.Appointment, 0, capacity),
		fl: fl,
		saveInterval: saveInterval,
		doneCh: make(chan bool),
	}

	s.init()

	return s
}

func (s* saver) Save(entity recording.Appointment) {
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.entities) == cap(s.entities) {
		s.doFlush()
	}

	s.entities = append(s.entities, entity)
}

func (s* saver) Close() {
	s.doneCh <- true
	close(s.doneCh)

	s.flush()
}

func (s* saver) doFlush() {
	fmt.Println("Going to flush", len(s.entities), "entities")

	s.fl.Flush(s.entities)

	s.entities = s.entities[:0]
}

func (s *saver) flush() {
	s.m.Lock()
	defer s.m.Unlock()

	s.doFlush()
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
					s.flush()
			}
		}

	}()
}
