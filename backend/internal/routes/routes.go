package routes

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/handlers"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func Setup(r *http.ServeMux, db *database.DB, l *logger.Logger) {
	pantry := handlers.NewPantryHandler(db, l)
	users := handlers.NewUserHandler(db, l)
	food := handlers.NewFoodHandler(db, l)

	// Pantry routes
	r.HandleFunc("GET /pantry", pantry.GetAll)
	r.HandleFunc("GET /pantry/{id}", pantry.GetByID)
	r.HandleFunc("POST /pantry", pantry.Create)
	r.HandleFunc("PATCH /pantry/{id}", pantry.Update)
	r.HandleFunc("DELETE /pantry/{id}", pantry.Delete)

	// User routes
	r.HandleFunc("GET /users/{id}", users.GetByID)
	r.HandleFunc("POST /users", users.Create)
	r.HandleFunc("PUT /users/{id}", users.Update)

	// Food routes
	r.HandleFunc("GET /food/{id}", food.GetByID)
	r.HandleFunc("POST /food", food.Create)
}
