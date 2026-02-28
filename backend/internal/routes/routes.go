package routes

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/food"
	"github.com/Jayyk09/CUHackIt/internal/pantry"
	"github.com/Jayyk09/CUHackIt/internal/users"
)

func Setup(r *http.ServeMux, db *database.DB) {
	pantry.RegisterRoutes(r, db)
	users.RegisterRoutes(r, db)
	food.RegisterRoutes(r, db)
}
