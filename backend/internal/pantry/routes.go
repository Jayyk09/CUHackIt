package pantry

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func RegisterRoutes(r *http.ServeMux, db *database.DB) {
	cfg := config.GetConfig()
	_ = &pantryHandler{
		db:  db,
		log: logger.GetLogger(cfg.Log.Level),
		cfg: cfg,
	}

	// TODO: replace stubs with real handlers
	r.HandleFunc("GET /pantry", http.NotFound)
	r.HandleFunc("POST /pantry", http.NotFound)
	r.HandleFunc("GET /pantry/{id}", http.NotFound)
	r.HandleFunc("PATCH /pantry/{id}", http.NotFound)
	r.HandleFunc("DELETE /pantry/{id}", http.NotFound)
}
