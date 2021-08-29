package recording

import (
  "context"
  "github.com/ozonva/ova-recording-api/internal/repo"
  "github.com/ozonva/ova-recording-api/pkg/recording"
  desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
  "google.golang.org/protobuf/types/known/emptypb"
  "google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceAPI struct {
  desc.UnimplementedRecordingServiceServer
  r repo.Repo
}

func NewRecordingServiceAPI(inRepo repo.Repo) desc.RecordingServiceServer {
  return &ServiceAPI{
    r: inRepo,
  }
}

func (a *ServiceAPI) CreateAppointmentV1(ctx context.Context, req *desc.CreateAppointmentV1Request) (*emptypb.Empty, error) {
  GetLogger(ctx).Infof("Got CreateAppointmentV1 request: %s", req)

  app := recording.Appointment{
    UserID: req.Appointment.UserId,
    Name: req.Appointment.Name,
    Description: req.Appointment.Description,
    StartTime: req.Appointment.StartTime.AsTime(),
    EndTime: req.Appointment.EndTime.AsTime(),
  }

  err := a.r.AddEntities(ctx, []recording.Appointment{app})
  if err != nil {
    GetLogger(ctx).Errorf("Cannot add entity: %s", err)
  }

  return &emptypb.Empty{}, err
}

func (a *ServiceAPI) DescribeAppointmentV1(ctx context.Context, req *desc.DescribeAppointmentV1Request) (*desc.DescribeAppointmentV1Response, error) {
  GetLogger(ctx).Infof("Got DescribeAppointmentV1 request: %s", req)

  app, err := a.r.DescribeEntity(ctx, req.AppointmentId)
  if err != nil {
    GetLogger(ctx).Errorf("cannot describe appointment: %s", err)
  }

  out := desc.OutAppointmentV1{
    AppointmentId: app.AppointmentID,
    UserId: app.UserID,
    Name: app.Name,
    Description: app.Description,
    StartTime: timestamppb.New(app.StartTime),
    EndTime: timestamppb.New(app.EndTime),
  }

  return &desc.DescribeAppointmentV1Response{Appointment: &out}, nil
}

func (a *ServiceAPI) ListAppointmentsV1(ctx context.Context, req *desc.ListAppointmentsV1Request) (*desc.ListAppointmentsV1Response, error) {
  GetLogger(ctx).Infof("Got ListAppointmentsV1 request: %s", req)

  res, err := a.r.ListEntities(ctx, req.Num, req.FromId)
  if err != nil {
    GetLogger(ctx).Errorf("Cannot list: %s", err)
    return nil ,err
  }

  out := &desc.ListAppointmentsV1Response{Appointments: make([]*desc.OutAppointmentV1, len(res))}
  for i, app := range res {
    out.Appointments[i] = &desc.OutAppointmentV1{
      AppointmentId: app.AppointmentID,
      UserId: app.UserID,
      Name: app.Name,
      Description: app.Description,
      StartTime: timestamppb.New(app.StartTime),
      EndTime: timestamppb.New(app.EndTime),
    }
  }

  return out, nil
}

func (a *ServiceAPI) RemoveAppointmentV1(ctx context.Context, req *desc.RemoveAppointmentV1Request) (*emptypb.Empty, error) {
  GetLogger(ctx).Infof("Got RemoveAppointmentV1 request: %s", req)

  err := a.r.RemoveEntity(ctx, req.AppointmentId)
  if err != nil {
    GetLogger(ctx).Errorf("Cannot remove entity %d: %s", req.AppointmentId, err)
  }

  return &emptypb.Empty{}, err
}
