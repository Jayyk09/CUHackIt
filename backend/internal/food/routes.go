package food

import (
	"context"
	"net/http"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/auth"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
	"github.com/gorilla/sessions"
)

// RegisterRoutes registers all food routes
func RegisterRoutes(r *http.ServeMux, db *database.DB, store sessions.Store) {
	cfg := config.GetConfig()
	log := logger.GetLogger(cfg.Log.Level)
	aiClient, err := gemini.NewClient(context.Background(), cfg.Gemini.APIKey, cfg.Gemini.Model, log)
	if err != nil {
		log.Warn("failed to initialize gemini client: %v", err)
	}

	h := NewHandler(db, aiClient, cfg, log)
	r.Handle("GET /food", auth.IsAuthenticated(store, http.HandlerFunc(h.List)))
	r.Handle("GET /food/{id}", auth.IsAuthenticated(store, http.HandlerFunc(h.GetProduct)))
	r.Handle("PATCH /food/{id}/metadata", auth.IsAuthenticated(store, http.HandlerFunc(h.UpdateMetadata)))
}
