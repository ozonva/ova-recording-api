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
	"github.com/ozonva/ova-recording-api/pkg/recording"
	desc "github.com/ozonva/ova-recording-api/pkg/recording/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
  grpcPort = ":8888"
  grpcServerEndpoint = "localhost:8888"
)

type User struct {
	UserId int64 `db:"user_id"`
	Name string `db:"name"`
}

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

	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()
	fmt.Printf("Path to config: %s\n", configPath)

	cfg := api.ReadConfig(configPath)
	fmt.Printf("Parsed config: %v, conn string: %s\n", cfg, cfg.GetConnString())

	db := ConnectToDatabase(cfg)

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	currRepo := repo.NewRepo(db)

	currRepo.AddEntities([]recording.Appointment{{UserID: 2, Name: "some name", StartTime: time.Now(), EndTime: time.Now()}})

	//var user User
	//err := db.QueryRowx("select user_id, name from users limit 1").StructScan(&user)
	//if err != nil {
	//	log.Errorf("Cannot select: %s\n", err)
	//	return
	//}
	//
	//fmt.Printf("User: %v\n", user)


	//go runClient()
	//
	//if err := run(); err != nil {
	//	log.Fatal(err)
	//}
}

func run() error {
  listen, err := net.Listen("tcp", grpcPort)
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  s := grpc.NewServer(
  	grpc.UnaryInterceptor(
  		middleware.ChainUnaryServer(
  			api.RequestIdInterceptor,
  		),
  	),
  )
  desc.RegisterRecordingServiceServer(s, api.NewRecordingServiceAPI())

  log.Infof("Start serving on port %s...", grpcPort)

  if err := s.Serve(listen); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := desc.NewRecordingServiceClient(conn)

	for i := 0; i < 10; i++ {
		a := desc.InAppointmentV1{Name: "Some hello name"}
		_, err = client.CreateAppointmentV1(ctx, &desc.CreateAppointmentRequestV1{Appointment: &a})
		if err != nil {
			log.Errorf("Got error from server: %s", err)
		} else {
			log.Info("Successfully sent request to server")
			break
		}
		time.Sleep(time.Second*5)
	}
}
