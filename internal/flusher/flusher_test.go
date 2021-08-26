package flusher_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/ozonva/ova-recording-api/internal/flusher"
	mock_repo "github.com/ozonva/ova-recording-api/internal/repo/mock"
	"github.com/ozonva/ova-recording-api/pkg/recording"
)

var _ = Describe("Flusher", func() {
	var (
		someRepo *mock_repo.MockRepo
		ctrl *gomock.Controller
		someFlusher flusher.Flusher
		entities []recording.Appointment
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		someRepo = mock_repo.NewMockRepo(ctrl)
		someFlusher = flusher.NewFlusher(2, someRepo)
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

	Describe("Flushing entries", func() {
		Context("ok scenario", func() {
			It("should flush all", func() {
				someRepo.EXPECT().AddEntities(entities[:2]).Return(nil).Times(1)
				someRepo.EXPECT().AddEntities(entities[2:]).Return(nil).Times(1)

				unhandled := someFlusher.Flush(entities)
				gomega.Expect(unhandled).To(gomega.BeNil())

				someRepo.EXPECT().GetAddedCount().Return(4).Times(1)
				someRepo.GetAddedCount()
			})
		})
		Context("fail scenario", func() {
			It("should return unhandled entities", func() {
				someRepo.EXPECT().AddEntities(entities[:2]).Return(nil).Times(1)
				someRepo.EXPECT().AddEntities(entities[2:]).Return(errors.New("repoError")).Times(1)
				unhandled := someFlusher.Flush(entities)
				gomega.Expect(unhandled).To(gomega.Equal(entities[2:]))
			})
		})
	})
})

var _ = Describe("Flusher errors", func() {
	var (
		someRepo    *mock_repo.MockRepo
		ctrl        *gomock.Controller
		someFlusher flusher.Flusher
		entities    []recording.Appointment
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		someRepo = mock_repo.NewMockRepo(ctrl)
		someFlusher = flusher.NewFlusher(-1, someRepo)
		entities = []recording.Appointment{
			{
				UserID: 100,
				AppointmentID: 1,
				Name: "Some appointment1",
			},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("negative chunk size", func() {
			It("should return all entries", func() {
				unhandled := someFlusher.Flush(entities)
				gomega.Expect(unhandled).To(gomega.Equal(entities))
			})
	})
	Context("nil entries", func() {
			It("should return nil", func() {
				unhandled := someFlusher.Flush(nil)
				gomega.Expect(unhandled).To(gomega.BeNil())
			})
	})
})
