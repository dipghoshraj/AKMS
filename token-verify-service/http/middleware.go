package http

import (
	"net/http"
)

func (s *Server) TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("x-api-key")
		if authHeader == "" {
			http.Error(w, "Missing x-api-key header", http.StatusUnauthorized)
			return
		}

		valid, err := s.tokenOps.GetRedisToken(r.Context(), authHeader)
		if !valid {
			if err != nil && err.Error() == "rate limit exceeded" {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
