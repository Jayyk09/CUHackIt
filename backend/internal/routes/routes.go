package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/auth"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/food"
	"github.com/Jayyk09/CUHackIt/internal/pantry"
	"github.com/Jayyk09/CUHackIt/internal/recipes"
	"github.com/Jayyk09/CUHackIt/internal/users"
	"github.com/Jayyk09/CUHackIt/internal/ws"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

// Setup registers all application routes
func Setup(r *http.ServeMux, db *database.DB, cfg *config.Config, log *logger.Logger) error {
	// Initialize Gemini client (optional - can work without it)
	var geminiClient *gemini.Client
	if cfg.Gemini.APIKey != "" {
		var err error
		geminiClient, err = gemini.NewClient(context.Background(), cfg.Gemini.APIKey, cfg.Gemini.Model, log)
		if err != nil {
			log.Warn("Failed to initialize Gemini client: %v - recipe generation will be disabled", err)
		} else {
			log.Info("Gemini client initialized with model: %s", cfg.Gemini.Model)
		}
	} else {
		log.Warn("GEMINI_API_KEY not set - recipe generation will be disabled")
	}

	// User routes
	users.RegisterRoutes(r, db, log)

	// Pantry routes
	pantry.RegisterRoutes(r, db, log)

	// Food search routes (no DB needed - mock data)
	food.RegisterRoutes(r, log)

	// Recipe routes (with optional Gemini client)
	recipes.RegisterRoutes(r, db, geminiClient, log)

	// WebSocket routes for real-time recipe streaming
	ws.RegisterRoutes(r, db, geminiClient, log)

	// Auth routes
	if err := auth.RegisterRoutes(r, cfg); err != nil {
		return fmt.Errorf("auth routes: %w", err)
	}

	// Health check
	r.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API info
	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"name": "Sift API",
			"version": "1.0.0",
			"description": "AI-powered pantry management and recipe generation",
			"endpoints": {
				"health": "GET /health",
				"users": "GET/POST /users",
				"pantry": "GET/POST /users/{user_id}/pantry",
				"recipes": "GET/POST /users/{user_id}/recipes",
				"generate": "POST /users/{user_id}/recipes/generate",
				"food_search": "GET /food/search?q=...",
				"websocket": "GET /ws (real-time recipe streaming)"
			}
		}`))
	})

	return nil
}
