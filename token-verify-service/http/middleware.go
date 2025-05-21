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

		valid, limit, err := s.tokenOps.GetRedisToken(r.Context(), authHeader)
		if !valid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		rateExceed, err := s.tokenOps.CheckRateLimit(r.Context(), authHeader, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		if rateExceed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
