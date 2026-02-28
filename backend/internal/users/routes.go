package users

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

// RegisterRoutes registers all user routes
func RegisterRoutes(r *http.ServeMux, db *database.DB, log *logger.Logger) {
	h := NewHandler(db, log)

	// User CRUD
	r.HandleFunc("GET /users/{id}", h.GetUser)
	r.HandleFunc("POST /users", h.CreateUser)
	r.HandleFunc("DELETE /users/{id}", h.DeleteUser)

	// Auth0 ID lookup
	r.HandleFunc("GET /auth0-users/{auth0_id}", h.GetUserByAuth0ID)

	// Find or create (for Auth0 callback)
	r.HandleFunc("POST /users/find-or-create", h.FindOrCreateUser)

	// Profile management
	r.HandleFunc("PUT /users/{id}/profile", h.UpdateProfile)

	// Onboarding
	r.HandleFunc("POST /users/{id}/onboarding", h.CompleteOnboarding)

	// Current user (requires auth)
	r.HandleFunc("GET /users/me", h.GetCurrentUser)
}
