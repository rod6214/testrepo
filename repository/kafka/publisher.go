package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

type Publisher struct {
	writer *kafka.Writer
}

func (publisher *Publisher) Connect(clientId string, brokers []string, topic string) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}
	config := kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	publisher.writer = kafka.NewWriter(config)
}

func (publisher *Publisher) Publish(ctx context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	log.Println("Writing message...")
	writeMessage := publisher.writer.WriteMessages(ctx, message)
	if writeMessage != nil {
		log.Println(writeMessage)
	}
	log.Printf("Message: %s", value)
	return writeMessage
}
