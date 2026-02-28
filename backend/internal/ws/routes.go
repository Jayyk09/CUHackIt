package ws

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/pantry"
	"github.com/Jayyk09/CUHackIt/internal/users"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

// RegisterRoutes registers the WebSocket endpoint and starts the hub
func RegisterRoutes(r *http.ServeMux, db *database.DB, geminiClient *gemini.Client, log *logger.Logger) *Hub {
	pantryRepo := pantry.NewRepository(db.Pool)
	userRepo := users.NewRepository(db.Pool)

	hub := NewHub(geminiClient, pantryRepo, userRepo, log)

	// Start the hub in a goroutine
	go hub.Run()

	// Register WebSocket endpoint
	r.HandleFunc("GET /ws", hub.HandleWebSocket)
	r.HandleFunc("GET /ws/recipes", hub.HandleWebSocket)

	log.Info("WebSocket endpoint registered at /ws")

	return hub
}
