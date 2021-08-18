package flusher_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/ozonva/ova-recording-api/internal/flusher"
	"github.com/ozonva/ova-recording-api/internal/repo"
	"github.com/ozonva/ova-recording-api/pkg/recording"
)

var _ = Describe("Flusher", func() {
	var (
		someRepo repo.Repo
		ctrl *gomock.Controller
		someFlusher flusher.Flusher
		entities []recording.Appointment
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		someRepo = repo.NewDummyRepo()
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
				_, err := someFlusher.Flush(entities)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
