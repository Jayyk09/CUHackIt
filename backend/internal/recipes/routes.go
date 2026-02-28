package recipes

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

// RegisterRoutes registers all recipe routes
func RegisterRoutes(r *http.ServeMux, db *database.DB, geminiClient *gemini.Client, log *logger.Logger) {
	h := NewHandler(db, geminiClient, log)

	// Recipe generation
	r.HandleFunc("POST /users/{user_id}/recipes/generate", h.GenerateRecipes)

	// Recipe CRUD
	r.HandleFunc("GET /users/{user_id}/recipes", h.ListRecipes)
	r.HandleFunc("POST /users/{user_id}/recipes", h.SaveRecipe)
	r.HandleFunc("GET /users/{user_id}/recipes/{id}", h.GetRecipe)
	r.HandleFunc("PUT /users/{user_id}/recipes/{id}", h.UpdateRecipe)
	r.HandleFunc("DELETE /users/{user_id}/recipes/{id}", h.DeleteRecipe)

	// Recipe actions
	r.HandleFunc("POST /users/{user_id}/recipes/{id}/favorite", h.ToggleFavorite)
	r.HandleFunc("POST /users/{user_id}/recipes/{id}/cooked", h.MarkAsCooked)
}
