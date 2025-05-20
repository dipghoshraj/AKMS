package main

import (
	"fmt"
	"tvs/consumer"
)

func main() {
	fmt.Println("Starting Kafka Consumer...")
	brokers := []string{"localhost:9092"}

	topic := "akm"
	consumer.ConsumeKafkaMessages(brokers, topic)
}
