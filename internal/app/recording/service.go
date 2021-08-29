package recording

import (
  "context"
  desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
  "google.golang.org/protobuf/types/known/emptypb"
)

type ServiceAPI struct {
  desc.UnimplementedRecordingServiceServer
}

func NewRecordingServiceAPI() desc.RecordingServiceServer {
  return &ServiceAPI{
  }
}

func (a *ServiceAPI) CreateAppointmentV1(ctx context.Context, req *desc.CreateAppointmentV1Request) (*emptypb.Empty, error) {
  GetLogger(ctx).Infof("Got CreateAppointmentV1 request: %s", req)
  return &emptypb.Empty{}, nil
}

func (a *ServiceAPI) DescribeAppointmentV1(ctx context.Context, req *desc.DescribeAppointmentV1Request) (*desc.DescribeAppointmentV1Response, error) {
  GetLogger(ctx).Infof("Got DescribeAppointmentV1 request: %s", req)
  return &desc.DescribeAppointmentV1Response{Appointment: &desc.Appointment{Name: "not implemented"}}, nil
}

func (a *ServiceAPI) ListAppointmentsV1(ctx context.Context, req *desc.ListAppointmentsV1Request) (*desc.ListAppointmentsV1Response, error) {
  GetLogger(ctx).Infof("Got ListAppointmentsV1 request: %s", req)
  return &desc.ListAppointmentsV1Response{}, nil
}

func (a *ServiceAPI) RemoveAppointmentV1(ctx context.Context, req *desc.RemoveAppointmentV1Request) (*emptypb.Empty, error) {
  GetLogger(ctx).Infof("Got RemoveAppointmentV1 request: %s", req)
  return &emptypb.Empty{}, nil
}
