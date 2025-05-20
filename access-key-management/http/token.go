package http

import (
	"akm/dbops/model"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *ServiceOps) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var tokenRequest model.TokenCreateInput
	if err := json.NewDecoder(r.Body).Decode(&tokenRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := s.tokenOps.Create(r.Context(), &tokenRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	// Return the token in the response
	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    token,
	})
}

func (s *ServiceOps) getTokensHandler(w http.ResponseWriter, r *http.Request) {
	// Get the tokens from the database
	tokens, err := s.tokenOps.GetAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve tokens")
		return
	}

	// Return the tokens in the response
	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    tokens,
	})
}

func (s *ServiceOps) getTokenPlan(w http.ResponseWriter, r *http.Request) {
	// Get the token key from the URL parameters
	accessKey := r.Header.Get("x-api-key")

	// Retrieve the token plan based on the key
	token, err := s.tokenOps.GetByKey(r.Context(), accessKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve token plan")
		return
	}

	// Return the token plan in the response
	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    token,
	})
}

func (s *ServiceOps) disableTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Get the token key from the URL parameters
	accessKey := r.Header.Get("x-api-key")

	// Disable the token in the database
	err := s.tokenOps.Disable(r.Context(), accessKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to disable token")
		return
	}

	// Return a success message in the response
	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Token disabled successfully",
	})
}

func (s *ServiceOps) updateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Get the token key from the URL parameters
	accessKey := mux.Vars(r)["key"]
	if accessKey == "" {
		respondWithError(w, http.StatusInternalServerError, "Token key is required")
		return
	}

	var tokenRequest model.TokenUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&tokenRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := s.tokenOps.Update(r.Context(), accessKey, &tokenRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update token")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    token,
	})
}
