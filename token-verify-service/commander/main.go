package main

import (
	"fmt"
	"tvs/consumer"
	"tvs/store"
)

func main() {
	store.InitDB()
	Migrate()

	fmt.Println("Starting Kafka Consumer...")
	brokers := []string{"localhost:9092"}

	topic := "akm"
	consumer.ConsumeKafkaMessages(brokers, topic)
}
