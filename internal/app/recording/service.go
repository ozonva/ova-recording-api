package recording

import (
  "context"
  "fmt"
  "github.com/ozonva/ova-recording-api/internal/repo"
  "github.com/ozonva/ova-recording-api/pkg/recording"
  desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
  "google.golang.org/protobuf/types/known/emptypb"
  "google.golang.org/protobuf/types/known/timestamppb"
  "time"
)

type Repo interface {
	AddEntities(ctx context.Context, entities []recording.Appointment) error
    UpdateEntity(ctx context.Context,
				entityId uint64,
				userId uint64,
				name string,
				description string,
				startTime time.Time,
				endTime time.Time) error
    ListEntities(ctx context.Context, limit, offset uint64) ([]recording.Appointment, error)
    DescribeEntity(ctx context.Context, entityId uint64) (*recording.Appointment, error)
	RemoveEntity(ctx context.Context, entityId uint64) error
	GetAddedCount(ctx context.Context) int
}

type Flusher interface {
	Flush(entities []recording.Appointment) []recording.Appointment
}

type Saver interface {
	Save(entity recording.Appointment) error
	Close()
}

type ServiceAPI struct {
  desc.UnimplementedRecordingServiceServer
  repository Repo
  currSaver Saver
}

func NewRecordingServiceAPI(inRepo repo.Repo, sv Saver) desc.RecordingServiceServer {
  api := &ServiceAPI{
    repository: inRepo,
    currSaver: sv,
  }
  return api
}

func AppointmentToApiInput (appointment *recording.Appointment) *desc.InAppointmentV1 {
  return &desc.InAppointmentV1{
    UserId: appointment.UserID,
    Name: appointment.Name,
    Description: appointment.Description,
    StartTime: timestamppb.New(appointment.StartTime.UTC()),
    EndTime: timestamppb.New(appointment.EndTime.UTC()),
  }
}

func AppointmentFromApiInput (appointment *desc.InAppointmentV1) recording.Appointment {
  return recording.Appointment{
    UserID: appointment.UserId,
    Name: appointment.Name,
    Description: appointment.Description,
    StartTime: appointment.StartTime.AsTime(),
    EndTime: appointment.EndTime.AsTime(),
  }
}

func AppointmentToApiOutput (appointment *recording.Appointment) *desc.OutAppointmentV1 {
  return &desc.OutAppointmentV1{
    AppointmentId: appointment.AppointmentID,
    UserId: appointment.UserID,
    Name: appointment.Name,
    Description: appointment.Description,
    StartTime: timestamppb.New(appointment.StartTime),
    EndTime: timestamppb.New(appointment.EndTime),
  }
}

func (service *ServiceAPI) CreateAppointmentV1(ctx context.Context, req *desc.CreateAppointmentV1Request) (out *emptypb.Empty, err error) {
  GetLogger(ctx).Infof("Got CreateAppointmentV1 request: %s", req)

  if req.Appointment == nil {
    err = fmt.Errorf("request field `Appointment` is nil")
    GetLogger(ctx).Error(err)
    return
  }

  app := AppointmentFromApiInput(req.Appointment)

  GetLogger(ctx).Infof("Try to add %v", []recording.Appointment{app})

  err = service.repository.AddEntities(ctx, []recording.Appointment{app})
  if err != nil {
    GetLogger(ctx).Errorf("Cannot add entity: %s", err)
  }

  return &emptypb.Empty{}, err
}

func (service *ServiceAPI) UpdateAppointmentV1(ctx context.Context, req *desc.UpdateAppointmentV1Request) (out *emptypb.Empty, err error) {
  GetLogger(ctx).Infof("Got UpdateAppointmentV1 request: %s", req)

  err = service.repository.UpdateEntity(ctx,
    req.AppointmentId, req.UserId,
    req.Name, req.Description,
    req.StartTime.AsTime(), req.EndTime.AsTime(),
  )

  if err != nil {
    GetLogger(ctx).Errorf("Cannot update entity: %s", err)
  }

  return &emptypb.Empty{}, err
}

func (service *ServiceAPI) MultiCreateAppointmentsV1(ctx context.Context, req *desc.MultiCreateAppointmentsV1Request) (out *emptypb.Empty, err error) {
  GetLogger(ctx).Infof("Got MultiCreateAppointmentsV1 request: %s", req)

  if req.Appointments == nil {
    err = fmt.Errorf("request field `Appointments` is nil")
    GetLogger(ctx).Error(err)
    return
  }

  for _, inApp := range req.Appointments {
    app := AppointmentFromApiInput(inApp)
    GetLogger(ctx).Infof("Try to add %v", app)
    err := service.currSaver.Save(app)
    if err != nil {
      GetLogger(ctx).Errorf("Cannot Save entity: %s", err)
      break
    }
  }

  return &emptypb.Empty{}, err
}

func (service *ServiceAPI) DescribeAppointmentV1(ctx context.Context, req *desc.DescribeAppointmentV1Request) (*desc.DescribeAppointmentV1Response, error) {
  GetLogger(ctx).Infof("Got DescribeAppointmentV1 request: %s", req)

  app, err := service.repository.DescribeEntity(ctx, req.AppointmentId)
  if err != nil {
    GetLogger(ctx).Errorf("cannot describe appointment: %s", err)
  }

  out := AppointmentToApiOutput(app)

  return &desc.DescribeAppointmentV1Response{Appointment: out}, nil
}

func (service *ServiceAPI) ListAppointmentsV1(ctx context.Context, req *desc.ListAppointmentsV1Request) (*desc.ListAppointmentsV1Response, error) {
  GetLogger(ctx).Infof("Got ListAppointmentsV1 request: %s", req)

  res, err := service.repository.ListEntities(ctx, req.Limit, req.Offset)
  if err != nil {
    GetLogger(ctx).Errorf("Cannot list: %s", err)
    return nil ,err
  }

  out := &desc.ListAppointmentsV1Response{Appointments: make([]*desc.OutAppointmentV1, len(res))}
  for i, app := range res {
    out.Appointments[i] = AppointmentToApiOutput(&app)
  }

  return out, nil
}

func (service *ServiceAPI) RemoveAppointmentV1(ctx context.Context, req *desc.RemoveAppointmentV1Request) (*emptypb.Empty, error) {
  GetLogger(ctx).Infof("Got RemoveAppointmentV1 request: %s", req)

  err := service.repository.RemoveEntity(ctx, req.AppointmentId)
  if err != nil {
    GetLogger(ctx).Errorf("Cannot remove entity %d: %s", req.AppointmentId, err)
  }

  return &emptypb.Empty{}, err
}
