package users

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func RegisterRoutes(r *http.ServeMux, db *database.DB) {
	cfg := config.GetConfig()
	h := &userHandler{
		db:  db,
		log: logger.GetLogger(cfg.Log.Level),
		cfg: cfg,
	}

	r.HandleFunc("GET /users/{id}", _)
	r.HandleFunc("POST /users", _)
	r.HandleFunc("PUT /users/{id}", _)
}
