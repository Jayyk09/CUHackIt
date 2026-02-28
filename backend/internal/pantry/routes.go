package pantry

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

// RegisterRoutes registers all pantry routes
func RegisterRoutes(r *http.ServeMux, db *database.DB, log *logger.Logger) {
	h := NewHandler(db, log)

	// Pantry item reads (nested under users)
	r.HandleFunc("GET /users/{user_id}/pantry", h.ListItems)
	r.HandleFunc("GET /users/{user_id}/pantry/{id}", h.GetItem)
	r.HandleFunc("DELETE /users/{user_id}/pantry/{id}", h.DeleteItem)

	// Category summary
	r.HandleFunc("GET /users/{user_id}/pantry/summary", h.GetCategorySummary)

	// Simplified pantry endpoint (uses auth0_id to resolve user)
	r.HandleFunc("POST /pantry", h.AddToPantry)
}
