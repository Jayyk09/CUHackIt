package routes

import (
	"fmt"
	"net/http"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/auth"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/food"
	"github.com/Jayyk09/CUHackIt/internal/pantry"
	"github.com/Jayyk09/CUHackIt/internal/users"
)

func Setup(r *http.ServeMux, db *database.DB, cfg *config.Config) error {
	store := auth.NewSessionStore(cfg)
	pantry.RegisterRoutes(r, db, store)
	users.RegisterRoutes(r, db, store)
	food.RegisterRoutes(r, db, store)

	if err := auth.RegisterRoutes(r, cfg, store); err != nil {
		return fmt.Errorf("auth routes: %w", err)
	}

	return nil
}
