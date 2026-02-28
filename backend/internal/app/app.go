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

	// Run auto-migrations (creates tables if they don't exist).
	if err := db.Migrate(context.Background()); err != nil {
		l.Fatal("Migration error: %v", err)
	}
	l.Info("Database migrations applied")

	// Create basic router
	r := http.NewServeMux()

	if err := routes.Setup(r, db, cfg); err != nil {
		l.Fatal("Failed to setup routes: %v", err)
	}

	if err := server.Start(cfg, middleware.CORS(cfg, r), l); err != nil {
		l.Fatal("Server error: %v", err)
	}
}
