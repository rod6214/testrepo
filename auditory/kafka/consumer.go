package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go/snappy"
)

type Consumer struct {
	reader *kafka.Reader
}

func (consumer *Consumer) Connect(brokers []string, groupID string, topic string) {
	config := kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	}
	log.Println("Connecting to Kafka as new reader...")
	consumer.reader = kafka.NewReader(config)
}

func (consumer *Consumer) Consume() error {
	log.Println("Waiting for new messages...")
	message, err := consumer.reader.ReadMessage(context.Background())
	if err != nil {
		log.Println("ERR: ", err)
		return err
	}
	log.Printf("New message: Topic: %s - Partition: %d - Offset: %d - Key: %s - Message: %s",
		message.Topic, message.Partition, message.Offset, message.Key, message.Value)

	return nil
}

func (consumer *Consumer) Close() {
	consumer.reader.Close()
}
