package recording

import (
  "context"
  desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
  log "github.com/sirupsen/logrus"
  "google.golang.org/protobuf/types/known/emptypb"
)

type ServiceAPI struct {
  desc.UnimplementedRecordingServiceServer
}

func NewRecordingServiceAPI() desc.RecordingServiceServer {
  return &ServiceAPI{
  }
}

func (a *ServiceAPI) CreateAppointmentV1(ctx context.Context, req *desc.CreateAppointmentRequestV1) (*emptypb.Empty, error) {
  log.Infof("Got CreateAppointmentV1 request: %s", req)
  return &emptypb.Empty{}, nil
}

func (a *ServiceAPI) DescribeAppointmentV1(ctx context.Context, req *desc.DescribeAppointmentRequestV1) (*desc.Appointment, error) {
  log.Infof("Got DescribeAppointmentV1 request: %s", req)
  return &desc.Appointment{Name: "not implemented"}, nil
}

func (a *ServiceAPI) ListAppointmentsV1(ctx context.Context, req *desc.ListAppointmentsRequestV1) (*desc.ListAppointmentsResponseV1, error) {
  log.Infof("Got ListAppointmentsV1 request: %s", req)
  return &desc.ListAppointmentsResponseV1{}, nil
}

func (a *ServiceAPI) RemoveAppointmentV1(ctx context.Context, req *desc.RemoveAppointmentRequestV1) (*emptypb.Empty, error) {
  log.Infof("Got RemoveAppointmentV1 request: %s", req)
  return &emptypb.Empty{}, nil
}
