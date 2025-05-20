package main

import (
	"akm/dbops"
	"akm/store"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"akm/config"
	"akm/http"
	net "net/http"

	"github.com/gorilla/mux"
)

func setupRouter(service *http.ServiceOps) *mux.Router {
	router := mux.NewRouter()
	http.SetupRoutes(router, service)

	return router
}

func main() {
	fmt.Println("Hello, World!")
	store.InitDB()
	Migrate()

	repoOps := dbops.NewOpsManager(store.DataBase)
	service := http.NewServiceOps(repoOps)
	router := setupRouter(service)

	// Start server
	server := &net.Server{
		Addr:         ":" + config.GetEnv("SERVER_PORT", "8080"),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", config.GetEnv("SERVER_PORT", "8080"))
		if err := server.ListenAndServe(); err != nil && err != net.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

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
