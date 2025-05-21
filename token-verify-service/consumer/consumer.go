package consumer

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func (c *Consumer) ConsumeKafkaMessages() {
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
		Brokers:        c.brokers,
		Topic:          c.Topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		GroupID:        "my-group",
		StartOffset:    kafka.LastOffset, // for simplicity, start from the latest offset in scale we should use kafka.FirstOffset
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

		if err := c.processMessage(ctx, msg); err != nil {
			log.Printf("Error processing message: %v", err)
			// TODO: Use dead-letter queue or retry logic
			continue
		}

		// commit the message
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Failed to commit offset: %v", err)
		} else {
			log.Printf("Offset committed")
		}

	}
}
