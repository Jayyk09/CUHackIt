package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/models"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

type FoodHandler struct {
	db  *database.DB
	log logger.Interface
}

func NewFoodHandler(db *database.DB, log logger.Interface) *FoodHandler {
	return &FoodHandler{db: db, log: log}
}

// GetByID handles GET /food/{id}
func (h *FoodHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid food id")
		return
	}

	var food models.Food
	err = h.db.Pool.QueryRow(r.Context(), `
		SELECT id, product_name, environmental_score, nutriscore_score,
			   labels_en, allergens_en, traces_en,
			   image_url, image_small_url, shelf_life, category
		FROM food
		WHERE id = $1
	`, id).Scan(
		&food.ID,
		&food.ProductName,
		&food.EnvironmentalScore,
		&food.NutriscoreScore,
		&food.LabelsEn,
		&food.AllergensEn,
		&food.TracesEn,
		&food.ImageURL,
		&food.ImageSmallURL,
		&food.ShelfLife,
		&food.Category,
	)
	if err != nil {
		h.log.Error("failed to get food %d: %v", id, err)
		writeError(w, http.StatusNotFound, "food item not found")
		return
	}

	if food.LabelsEn == nil {
		food.LabelsEn = []string{}
	}
	if food.AllergensEn == nil {
		food.AllergensEn = []string{}
	}
	if food.TracesEn == nil {
		food.TracesEn = []string{}
	}

	writeJSON(w, http.StatusOK, food)
}

// Create handles POST /food/{id}
// The {id} in the path is ignored; the food ID is auto-generated.
// If you want to upsert by a specific ID, this can be adjusted.
func (h *FoodHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateFoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ProductName == "" {
		writeError(w, http.StatusBadRequest, "product_name is required")
		return
	}

	if req.LabelsEn == nil {
		req.LabelsEn = []string{}
	}
	if req.AllergensEn == nil {
		req.AllergensEn = []string{}
	}
	if req.TracesEn == nil {
		req.TracesEn = []string{}
	}

	var id int
	err := h.db.Pool.QueryRow(r.Context(), `
		INSERT INTO food (product_name, environmental_score, nutriscore_score,
						  labels_en, allergens_en, traces_en,
						  image_url, image_small_url, shelf_life, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`,
		req.ProductName,
		req.EnvironmentalScore,
		req.NutriscoreScore,
		req.LabelsEn,
		req.AllergensEn,
		req.TracesEn,
		req.ImageURL,
		req.ImageSmallURL,
		req.ShelfLife,
		req.Category,
	).Scan(&id)
	if err != nil {
		h.log.Error("failed to create food: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create food item")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int{"id": id})
}
