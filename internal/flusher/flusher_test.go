package flusher_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/ozonva/ova-recording-api/internal/flusher"
	repo_mock "github.com/ozonva/ova-recording-api/internal/repo/mock"
	"github.com/ozonva/ova-recording-api/pkg/recording"
)

var _ = Describe("Flusher", func() {
	var (
		someRepo *repo_mock.MockRepo
		ctrl *gomock.Controller
		someFlusher flusher.Flusher
		entities []recording.Appointment
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		someRepo = repo_mock.NewMockRepo(ctrl)
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
		Context("Hmm", func() {
			It("should flush all", func() {
				someRepo.EXPECT().AddEntities(entities[:2]).Return(nil).Times(1)
				someRepo.EXPECT().AddEntities(entities[2:]).Return(nil).Times(1)
				unhandled := someFlusher.Flush(entities)
				gomega.Expect(unhandled).To(gomega.BeNil())
			})
		})
	})
})
