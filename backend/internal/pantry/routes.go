package pantry

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func RegisterRoutes(r *http.ServeMux, db *database.DB) {
	cfg := config.GetConfig()
	h := &pantryHandler{
		db:  db,
		log: logger.GetLogger(cfg.Log.Level),
		cfg: cfg,
	}

	r.HandleFunc("GET /pantry", _)
	r.HandleFunc("POST /pantry", _)
	r.HandleFunc("GET /pantry/{id}", _)
	r.HandleFunc("PATCH /pantry/{id}", _)
	r.HandleFunc("DELETE /pantry/{id}", _)
}
