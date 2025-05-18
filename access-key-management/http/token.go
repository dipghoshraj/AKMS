package http

import (
	"akm/dbops/model"
	"encoding/json"
	"net/http"
)

func (s *ServiceOps) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var tokenRequest model.TokenCreateInput
	if err := json.NewDecoder(r.Body).Decode(&tokenRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the request
	// if err := validateTokenRequest(tokenRequest); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

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
