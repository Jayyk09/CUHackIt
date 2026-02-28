package users

import (
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/auth"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func RegisterRoutes(r *http.ServeMux, db *database.DB, store sessions.Store) {
	cfg := config.GetConfig()
	_ = &userHandler{
		db:  db,
		log: logger.GetLogger(cfg.Log.Level),
		cfg: cfg,
	}

	// TODO: replace stubs with real handlers
	r.Handle("GET /users/{id}", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
	r.Handle("POST /users", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
	r.Handle("PUT /users/{id}", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
}
