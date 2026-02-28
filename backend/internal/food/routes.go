package food

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func RegisterRoutes(r *http.ServeMux, db *database.DB) {
	cfg := config.GetConfig()
	_ = &foodHandler{
		db:  db,
		log: logger.GetLogger(cfg.Log.Level),
		cfg: cfg,
	}

	// TODO: replace stubs with real handlers
	r.HandleFunc("GET /food/{id}", http.NotFound)
	r.HandleFunc("POST /food/{id}", http.NotFound)
}
