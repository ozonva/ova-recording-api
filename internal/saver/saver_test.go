package saver_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/ozonva/ova-recording-api/internal/flusher"
	"github.com/ozonva/ova-recording-api/internal/repo"
	mock_repo "github.com/ozonva/ova-recording-api/internal/repo/mock"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/ozonva/ova-recording-api/internal/saver"
)

var _ = Describe("Saver", func() {
	var (
		someRepo *mock_repo.MockRepo
		ctrl *gomock.Controller
		someFlusher flusher.Flusher
		someSaver saver.Saver
		entities []recording.Appointment
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		someRepo = mock_repo.NewMockRepo(ctrl)
		someFlusher = flusher.NewFlusher(2, someRepo)
		someSaver = saver.NewSaver(2, someFlusher, time.Second*3)
		entities = []recording.Appointment{
			{
				UserID: 100,
				AppointmentID: 1,
				Name: "Some appointment1",
			},
			{
				UserID: 200,
				AppointmentID: 2,
				Name: "Some appointment2",
			},
			{
				UserID: 200,
				AppointmentID: 3,
				Name: "Some appointment3",
			},
			{
				UserID: 200,
				AppointmentID: 4,
				Name: "Some appointment4",
			},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Basic saving entries", func() {
		Context("Usual case", func() {
			It("should save all", func() {
				someRepo.EXPECT().AddEntities(entities[:2]).Return(nil).Times(1)
				someRepo.EXPECT().AddEntities(entities[2:]).Return(nil).Times(1)
				someRepo.EXPECT().GetAddedCount().Return(4).Times(1)
				for _, entity := range entities {
					err := someSaver.Save(entity)
					gomega.Expect(err).To(gomega.BeNil())
				}

				time.Sleep(time.Second * 5)

				someSaver.Close()

				someRepo.GetAddedCount()

			})
		})
	})
})

var _ = Describe("Saver Multi thread", func() {
	var (
		someRepo repo.Repo
		someFlusher flusher.Flusher
		someSaver saver.Saver
	)

	BeforeEach(func() {
		someRepo = repo.NewDummyRepo()
		someFlusher = flusher.NewFlusher(10, someRepo)
		someSaver = saver.NewSaver(20, someFlusher, time.Second*5)
		log.SetLevel(log.WarnLevel)
	})

	Describe("Basic saving entries", func() {
		Context("Multi threaded test", func() {
			It("should save all", func() {
				wg := sync.WaitGroup{}
				numGoroutines := 10
				numEntitiesPerGoroutine := 100
				wg.Add(numGoroutines)
				for i := 0; i < numGoroutines; i++ {
					go func (myNum int) {
						for j := 0; j < numEntitiesPerGoroutine; j++ {
							err := someSaver.Save(recording.Appointment{AppointmentID: uint64(j + numEntitiesPerGoroutine*myNum)})
							if err != nil {
								fmt.Printf("Cannot Save entitiy. what: %s\n", err)
							}
						}
						wg.Done()
					}(i)
				}

				wg.Wait()

				someSaver.Close()

				gomega.Expect(someRepo.GetAddedCount()).To(gomega.Equal(numEntitiesPerGoroutine*numGoroutines))
			})
		})
	})
})
