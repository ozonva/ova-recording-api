package recording

import (
  "context"
  "fmt"
  "github.com/opentracing/opentracing-go"
  opLog "github.com/opentracing/opentracing-go/log"
  "github.com/ozonva/ova-recording-api/internal/kafka_client"
  "github.com/ozonva/ova-recording-api/pkg/recording"
  desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
  kafkaproto "github.com/ozonva/ova-recording-api/pkg/recording/kafka"
  log "github.com/sirupsen/logrus"
  "github.com/uber/jaeger-client-go"
  jaegercfg "github.com/uber/jaeger-client-go/config"
  jaegerlog "github.com/uber/jaeger-client-go/log"
  "github.com/uber/jaeger-lib/metrics"
  "google.golang.org/protobuf/proto"
  "google.golang.org/protobuf/types/known/emptypb"
  "google.golang.org/protobuf/types/known/timestamppb"
  "io"
  "time"
)

type Repo interface {
	AddEntities(ctx context.Context, entities []recording.Appointment) ([]uint64, error)
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

type ServiceAPI struct {
  desc.UnimplementedRecordingServiceServer
  repository Repo
  batchSize int
  kfkClient kafka_client.Client
}

func NewRecordingServiceAPI(inRepo Repo, batchSize int, client kafka_client.Client) desc.RecordingServiceServer {
  api := &ServiceAPI{
    repository: inRepo,
    batchSize: batchSize,
    kfkClient: client,
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

func SetupTracing() io.Closer {
	cfg := jaegercfg.Configuration{
        ServiceName: "grpc_recording_api",
        Sampler:     &jaegercfg.SamplerConfig{
            Type:  jaeger.SamplerTypeConst,
            Param: 1,
        },
        Reporter:    &jaegercfg.ReporterConfig{
            LogSpans: true,
        },
    }

    // Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
    // and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
    // frameworks.
    jLogger := jaegerlog.StdLogger
    jMetricsFactory := metrics.NullFactory

    // Initialize tracer with a logger and a metrics factory
    tracer, tracingCloser, err := cfg.NewTracer(
        jaegercfg.Logger(jLogger),
        jaegercfg.Metrics(jMetricsFactory),
    )

	if err != nil {
		log.Errorf("cannot create tracer: %s", err)
	}

    // Set the singleton opentracing.Tracer with the Jaeger tracer.
    opentracing.SetGlobalTracer(tracer)

    return tracingCloser
}

func (service *ServiceAPI) sendMessageToKafka(m proto.Message) error {
  msg, err := proto.Marshal(m)
  if err != nil {
    return err
  }

  err = service.kfkClient.SendMessage(msg)
  return err
}

func (service *ServiceAPI) sendCreatedEvent(entityId uint64) error {
  event := kafkaproto.KafkaMessage{
    Kind: kafkaproto.KafkaMessage_CREATED,
    Producer: service.kfkClient.Name(),
    Body: &kafkaproto.KafkaMessage_Created{Created: &kafkaproto.AppointmentCreatedV1{AppointmentId: entityId}}}

  return service.sendMessageToKafka(&event)
}

func (service *ServiceAPI) sendUpdatedEvent(entityId uint64) error {
  event := kafkaproto.KafkaMessage{
    Kind: kafkaproto.KafkaMessage_UPDATED,
    Producer: service.kfkClient.Name(),
    Body: &kafkaproto.KafkaMessage_Updated{Updated: &kafkaproto.AppointmentUpdatedV1{AppointmentId: entityId}}}

  return service.sendMessageToKafka(&event)
}

func (service *ServiceAPI) sendDeletedEvent(entityId uint64) error {
  event := kafkaproto.KafkaMessage{
    Kind: kafkaproto.KafkaMessage_DELETED,
    Producer: service.kfkClient.Name(),
    Body: &kafkaproto.KafkaMessage_Deleted{Deleted: &kafkaproto.AppointmentDeletedV1{AppointmentId: entityId}}}

  return service.sendMessageToKafka(&event)
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

  res, err := service.repository.AddEntities(ctx, []recording.Appointment{app})
  if err != nil {
    GetLogger(ctx).Errorf("Cannot add entity: %s", err)
  }

  GetLogger(ctx).Infof("Added entity with id %d", res[0])

  err = service.sendCreatedEvent(res[0])
  if err != nil {
    GetLogger(ctx).Warnf("Cannot send CUD event: %s", err)
    return out, err
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
    return out, err
  }

  err = service.sendUpdatedEvent(req.AppointmentId)
  if err != nil {
    GetLogger(ctx).Warnf("Cannot send CUD event: %s", err)
  }

  return &emptypb.Empty{}, err
}

func (service *ServiceAPI) MultiCreateAppointmentsV1(ctx context.Context, req *desc.MultiCreateAppointmentsV1Request) (out *emptypb.Empty, err error) {
  GetLogger(ctx).Infof("Got MultiCreateAppointmentsV1 request: %s", req)
  out = &emptypb.Empty{}

  parentSpan := ctx.Value(SpanKey).(opentracing.Span)

  if req.Appointments == nil {
    err = fmt.Errorf("request field `Appointments` is nil")
    GetLogger(ctx).Error(err)
    return
  }

  currSlice := make([]recording.Appointment, 0, service.batchSize)
  for _, inApp := range req.Appointments {
    if len(currSlice) < service.batchSize {
      currSlice = append(currSlice, AppointmentFromApiInput(inApp))
    } else {
      err = service.insertBatch(ctx, service.repository, parentSpan, currSlice)
      if err != nil {
        return out, err
      }
      currSlice = make([]recording.Appointment, 0, service.batchSize)
    }
  }

  if len(currSlice) > 0 {
    err = service.insertBatch(ctx, service.repository, parentSpan, currSlice)
  }

  return out, err
}

func (service *ServiceAPI) insertBatch(ctx context.Context, repository Repo, parentSpan opentracing.Span, entities []recording.Appointment) error {
  tracer := opentracing.GlobalTracer()
  childSpan := tracer.StartSpan("Batch", opentracing.ChildOf(parentSpan.Context()))
  childSpan.LogFields(opLog.Int("batch size", len(entities)))
  defer childSpan.Finish()
  res, err := repository.AddEntities(ctx, entities)
  if err != nil {
    GetLogger(ctx).Errorf("Cannot add entities: %s", err)
    return err
  }

  for _, entityId := range res {
    err = service.sendCreatedEvent(entityId)
    if err != nil {
      GetLogger(ctx).Warnf("Cannot send CUD event: %s", err)
    }
  }

  return err
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
    return &emptypb.Empty{}, err
  }

  err = service.sendDeletedEvent(req.AppointmentId)
  if err != nil {
    GetLogger(ctx).Warnf("Cannot send CUD event: %s", err)
  }

  return &emptypb.Empty{}, err
}
