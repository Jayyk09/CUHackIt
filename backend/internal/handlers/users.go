package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/models"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

type UserHandler struct {
	db  *database.DB
	log logger.Interface
}

func NewUserHandler(db *database.DB, log logger.Interface) *UserHandler {
	return &UserHandler{db: db, log: log}
}

// GetByID handles GET /users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var account models.Account
	err = h.db.Pool.QueryRow(r.Context(), `
		SELECT account_id, labels, allergens
		FROM account
		WHERE account_id = $1
	`, id).Scan(&account.AccountID, &account.Labels, &account.Allergens)
	if err != nil {
		h.log.Error("failed to get user %d: %v", id, err)
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	if account.Labels == nil {
		account.Labels = []string{}
	}
	if account.Allergens == nil {
		account.Allergens = []string{}
	}

	writeJSON(w, http.StatusOK, account)
}

// Create handles POST /users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Labels == nil {
		req.Labels = []string{}
	}
	if req.Allergens == nil {
		req.Allergens = []string{}
	}

	var id int
	err := h.db.Pool.QueryRow(r.Context(), `
		INSERT INTO account (labels, allergens)
		VALUES ($1, $2)
		RETURNING account_id
	`, req.Labels, req.Allergens).Scan(&id)
	if err != nil {
		h.log.Error("failed to create user: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int{"account_id": id})
}

// Update handles PUT /users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req models.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Labels == nil {
		req.Labels = []string{}
	}
	if req.Allergens == nil {
		req.Allergens = []string{}
	}

	ct, err := h.db.Pool.Exec(r.Context(), `
		UPDATE account
		SET labels = $1, allergens = $2
		WHERE account_id = $3
	`, req.Labels, req.Allergens, id)
	if err != nil {
		h.log.Error("failed to update user %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	if ct.RowsAffected() == 0 {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
