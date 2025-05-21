package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"tvs/consumer"
	"tvs/dbops"
	rest "tvs/http"
	"tvs/store"

	"github.com/gorilla/mux"
)

func setupRouter(service *rest.Server) *mux.Router {
	router := mux.NewRouter()
	rest.SetupRoutes(router, service)

	return router
}

func main() {
	store.InitDB()
	store.InitRedis()
	Migrate()

	// Rest API setup
	dataops := dbops.NewOpsManager(store.DataBase, store.RedisClient)
	services := rest.NewServerOps(dataops)

	router := setupRouter(services)

	server := &http.Server{
		Addr:         ":" + store.GetEnv("SERVER_PORT", "8000"),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Printf("Starting server on port %s\n", store.GetEnv("SERVER_PORT", "8080"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()

	go func() {
		// kafka consumer setup
		fmt.Println("Starting Kafka Consumer...")
		brokers := []string{"localhost:9092"}

		topic := "akm"
		consumer := consumer.NewComsumer(brokers, topic, dataops)
		consumer.ConsumeKafkaMessages()
	}()
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown: %v\n", err)
	}

	log.Println("Server exited properly")

}
