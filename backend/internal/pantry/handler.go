package pantry

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for pantry items
type Handler struct {
	repo *Repository
	log  *logger.Logger
}

// NewHandler creates a new pantry handler
func NewHandler(db *database.DB, log *logger.Logger) *Handler {
	return &Handler{
		repo: NewRepository(db.Pool),
		log:  log,
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

// getUserID extracts user ID from request (path param or context)
func (h *Handler) getUserID(r *http.Request) (string, error) {
	if idStr := r.PathValue("user_id"); idStr != "" {
		return idStr, nil
	}
	if userID, ok := r.Context().Value("user_id").(string); ok {
		return userID, nil
	}
	return "", errors.New("user id not found")
}

// ListItems handles GET /users/{user_id}/pantry
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// Check for category filter
	category := r.URL.Query().Get("category")
	var items []PantryItem

	if category != "" {
		items, err = h.repo.ListByCategory(r.Context(), userID, FoodCategory(category))
	} else {
		items, err = h.repo.ListByUserID(r.Context(), userID)
	}

	if err != nil {
		h.log.Error("Failed to list pantry items: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Return empty array instead of null
	if items == nil {
		items = []PantryItem{}
	}

	h.writeJSON(w, http.StatusOK, items)
}

// GetItem handles GET /users/{user_id}/pantry/{id}
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	itemID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := h.repo.GetByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			h.writeError(w, http.StatusNotFound, "item not found")
			return
		}
		h.log.Error("Failed to get pantry item: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Verify ownership
	if item.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	h.writeJSON(w, http.StatusOK, item)
}

// CreateItem handles POST /users/{user_id}/pantry
func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var input CreateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Name == "" {
		h.writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	if input.Category == "" {
		h.writeError(w, http.StatusBadRequest, "category is required")
		return
	}

	item, err := h.repo.Create(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			h.writeError(w, http.StatusBadRequest, "invalid input")
			return
		}
		h.log.Error("Failed to create pantry item: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusCreated, item)
}

// UpdateItem handles PUT /users/{user_id}/pantry/{id}
func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	itemID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	// Verify ownership first
	existing, err := h.repo.GetByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			h.writeError(w, http.StatusNotFound, "item not found")
			return
		}
		h.log.Error("Failed to get pantry item: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if existing.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	var input UpdateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.repo.Update(r.Context(), itemID, input)
	if err != nil {
		h.log.Error("Failed to update pantry item: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, item)
}

// DeleteItem handles DELETE /users/{user_id}/pantry/{id}
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	itemID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	// Verify ownership first
	existing, err := h.repo.GetByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			h.writeError(w, http.StatusNotFound, "item not found")
			return
		}
		h.log.Error("Failed to get pantry item: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if existing.UserID != userID {
		h.writeError(w, http.StatusForbidden, "access denied")
		return
	}

	if err := h.repo.Delete(r.Context(), itemID); err != nil {
		h.log.Error("Failed to delete pantry item: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListExpiringSoon handles GET /users/{user_id}/pantry/expiring
func (h *Handler) ListExpiringSoon(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	items, err := h.repo.ListExpiringSoon(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to list expiring items: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if items == nil {
		items = []PantryItem{}
	}

	h.writeJSON(w, http.StatusOK, items)
}

// GetCategorySummary handles GET /users/{user_id}/pantry/summary
func (h *Handler) GetCategorySummary(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	summary, err := h.repo.GetCategorySummary(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get category summary: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, summary)
}

// AddToPantry handles POST /pantry
func (h *Handler) AddToPantry(w http.ResponseWriter, r *http.Request) {
	var input AddToPantryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Auth0ID == "" {
		h.writeError(w, http.StatusBadRequest, "auth0_id is required")
		return
	}

	if input.FoodID <= 0 {
		h.writeError(w, http.StatusBadRequest, "food_id is required")
		return
	}

	entry, err := h.repo.AddToPantry(r.Context(), input)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			h.writeError(w, http.StatusBadRequest, "invalid input")
			return
		}
		h.log.Error("Failed to add pantry entry: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusCreated, entry)
}
