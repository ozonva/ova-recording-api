package recording_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opentracing/opentracing-go"
	api "github.com/ozonva/ova-recording-api/internal/app/recording"
	mock_kafka_client "github.com/ozonva/ova-recording-api/internal/kafka_client/mock"
	mock_repo "github.com/ozonva/ova-recording-api/internal/repo/mock"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"time"
)

var _ = Describe("Service", func() {
	var (
		ctrl *gomock.Controller
		someRepo *mock_repo.MockRepo
		srv desc.RecordingServiceServer
		ctx context.Context
		tracingCloser io.Closer
		span opentracing.Span
		kfkClient *mock_kafka_client.MockClient
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		someRepo = mock_repo.NewMockRepo(ctrl)
		kfkClient = mock_kafka_client.NewMockClient(ctrl)
		srv = api.NewRecordingServiceAPI(someRepo, 10, kfkClient)
		ctx = api.AddValue(context.Background(), "test", 1)
		tracingCloser = api.SetupTracing()
		tracer := opentracing.GlobalTracer()
		span = tracer.StartSpan("TEST")
		ctx = context.WithValue(ctx, api.SpanKey, span)
	})

	AfterEach(func() {
		span.Finish()
		err := tracingCloser.Close()
		if err != nil {
			log.Errorf("Cannot close tracing %s", err)
		}
		ctrl.Finish()
	})

	Context("ok scenario", func() {
		It("create", func() {
			entity := recording.Appointment{
				UserID: 1,
				StartTime: time.Now().UTC(),
				EndTime: time.Now().UTC(),
			}
			someRepo.EXPECT().AddEntities(gomock.Any(), []recording.Appointment{entity}).Return([]uint64{1}, nil).Times(1)
			kfkClient.EXPECT().Name().Return("test").Times(1)
			kfkClient.EXPECT().SendMessage(gomock.Any()).Return(nil).Times(1)
			_, err := srv.CreateAppointmentV1(ctx, &desc.CreateAppointmentV1Request{
				Appointment: &desc.InAppointmentV1{
					UserId: entity.UserID,
					StartTime: timestamppb.New(entity.StartTime),
					EndTime: timestamppb.New(entity.EndTime),
				},
			})

			Expect(err).To(BeNil())
		})

		It("update", func() {
			entity := recording.Appointment{
				AppointmentID: 1,
				UserID: 1,
				Name: "name",
				Description: "desc",
				StartTime: time.Now().UTC(),
				EndTime: time.Now().UTC(),
			}
			someRepo.EXPECT().UpdateEntity(gomock.Any(),
				entity.AppointmentID, entity.UserID,
				entity.Name, entity.Description,
				entity.StartTime, entity.EndTime).Return(nil).Times(1)
			kfkClient.EXPECT().Name().Return("test").Times(1)
			kfkClient.EXPECT().SendMessage(gomock.Any()).Return(nil).Times(1)

			_, err := srv.UpdateAppointmentV1(ctx, &desc.UpdateAppointmentV1Request{
				AppointmentId: entity.AppointmentID,
				UserId: entity.UserID,
				Name: entity.Name,
				Description: entity.Description,
				StartTime: timestamppb.New(entity.StartTime),
				EndTime: timestamppb.New(entity.EndTime),
			})

			Expect(err).To(BeNil())
		})

		It("multi create", func() {
			entities := []recording.Appointment{
				{
					UserID:    1,
					StartTime: time.Now().UTC(),
					EndTime:   time.Now().UTC(),
				},
				{
					UserID:    2,
					StartTime: time.Now().UTC(),
					EndTime:   time.Now().UTC(),
				},
			}
			someRepo.EXPECT().AddEntities(gomock.Any(), entities).Return([]uint64{1, 2}, nil).Times(1)
			kfkClient.EXPECT().Name().Return("test").Times(2)
			kfkClient.EXPECT().SendMessage(gomock.Any()).Return(nil).Times(2)

			_, err := srv.MultiCreateAppointmentsV1(ctx, &desc.MultiCreateAppointmentsV1Request{
				Appointments: []*desc.InAppointmentV1{
					api.AppointmentToApiInput(&entities[0]),
					api.AppointmentToApiInput(&entities[1]),
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

			resp, err := srv.DescribeAppointmentV1(ctx, &desc.DescribeAppointmentV1Request{
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

			resp, err := srv.ListAppointmentsV1(ctx, &desc.ListAppointmentsV1Request{
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

		It("remove", func() {

			someRepo.EXPECT().RemoveEntity(gomock.Any(), uint64(1)).Return(nil).Times(1)
			kfkClient.EXPECT().Name().Return("test").Times(1)
			kfkClient.EXPECT().SendMessage(gomock.Any()).Return(nil).Times(1)

			_, err := srv.RemoveAppointmentV1(ctx, &desc.RemoveAppointmentV1Request{
				AppointmentId: 1,
			})

			Expect(err).To(BeNil())
		})
	})
})
