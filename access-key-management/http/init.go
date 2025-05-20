package http

import (
	"akm/dbops"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ServiceOps struct {
	tokenOps dbops.TokenOps
}

func NewServiceOps(ops *dbops.OpsManager) *ServiceOps {
	return &ServiceOps{
		tokenOps: ops.TokenOps,
	}
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func SetupRoutes(router *mux.Router, service *ServiceOps) {
	// Initialize the router
	// Define your routes here

	router.HandleFunc("/key/info", service.getTokenPlan).Methods("GET")
	router.HandleFunc("/key/disable", service.disableTokenHandler).Methods("POST")

	protected := router.PathPrefix("/admin").Subrouter()
	protected.Use(JWTMiddleware)
	protected.HandleFunc("/key", service.createTokenHandler).Methods("POST")
	protected.HandleFunc("/key", service.getTokensHandler).Methods("GET")
	protected.HandleFunc("/keys/{key}", service.updateTokenHandler).Methods("PUT")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, APIResponse{
		Success: false,
		Error:   message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
