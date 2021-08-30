package recording_test

import (
	"context"
	"github.com/golang/mock/gomock"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	api "github.com/ozonva/ova-recording-api/internal/app/recording"
	mock_repo "github.com/ozonva/ova-recording-api/internal/repo/mock"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
  grpcPort = ":8888"
  grpcServerEndpoint = "localhost:8888"
)

var _ = Describe("Service", func() {
	var (
		ctrl *gomock.Controller
		ctx context.Context
		someRepo *mock_repo.MockRepo
		listen net.Listener
		srv *grpc.Server
		conn *grpc.ClientConn
		client desc.RecordingServiceClient
		entities []recording.Appointment
		err error
	)
	BeforeEach(func() {
		entities = []recording.Appointment{
			{
				AppointmentID: 1,
				Name: "Name1",
				UserID: 1,
				Description: "some desc 1",
				StartTime: time.Now(),
				EndTime: time.Now(),
			},
			{
				AppointmentID: 2,
				Name: "Name2",
				UserID: 1,
				Description: "some desc 2",
				StartTime: time.Now(),
				EndTime: time.Now(),
			},
		}
		ctrl = gomock.NewController(GinkgoT())
		ctx = context.Background()
		someRepo = mock_repo.NewMockRepo(ctrl)
		srv = grpc.NewServer(
			grpc.UnaryInterceptor(
				middleware.ChainUnaryServer(
					api.RequestIdInterceptor,
				),
			),
		)
		desc.RegisterRecordingServiceServer(srv, api.NewRecordingServiceAPI(someRepo))

		listen, err = net.Listen("tcp", grpcPort)
		if err != nil {
			logrus.Fatalf("failed to listen: %v", err)
		}

		conn, err = grpc.Dial(grpcServerEndpoint, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logrus.Fatalf("did not connect: %v", err)
		}

		client = desc.NewRecordingServiceClient(conn)

		time.Sleep(time.Second*1)
	})

	AfterEach(func() {
		ctrl.Finish()
		err := conn.Close()
		if err != nil {
			logrus.Warnf("Cannot close client connection: %s", err)
		}
	})

	Describe("Serving", func() {
		Context("ok scenario", func() {
			It("Should handle all", func() {
				go func() {
					if err := srv.Serve(listen); err != nil {
						logrus.Fatalf("failed to serve: %v", err)
					}
				}()
				someRepo.EXPECT().AddEntities(gomock.Any(), entities[0]).Return(nil).Times(1)
				//someRepo.EXPECT().AddEntities(ctx, entities[1]).Return(nil).Times(1)

				_, err := client.CreateAppointmentV1(ctx,
					&desc.CreateAppointmentV1Request{
						Appointment: api.AppointmentToApiInput(&entities[0]),
					},
				)
				gomega.Expect(err).To(gomega.BeNil())
				//
				//_, err = client.CreateAppointmentV1(ctx,
				//	&desc.CreateAppointmentV1Request{
				//		Appointment: api.AppointmentToApiInput(&entities[1]),
				//	},
				//)
				//gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
