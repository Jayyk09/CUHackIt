package users

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func RegisterRoutes(r *http.ServeMux, db *database.DB) {
	cfg := config.GetConfig()
	_ = &userHandler{
		db:  db,
		log: logger.GetLogger(cfg.Log.Level),
		cfg: cfg,
	}

	// TODO: replace stubs with real handlers
	r.HandleFunc("GET /users/{id}", http.NotFound)
	r.HandleFunc("POST /users", http.NotFound)
	r.HandleFunc("PUT /users/{id}", http.NotFound)
}
