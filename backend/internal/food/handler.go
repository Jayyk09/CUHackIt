package food

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

// Handler handles HTTP requests for food products
type Handler struct {
	db  *database.DB
	log *logger.Logger
}

// NewHandler creates a new food handler
func NewHandler(db *database.DB, log *logger.Logger) *Handler {
	return &Handler{db: db, log: log}
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

// Search handles GET /food/search?q=...&limit=...&offset=...
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	products, err := fetchProducts(r.Context(), h.db, query, limit, offset)
	if err != nil {
		h.log.Error("Failed to search products: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to search products")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"query":    query,
		"count":    len(products),
		"products": products,
	})
}

// GetProduct handles GET /food/{id}
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		h.writeError(w, http.StatusBadRequest, "product id is required")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	products, err := fetchProducts(r.Context(), h.db, "", 1, 0)
	_ = products
	_ = id
	// Look up by ID directly
	rows, err2 := h.db.Pool.Query(r.Context(),
		`SELECT id, product_name, norm_environmental_score, nutriscore_score,
		 labels_en, allergens_en, traces_en, image_url, image_small_url, shelf_life, category
		 FROM foods WHERE id = $1`, id)
	if err2 != nil {
		h.log.Error("Failed to get product: %v", err2)
		h.writeError(w, http.StatusInternalServerError, "failed to get product")
		return
	}
	defer rows.Close()

	if !rows.Next() {
		h.writeError(w, http.StatusNotFound, "product not found")
		return
	}

	var product Product
	if err := rows.Scan(
		&product.ID, &product.ProductName, &product.NormEnvironmentalScore,
		&product.NutriscoreScore, &product.LabelsEn, &product.AllergensEn,
		&product.TracesEn, &product.ImageURL, &product.ImageSmallURL,
		&product.ShelfLife, &product.Category,
	); err != nil {
		h.log.Error("Failed to scan product: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to get product")
		return
	}

	h.writeJSON(w, http.StatusOK, product)
}

// UpdateMetadata handles PATCH /food/{id}/metadata
func (h *Handler) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		h.writeError(w, http.StatusBadRequest, "product id is required")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var body struct {
		Category  *string `json:"category"`
		ShelfLife *int    `json:"shelf_life"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := updateProductMetadata(r.Context(), h.db, id, body.Category, body.ShelfLife); err != nil {
		h.log.Error("Failed to update product metadata: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to update product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
