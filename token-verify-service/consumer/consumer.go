package consumer

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func ConsumeKafkaMessages(brokers []string, topic string) {
	// Create a new Kafka reader
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown listener
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Shutting down Kafka consumer...")
		cancel()
	}()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
	})

	defer reader.Close()

	for {
		// Read messages from the topic
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			log.Printf("Kafka fetch error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var reqID string
		for _, h := range msg.Headers {
			if strings.ToLower(h.Key) == "request_id" {
				reqID = string(h.Value)
			}
		}
		log.Printf("[request_id=%s] Fetched message: %s", reqID, string(msg.Value))

		if err := processMessage(ctx, reqID, msg); err != nil {
			log.Printf("[request_id=%s] Error processing message: %v", reqID, err)
			// TODO: Use dead-letter queue or retry logic
			continue
		}

		// commit the message
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("[request_id=%s] Failed to commit offset: %v", reqID, err)
		} else {
			log.Printf("[request_id=%s] Offset committed", reqID)
		}

	}
}
