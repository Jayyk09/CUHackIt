package recipes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Jayyk09/CUHackIt/internal/agents"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/pantry"
	"github.com/Jayyk09/CUHackIt/internal/users"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for recipes
type Handler struct {
	repo         *Repository
	pantryRepo   *pantry.Repository
	userRepo     *users.Repository
	orchestrator *agents.Orchestrator
	log          *logger.Logger
}

// NewHandler creates a new recipe handler
func NewHandler(db *database.DB, geminiClient *gemini.Client, log *logger.Logger) *Handler {
	var orchestrator *agents.Orchestrator
	if geminiClient != nil {
		orchestrator = agents.NewOrchestrator(geminiClient, log)
	}

	return &Handler{
		repo:         NewRepository(db.Pool),
		pantryRepo:   pantry.NewRepository(db.Pool),
		userRepo:     users.NewRepository(db.Pool),
		orchestrator: orchestrator,
		log:          log,
	}
}

// writeJSON writes a JSON response
func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("Failed to encode response: %v", err)
	}
}

// writeError writes an error response
func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

// getUserID extracts user ID from request
func (h *Handler) getUserID(r *http.Request) (string, error) {
	if idStr := r.PathValue("user_id"); idStr != "" {
		return idStr, nil
	}
	return "", errors.New("user id not found")
}

// GenerateRecipesRequest is the request body for generating recipes
type GenerateRecipesRequest struct {
	Mode        string `json:"mode"`         // "pantry_only", "flexible", "both"
	RecipeCount int    `json:"recipe_count"` // 1-3
}

// GenerateRecipes handles POST /users/{user_id}/recipes/generate
func (h *Handler) GenerateRecipes(w http.ResponseWriter, r *http.Request) {
	if h.orchestrator == nil {
		h.writeError(w, http.StatusServiceUnavailable, "recipe generation not available - Gemini API not configured")
		return
	}

	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req GenerateRecipesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Use defaults if no body provided
		req.Mode = "pantry_only"
		req.RecipeCount = 2
	}

	if req.RecipeCount <= 0 {
		req.RecipeCount = 2
	}
	if req.RecipeCount > 3 {
		req.RecipeCount = 3
	}

	// Get user's pantry items
	pantryItems, err := h.pantryRepo.ListByUserID(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get pantry items: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get pantry items")
		return
	}

	if len(pantryItems) == 0 {
		h.writeError(w, http.StatusBadRequest, "pantry is empty - add some items first")
		return
	}

	// Get user's preferences
	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get user: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get user preferences")
		return
	}

	// Convert pantry items to agent format
	agentPantryItems := make([]agents.PantryItem, len(pantryItems))
	for i, item := range pantryItems {
		category := ""
		if item.Category != nil {
			category = *item.Category
		}
		agentPantryItems[i] = agents.PantryItem{
			ID:       strconv.Itoa(item.ID),
			Name:     item.ProductName,
			Category: category,
			Quantity: float64(item.Quantity),
			Unit:     "item",
		}
	}

	// Determine mode
	var mode agents.OrchestratorMode
	switch req.Mode {
	case "flexible":
		mode = agents.ModeFlexible
	case "both":
		mode = agents.ModeBoth
	default:
		mode = agents.ModePantryOnly
	}

	// Generate recipes
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	result, err := h.orchestrator.Generate(ctx, agents.GenerateRequest{
		RecipeRequest: agents.RecipeRequest{
			PantryItems:        agentPantryItems,
			Allergens:          user.Allergens,
			DietaryPreferences: user.DietaryPreferences,
			NutritionalGoals:   user.NutritionalGoals,
			CookingSkill:       user.CookingSkill,
			CuisinePreferences: user.CuisinePreferences,
			RecipeCount:        req.RecipeCount,
		},
		Mode: mode,
	})
	if err != nil {
		if errors.Is(err, agents.ErrAllRecipesFiltered) {
			h.writeJSON(w, http.StatusOK, map[string]interface{}{
				"message":        "all generated recipes contained allergens and were filtered",
				"recipes":        []interface{}{},
				"filtered_count": result.FilteredCount,
			})
			return
		}
		h.log.Error("Failed to generate recipes: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to generate recipes")
		return
	}

	h.writeJSON(w, http.StatusOK, result)
}

