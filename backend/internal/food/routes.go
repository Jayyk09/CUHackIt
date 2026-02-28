package food

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/auth"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

func RegisterRoutes(r *http.ServeMux, db *database.DB, store sessions.Store) {
	cfg := config.GetConfig()
	log := logger.GetLogger(cfg.Log.Level)
	aiClient, err := gemini.New(context.Background(), cfg.Gemini.APIKey, cfg.Gemini.Model)
	if err != nil {
		log.Warn("failed to initialize gemini client: %v", err)
	}

	h := &foodHandler{
		db:  db,
		log: log,
		cfg: cfg,
		ai:  aiClient,
	}

	r.Handle("GET /food", http.HandlerFunc(h.List))
	r.Handle("GET /food/{id}", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
	r.Handle("POST /food/{id}", auth.IsAuthenticated(store, http.HandlerFunc(http.NotFound)))
}
