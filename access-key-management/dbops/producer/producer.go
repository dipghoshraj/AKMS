package producer

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	producer *kafka.Writer
}

func NewProducer() *kafkaProducer {
	return &kafkaProducer{
		producer: creation(),
	}
}

/*
TODO : need to use gorouting and sync
also need to user defer for writer close
*/

func creation() *kafka.Writer {

	kafkaURL := "localhost:9092"

	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    "service",
		Balancer: &kafka.LeastBytes{},
	}
	return kafkaWriter
}

func (kf *kafkaProducer) PushMessage(request_id string, message map[string]string) error {

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(request_id),
		Value: []byte(jsonData),
	}

	err = kf.producer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}
	return nil
}
