package main

import (
	"context"
	"flag"
	"fmt"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	api "github.com/ozonva/ova-recording-api/internal/app/recording"
	"github.com/ozonva/ova-recording-api/internal/repo"
	desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net"
	"time"
)

const (
  grpcPort = ":8888"
  grpcServerEndpoint = "localhost:8888"
)

var (
	server *grpc.Server
)

func ConnectToDatabase(cfg api.Config) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.GetConnString())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	return db
}

func main() {
	go runClient()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	cfg := api.ReadConfig(configPath)

	db := ConnectToDatabase(cfg)

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	currRepo := repo.NewRepo(db)

	server = grpc.NewServer(
		grpc.UnaryInterceptor(
			middleware.ChainUnaryServer(
				api.RequestIdInterceptor,
				api.TracingInterceptor,
			),
		),
	)

	desc.RegisterRecordingServiceServer(server, api.NewRecordingServiceAPI(
		currRepo,
		cfg.ChunkSize,
	))

	tracingCloser := api.SetupTracing()
	defer func(tracingCloser io.Closer) {
		err := tracingCloser.Close()
		if err != nil {
			log.Errorf("Cannot close opentracing tracer: %s", err)
		}
	}(tracingCloser)

	log.Infof("Start serving on port %s...", grpcPort)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Info("Shutting down...")

	return nil
}

func runClient() {
	conn, err := grpc.Dial(grpcServerEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Warnf("Cannot close client connection: %s", err)
		}
	}(conn)

	log.Info("Dialed")

	ctx := context.Background()

	client := desc.NewRecordingServiceClient(conn)

	log.Info("Try to make requests")

	a := desc.InAppointmentV1{Name: "Some hello name"}
	_, err = client.CreateAppointmentV1(ctx, &desc.CreateAppointmentV1Request{Appointment: &a})
	if err != nil {
		log.Errorf("cannot create appointment: %s", err)
		return
	}
	resp, err := client.ListAppointmentsV1(ctx, &desc.ListAppointmentsV1Request{Offset: 0, Limit: 100})
	if err != nil {
		log.Errorf("cannot list appointments: %s", err)
		return
	}

	log.Infof("List: %v", resp)

	req := &desc.MultiCreateAppointmentsV1Request{}
	for i := 0; i < 100; i++ {
		req.Appointments = append(req.Appointments, &desc.InAppointmentV1{
			UserId: uint64(i),
			Name: fmt.Sprintf("Name %d", i),
			StartTime: timestamppb.New(time.Now()),
			EndTime: timestamppb.New(time.Now()),
		})
	}

	_, err = client.MultiCreateAppointmentsV1(ctx, req)

	if err != nil {
		log.Errorf("cannot multi add appointments: %s", err)
		return
	}

	server.GracefulStop()
}
