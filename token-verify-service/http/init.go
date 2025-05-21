package http

import (
	"tvs/dbops"

	"github.com/gorilla/mux"
)

type Server struct {
	tokenOps dbops.TokenOps
}

func NewServerOps(ops *dbops.OpsManager) *Server {
	return &Server{
		tokenOps: ops.TokenOps,
	}
}

func SetupRoutes(router *mux.Router, service *Server) {
	router.Use(RateMiddleware)
	router.HandleFunc("/price", service.getPrice).Methods("GET")
}
