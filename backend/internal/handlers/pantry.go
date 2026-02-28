package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/models"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

type PantryHandler struct {
	db  *database.DB
	log logger.Interface
}

func NewPantryHandler(db *database.DB, log logger.Interface) *PantryHandler {
	return &PantryHandler{db: db, log: log}
}

// GetAll handles GET /pantry and GET /pantry?user_id=X&category=Y
// Returns the joined pantry+food response with is_spoiled computed.
func (h *PantryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Optional query params
	userIDStr := r.URL.Query().Get("user_id")
	category := r.URL.Query().Get("category")

	query := `
		SELECT
			f.product_name,
			f.environmental_score,
			f.nutriscore_score,
			f.labels_en,
			f.allergens_en,
			f.traces_en,
			f.image_url,
			f.shelf_life,
			f.category,
			p.quantity,
			p.units,
			p.is_frozen,
			p.added_at
		FROM pantry p
		JOIN food f ON p.food_id = f.id
		WHERE 1=1
	`

	args := []interface{}{}
	argIdx := 1

	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		query += " AND p.user_id = $" + strconv.Itoa(argIdx)
		args = append(args, userID)
		argIdx++
	}

	if category != "" {
		query += " AND p.category = $" + strconv.Itoa(argIdx)
		args = append(args, category)
		argIdx++
	}

	query += " ORDER BY p.added_at DESC"

	rows, err := h.db.Pool.Query(ctx, query, args...)
	if err != nil {
		h.log.Error("failed to query pantry: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to query pantry")
		return
	}
	defer rows.Close()

	results := []models.PantryResponse{}
	now := time.Now()

	for rows.Next() {
		var (
			productName        string
			environmentalScore float64
			nutriscoreScore    float64
			labelsEn           []string
			allergensEn        []string
			tracesEn           []string
			imageURL           *string
			shelfLife          *int
			category           *string
			quantity           int
			units              string
			isFrozen           bool
			addedAt            time.Time
		)

		err := rows.Scan(
			&productName,
			&environmentalScore,
			&nutriscoreScore,
			&labelsEn,
			&allergensEn,
			&tracesEn,
			&imageURL,
			&shelfLife,
			&category,
			&quantity,
			&units,
			&isFrozen,
			&addedAt,
		)
		if err != nil {
			h.log.Error("failed to scan pantry row: %v", err)
			writeError(w, http.StatusInternalServerError, "failed to read pantry data")
			return
		}

		// Compute is_spoiled: if shelf_life is set, check if added_at + shelf_life days < now
		isSpoiled := false
		if shelfLife != nil && *shelfLife > 0 {
			expiresAt := addedAt.AddDate(0, 0, *shelfLife)
			isSpoiled = now.After(expiresAt)
		}

		// Ensure arrays are not null in JSON
		if labelsEn == nil {
			labelsEn = []string{}
		}
		if allergensEn == nil {
			allergensEn = []string{}
		}
		if tracesEn == nil {
			tracesEn = []string{}
		}

		results = append(results, models.PantryResponse{
			ProductName:        productName,
			EnvironmentalScore: environmentalScore,
			NutriscoreScore:    nutriscoreScore,
			LabelsEn:           labelsEn,
			AllergensEn:        allergensEn,
			TracesEn:           tracesEn,
			ImageURL:           imageURL,
			ShelfLife:          shelfLife,
			Category:           category,
			Quantity:           quantity,
			Units:              units,
			IsSpoiled:          isSpoiled,
			IsFrozen:           isFrozen,
			AddedAt:            addedAt.Format("2006-01-02 15:04:05"),
		})
	}

	writeJSON(w, http.StatusOK, results)
}

// GetByID handles GET /pantry/{id}
func (h *PantryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pantry id")
		return
	}

	var (
		productName        string
		environmentalScore float64
		nutriscoreScore    float64
		labelsEn           []string
		allergensEn        []string
		tracesEn           []string
		imageURL           *string
		shelfLife          *int
		category           *string
		quantity           int
		units              string
		isFrozen           bool
		addedAt            time.Time
	)

	err = h.db.Pool.QueryRow(r.Context(), `
		SELECT
			f.product_name,
			f.environmental_score,
			f.nutriscore_score,
			f.labels_en,
			f.allergens_en,
			f.traces_en,
			f.image_url,
			f.shelf_life,
			f.category,
			p.quantity,
			p.units,
			p.is_frozen,
			p.added_at
		FROM pantry p
		JOIN food f ON p.food_id = f.id
		WHERE p.id = $1
	`, id).Scan(
		&productName,
		&environmentalScore,
		&nutriscoreScore,
		&labelsEn,
		&allergensEn,
		&tracesEn,
		&imageURL,
		&shelfLife,
		&category,
		&quantity,
		&units,
		&isFrozen,
		&addedAt,
	)
	if err != nil {
		h.log.Error("failed to get pantry item %d: %v", id, err)
		writeError(w, http.StatusNotFound, "pantry item not found")
		return
	}

	isSpoiled := false
	if shelfLife != nil && *shelfLife > 0 {
		expiresAt := addedAt.AddDate(0, 0, *shelfLife)
		isSpoiled = time.Now().After(expiresAt)
	}

	if labelsEn == nil {
		labelsEn = []string{}
	}
	if allergensEn == nil {
		allergensEn = []string{}
	}
	if tracesEn == nil {
		tracesEn = []string{}
	}

	result := models.PantryResponse{
		ProductName:        productName,
		EnvironmentalScore: environmentalScore,
		NutriscoreScore:    nutriscoreScore,
		LabelsEn:           labelsEn,
		AllergensEn:        allergensEn,
		TracesEn:           tracesEn,
		ImageURL:           imageURL,
		ShelfLife:          shelfLife,
		Category:           category,
		Quantity:           quantity,
		Units:              units,
		IsSpoiled:          isSpoiled,
		IsFrozen:           isFrozen,
		AddedAt:            addedAt.Format("2006-01-02 15:04:05"),
	}

	writeJSON(w, http.StatusOK, result)
}

