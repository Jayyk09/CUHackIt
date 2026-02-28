package app

import (
	"context"
	"net/http"

	"github.com/Jayyk09/CUHackIt/cmd/server"
	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/middleware"
	"github.com/Jayyk09/CUHackIt/internal/routes"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.GetLogger(cfg.Log.Level)

	l.Info("Starting application...")

	// Connect to database
	db, err := database.New(context.Background(), cfg.DB.URL)
	if err != nil {
		l.Fatal("Database connection error: %v", err)
	}
	defer db.Close()

	l.Info("Connected to database")

	// Create basic router
	r := http.NewServeMux()

	// Setup routes
	if err := routes.Setup(r, db, cfg, l); err != nil {
		l.Fatal("Failed to setup routes: %v", err)
	}

	l.Info("Routes registered")

	handler := middleware.CORS(cfg, r)

	if err := server.Start(cfg, handler, l); err != nil {
		l.Fatal("Server error: %v", err)
	}
}
