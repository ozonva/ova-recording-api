package main

import (
	"context"
	api "github.com/ozonva/ova-recording-api/internal/app/recording"
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


func run() error {
  listen, err := net.Listen("tcp", grpcPort)
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  s := grpc.NewServer()
  desc.RegisterRecordingServiceServer(s, api.NewRecordingServiceAPI())

  log.Infof("Start serving on port %s...", grpcPort)

  if err := s.Serve(listen); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }

  return nil
}

func main() {
	go func() {
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
			a := desc.Appointment{AppointmentId: 1, Name: "Some hello name"}
			_, err = client.CreateAppointmentV1(ctx, &desc.CreateAppointmentRequestV1{Appointment: &a})
			if err != nil {
				log.Errorf("Got error from server: %s", err)
			} else {
				log.Info("Successfully sent request to server")
				break
			}
			time.Sleep(time.Second*5)
		}
	}()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
