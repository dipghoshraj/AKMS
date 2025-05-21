package http

import (
	"encoding/json"
	"net/http"
)

func (s *Server) getPrice(w http.ResponseWriter, r *http.Request) {

	payload := map[string]any{
		"success":   true,
		"symbol":    "ETH",
		"price_usd": 3482.42,
		"timestamp": "2025-05-20T13:22:00Z",
	}
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}
