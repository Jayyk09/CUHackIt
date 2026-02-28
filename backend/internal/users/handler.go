package users

import (
	"encoding/json"
	"net/http"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

type userHandler struct {
	db  *database.DB
	log *logger.Logger
	cfg *config.Config
}

// GetByID returns a user by their Auth0 sub (the URL path value {id}).
func (h *userHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	user, err := GetUser(r.Context(), h.db, id)
	if err != nil {
		h.log.Error("get user: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Create creates a new user row. Expected JSON body: { "id", "email", "name" }.
func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	user, err := CreateUser(r.Context(), h.db, body.ID, body.Email, body.Name)
	if err != nil {
		h.log.Error("create user: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UpdatePrefs updates a user's labels and allergens.
// Expected JSON body: { "labels": [...], "allergens": [...] }
func (h *userHandler) UpdatePrefs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	var body struct {
		Labels    []string `json:"labels"`
		Allergens []string `json:"allergens"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user, err := UpdatePreferences(r.Context(), h.db, id, body.Labels, body.Allergens)
	if err != nil {
		h.log.Error("update preferences: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