// Create handles POST /pantry
func (h *PantryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePantryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	if req.Units == "" {
		req.Units = "unit"
	}

	var id int
	err := h.db.Pool.QueryRow(r.Context(), `
		INSERT INTO pantry (user_id, food_id, quantity, units, category, is_frozen)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, req.UserID, req.FoodID, req.Quantity, req.Units, req.Category, req.IsFrozen).Scan(&id)
	if err != nil {
		h.log.Error("failed to create pantry item: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create pantry item")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int{"id": id})
}

// Update handles PATCH /pantry/{id}
func (h *PantryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pantry id")
		return
	}

	var req models.UpdatePantryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Build dynamic UPDATE query
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Quantity != nil {
		setClauses = append(setClauses, "quantity = $"+strconv.Itoa(argIdx))
		args = append(args, *req.Quantity)
		argIdx++
	}
	if req.Units != nil {
		setClauses = append(setClauses, "units = $"+strconv.Itoa(argIdx))
		args = append(args, *req.Units)
		argIdx++
	}
	if req.Category != nil {
		setClauses = append(setClauses, "category = $"+strconv.Itoa(argIdx))
		args = append(args, *req.Category)
		argIdx++
	}
	if req.IsFrozen != nil {
		setClauses = append(setClauses, "is_frozen = $"+strconv.Itoa(argIdx))
		args = append(args, *req.IsFrozen)
		argIdx++
	}

	if len(setClauses) == 0 {
		writeError(w, http.StatusBadRequest, "no fields to update")
		return
	}

	query := "UPDATE pantry SET "
	for i, clause := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += clause
	}
	query += " WHERE id = $" + strconv.Itoa(argIdx)
	args = append(args, id)

	ct, err := h.db.Pool.Exec(r.Context(), query, args...)
	if err != nil {
		h.log.Error("failed to update pantry item %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "failed to update pantry item")
		return
	}

	if ct.RowsAffected() == 0 {
		writeError(w, http.StatusNotFound, "pantry item not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// Delete handles DELETE /pantry/{id}
func (h *PantryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pantry id")
		return
	}

	ct, err := h.db.Pool.Exec(r.Context(), "DELETE FROM pantry WHERE id = $1", id)
	if err != nil {
		h.log.Error("failed to delete pantry item %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "failed to delete pantry item")
		return
	}

	if ct.RowsAffected() == 0 {
		writeError(w, http.StatusNotFound, "pantry item not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// RunMigrations runs the SQL migration to create tables if they don't exist.
func RunMigrations(ctx context.Context, db *database.DB, log logger.Interface) error {
	migration := `
		CREATE TABLE IF NOT EXISTS account (
			account_id SERIAL PRIMARY KEY,
			labels TEXT[] DEFAULT '{}',
			allergens TEXT[] DEFAULT '{}'
		);

		CREATE TABLE IF NOT EXISTS food (
			id SERIAL PRIMARY KEY,
			product_name TEXT NOT NULL,
			environmental_score REAL NOT NULL,
			nutriscore_score REAL NOT NULL,
			labels_en TEXT[] DEFAULT '{}',
			allergens_en TEXT[] DEFAULT '{}',
			traces_en TEXT[] DEFAULT '{}',
			image_url TEXT,
			image_small_url TEXT,
			shelf_life INT,
			category TEXT
		);

		CREATE TABLE IF NOT EXISTS pantry (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES account(account_id) ON DELETE CASCADE,
			food_id INT NOT NULL REFERENCES food(id) ON DELETE CASCADE,
			added_at TIMESTAMP NOT NULL DEFAULT NOW(),
			quantity INT NOT NULL DEFAULT 1,
			units TEXT NOT NULL DEFAULT 'unit',
			category TEXT,
			is_frozen BOOLEAN DEFAULT FALSE
		);
	`

	_, err := db.Pool.Exec(ctx, migration)
	if err != nil {
		return err
	}

	log.Info("Database migrations applied successfully")
	return nil
}
