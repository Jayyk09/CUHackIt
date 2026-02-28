package users

import (
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func RegisterRoutes(r *http.ServeMux, db *database.DB, store sessions.Store) {
	cfg := config.GetConfig()
	h := &userHandler{
		db:  db,
		log: logger.GetLogger(cfg.Log.Level),
		cfg: cfg,
	}

	r.HandleFunc("GET /users/{id}", h.GetByID)
	r.HandleFunc("POST /users", h.Create)
	r.HandleFunc("PUT /users/{id}", h.UpdatePrefs)
}
