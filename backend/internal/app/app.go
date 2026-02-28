package app

import (
	"context"
	"net/http"

	"github.com/Jayyk09/CUHackIt/cmd/server"
	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
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

	routes.Setup(r, db, cfg)

	if err := server.Start(cfg, r, l); err != nil {
		l.Fatal("Server error: %v", err)
	}
}