// SaveRecipe handles POST /users/{user_id}/recipes
func (h *Handler) SaveRecipe(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var input CreateRecipeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Title == "" {
		h.writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	recipe, err := h.repo.Create(r.Context(), userID, input)
	if err != nil {
		h.log.Error("Failed to save recipe: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to save recipe")
		return
	}

	h.writeJSON(w, http.StatusCreated, recipe)
}

// ListRecipes handles GET /users/{user_id}/recipes
func (h *Handler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// Check for favorites filter
	favoritesOnly := r.URL.Query().Get("favorites") == "true"

	var recipes []Recipe
	if favoritesOnly {
		recipes, err = h.repo.ListFavorites(r.Context(), userID)
	} else {
		recipes, err = h.repo.ListByUserID(r.Context(), userID)
	}

	if err != nil {
		h.log.Error("Failed to list recipes: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to list recipes")
		return
	}

	if recipes == nil {
		recipes = []Recipe{}
	}

	h.writeJSON(w, http.StatusOK, recipes)
}

// GetRecipe handles GET /users/{user_id}/recipes/{id}
func (h *Handler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	recipeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	recipe, err := h.repo.GetByID(r.Context(), recipeID)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			h.writeError(w, http.StatusNotFound, "recipe not found")
			return
		}
		h.log.Error("Failed to get recipe: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}

	// Verify ownership
	if recipe.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	h.writeJSON(w, http.StatusOK, recipe)
}

// UpdateRecipe handles PUT /users/{user_id}/recipes/{id}
func (h *Handler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	recipeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	// Verify ownership
	existing, err := h.repo.GetByID(r.Context(), recipeID)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			h.writeError(w, http.StatusNotFound, "recipe not found")
			return
		}
		h.log.Error("Failed to get recipe: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}

	if existing.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	var input UpdateRecipeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	recipe, err := h.repo.Update(r.Context(), recipeID, input)
	if err != nil {
		h.log.Error("Failed to update recipe: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to update recipe")
		return
	}

	h.writeJSON(w, http.StatusOK, recipe)
}

// ToggleFavorite handles POST /users/{user_id}/recipes/{id}/favorite
func (h *Handler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	recipeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	// Verify ownership and get current state
	existing, err := h.repo.GetByID(r.Context(), recipeID)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			h.writeError(w, http.StatusNotFound, "recipe not found")
			return
		}
		h.log.Error("Failed to get recipe: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}

	if existing.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	// Toggle favorite
	newFavorite := !existing.IsFavorite
	recipe, err := h.repo.Update(r.Context(), recipeID, UpdateRecipeInput{
		IsFavorite: &newFavorite,
	})
	if err != nil {
		h.log.Error("Failed to toggle favorite: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to toggle favorite")
		return
	}

	h.writeJSON(w, http.StatusOK, recipe)
}

// MarkAsCooked handles POST /users/{user_id}/recipes/{id}/cooked
func (h *Handler) MarkAsCooked(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	recipeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	// Verify ownership
	existing, err := h.repo.GetByID(r.Context(), recipeID)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			h.writeError(w, http.StatusNotFound, "recipe not found")
			return
		}
		h.log.Error("Failed to get recipe: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}

	if existing.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	recipe, err := h.repo.MarkAsCooked(r.Context(), recipeID)
	if err != nil {
		h.log.Error("Failed to mark as cooked: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to mark as cooked")
		return
	}

	h.writeJSON(w, http.StatusOK, recipe)
}

// DeleteRecipe handles DELETE /users/{user_id}/recipes/{id}
func (h *Handler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	recipeID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	// Verify ownership
	existing, err := h.repo.GetByID(r.Context(), recipeID)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			h.writeError(w, http.StatusNotFound, "recipe not found")
			return
		}
		h.log.Error("Failed to get recipe: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get recipe")
		return
	}

	if existing.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	if err := h.repo.Delete(r.Context(), recipeID); err != nil {
		h.log.Error("Failed to delete recipe: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to delete recipe")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
