package main

import (
	"context"
	"flag"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ozonva/ova-recording-api/internal/app/metrics"
	api "github.com/ozonva/ova-recording-api/internal/app/recording"
	"github.com/ozonva/ova-recording-api/internal/kafka_client"
	"github.com/ozonva/ova-recording-api/internal/repo"
	desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"net"
	"net/http"
	"time"
)

const (
  grpcPort = ":8888"
  grpcServerEndpoint = "localhost:8888"
  prometheusMetricsPort = ":8081"
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
	go startPrometheusMetricsEndpoint()

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
			log.Errorf("cannot close db connection: %s", err)
		}
	}(db)

	currRepo := repo.NewRepo(db)
	kfkClient := kafka_client.NewKafkaClient("ova_recording_server")
	err = kfkClient.Connect(context.Background(), "localhost:9092", "appointments", 0)
	if err != nil {
		log.Fatalf("Cannot connect to Kafka: %s", err)
	}
	defer func(kfkClient kafka_client.Client) {
		err := kfkClient.Close()
		if err != nil {
			log.Errorf("cannot close Kafka client: %s", err)
		}
	}(kfkClient)

	prometheusMetrics := metrics.NewApiMetrics()

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
		kfkClient,
		prometheusMetrics,
	))

	grpc_prometheus.Register(server)

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

func startPrometheusMetricsEndpoint() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(prometheusMetricsPort, nil); err != nil {
		log.Errorf("Failed to start listen to metric requests, error %s", err)
	}
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
	for i := 0; i < 100; i++ {
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

		_, err = client.UpdateAppointmentV1(ctx, &desc.UpdateAppointmentV1Request{Appointment: &desc.OutAppointmentV1{AppointmentId: uint64(i),Name: "updated name"}})
		if err != nil {
			log.Infof("cannot update entity %d", i)
		}
		time.Sleep(time.Second*5)
	}

	server.GracefulStop()
}
