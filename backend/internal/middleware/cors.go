package middleware

import (
	"net/http"
	"strings"

	"github.com/Jayyk09/CUHackIt/config"
)

// CORS wraps an http.Handler and adds the necessary CORS headers so the
// Next.js frontend (on a different port) can call the API.
func CORS(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestOrigin := strings.TrimSpace(r.Header.Get("Origin"))
		frontendOrigin := strings.TrimRight(strings.TrimSpace(cfg.App.FrontendURL), "/")
		allowedOrigins := map[string]struct{}{}
		if frontendOrigin != "" {
			allowedOrigins[frontendOrigin] = struct{}{}
		}
		allowedOrigins["http://localhost:3000"] = struct{}{}
		allowedOrigins["http://127.0.0.1:3000"] = struct{}{}

		if requestOrigin != "" {
			if _, ok := allowedOrigins[requestOrigin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
			}
		}
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Handle preflight.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
