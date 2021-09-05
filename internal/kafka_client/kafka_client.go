package kafka_client

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Client interface {
	Name() string
	Connect(ctx context.Context, address string, topic string, partition int) error
	SendMessage(msg []byte) error
	Close() error
}

func NewKafkaClient(name string) Client {
	return &kafkaClient{name: name}
}

type kafkaClient struct {
	conn *kafka.Conn
	name string
}

func (receiver *kafkaClient) Name() string {
	return receiver.name
}

func (receiver *kafkaClient) Connect(ctx context.Context, address string, topic string, partition int) (err error) {
	receiver.conn, err = kafka.DialLeader(ctx, "tcp", address, topic, partition)
	return err
}

func (receiver *kafkaClient) SendMessage(msg []byte) (err error) {
	if receiver.conn == nil {
		logrus.Fatalf("Cannot send message to nil connection")
	}
	_, err = receiver.conn.WriteMessages(kafka.Message{Value: msg})
	return err
}

func (receiver *kafkaClient) Close() error {
	if receiver.conn == nil {
		return nil
	}
	return receiver.conn.Close()
}
