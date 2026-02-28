package pantry

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
	_ = &pantryHandler{
		db:  db,
		log: logger.GetLogger(cfg.Log.Level),
		cfg: cfg,
	}

	// TODO: replace stubs with real handlers
	r.Handle("GET /pantry", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
	r.Handle("POST /pantry", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
	r.Handle("GET /pantry/{id}", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
	r.Handle("PATCH /pantry/{id}", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
	r.Handle("DELETE /pantry/{id}", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
}
