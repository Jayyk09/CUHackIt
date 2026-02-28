package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

// Handler handles HTTP requests for users
type Handler struct {
	repo *Repository
	log  *logger.Logger
}

// NewHandler creates a new user handler
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

// GetUser handles GET /users/{id}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	user, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("Failed to get user: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// GetUserByAuth0ID handles GET /users/auth0/{auth0_id}
func (h *Handler) GetUserByAuth0ID(w http.ResponseWriter, r *http.Request) {
	auth0ID := r.PathValue("auth0_id")
	if auth0ID == "" {
		h.writeError(w, http.StatusBadRequest, "missing auth0 id")
		return
	}

	user, err := h.repo.GetByAuth0ID(r.Context(), auth0ID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("Failed to get user by auth0 id: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// CreateUser handles POST /users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Auth0ID == "" || input.Email == "" {
		h.writeError(w, http.StatusBadRequest, "auth0_id and email are required")
		return
	}

	user, err := h.repo.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			h.writeError(w, http.StatusConflict, "user already exists")
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			h.writeError(w, http.StatusBadRequest, "invalid input")
			return
		}
		h.log.Error("Failed to create user: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusCreated, user)
}

// FindOrCreateUser handles POST /users/find-or-create
func (h *Handler) FindOrCreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Auth0ID == "" || input.Email == "" {
		h.writeError(w, http.StatusBadRequest, "auth0_id and email are required")
		return
	}

	user, err := h.repo.FindOrCreate(r.Context(), input)
	if err != nil {
		h.log.Error("Failed to find or create user: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// UpdateProfile handles PUT /users/{id}/profile
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	var input UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.repo.UpdateProfile(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("Failed to update profile: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// CompleteOnboarding handles POST /users/{id}/onboarding
func (h *Handler) CompleteOnboarding(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	var input UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.repo.CompleteOnboarding(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("Failed to complete onboarding: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// DeleteUser handles DELETE /users/{id}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("Failed to delete user: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCurrentUser handles GET /users/me - gets user from session/context
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get auth0_id from context (set by auth middleware)
	auth0ID, ok := r.Context().Value("auth0_id").(string)
	if !ok || auth0ID == "" {
		h.writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, err := h.repo.GetByAuth0ID(r.Context(), auth0ID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("Failed to get current user: %v", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}
