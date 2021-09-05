package recording_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/ozonva/ova-recording-api/internal/app/recording"
	mock_repo "github.com/ozonva/ova-recording-api/internal/repo/mock"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var _ = Describe("Service", func() {
	var (
		ctrl *gomock.Controller
		someRepo *mock_repo.MockRepo
		srv desc.RecordingServiceServer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		someRepo = mock_repo.NewMockRepo(ctrl)
		srv = api.NewRecordingServiceAPI(someRepo)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("ok scenario", func() {
		It("create", func() {
			entity := recording.Appointment{
				UserID: 1,
				StartTime: time.Now().UTC(),
				EndTime: time.Now().UTC(),
			}
			someRepo.EXPECT().AddEntities(gomock.Any(), []recording.Appointment{entity}).Return(nil).Times(1)

			_, err := srv.CreateAppointmentV1(context.Background(), &desc.CreateAppointmentV1Request{
				Appointment: &desc.InAppointmentV1{
					UserId: entity.UserID,
					StartTime: timestamppb.New(entity.StartTime),
					EndTime: timestamppb.New(entity.EndTime),
				},
			})

			Expect(err).To(BeNil())
		})

		It("describe", func() {
			entity := recording.Appointment{
				AppointmentID: 1,
				UserID: 1,
				StartTime: time.Now().UTC(),
				EndTime: time.Now().UTC(),
			}
			someRepo.EXPECT().DescribeEntity(gomock.Any(), entity.AppointmentID).Return(&entity, nil).Times(1)

			resp, err := srv.DescribeAppointmentV1(context.Background(), &desc.DescribeAppointmentV1Request{
				AppointmentId: entity.AppointmentID,
			})

			Expect(err).To(BeNil())

			Expect(resp.Appointment.AppointmentId).To(Equal(entity.AppointmentID))
			Expect(resp.Appointment.UserId).To(Equal(entity.UserID))
			Expect(resp.Appointment.Name).To(Equal(entity.Name))
			Expect(resp.Appointment.Description).To(Equal(entity.Description))
			Expect(resp.Appointment.StartTime).To(Equal(timestamppb.New(entity.StartTime)))
			Expect(resp.Appointment.EndTime).To(Equal(timestamppb.New(entity.EndTime)))
		})

		It("list", func() {
			entity := recording.Appointment{
				AppointmentID: 1,
				UserID: 1,
				Name: "some name",
				Description: "some desc",
				StartTime: time.Now().UTC(),
				EndTime: time.Now().UTC(),
			}
			someRepo.EXPECT().ListEntities(gomock.Any(), uint64(1), uint64(0)).Return([]recording.Appointment{entity}, nil).Times(1)

			resp, err := srv.ListAppointmentsV1(context.Background(), &desc.ListAppointmentsV1Request{
				Offset: 0, Limit: 1,
			})

			Expect(err).To(BeNil())

			Expect(resp.Appointments[0].AppointmentId).To(Equal(entity.AppointmentID))
			Expect(resp.Appointments[0].UserId).To(Equal(entity.UserID))
			Expect(resp.Appointments[0].Name).To(Equal(entity.Name))
			Expect(resp.Appointments[0].Description).To(Equal(entity.Description))
			Expect(resp.Appointments[0].StartTime).To(Equal(timestamppb.New(entity.StartTime)))
			Expect(resp.Appointments[0].EndTime).To(Equal(timestamppb.New(entity.EndTime)))
		})

		It("remote", func() {

			someRepo.EXPECT().RemoveEntity(gomock.Any(), uint64(1)).Return(nil).Times(1)

			_, err := srv.RemoveAppointmentV1(context.Background(), &desc.RemoveAppointmentV1Request{
				AppointmentId: 1,
			})

			Expect(err).To(BeNil())
		})
	})
})
