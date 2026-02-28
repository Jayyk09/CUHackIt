package food

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

// Handler handles HTTP requests for food products
type Handler struct {
	db  *database.DB
	log *logger.Logger
	ai  *gemini.Client
}

const (
	defaultLimit = 12
	maxLimit     = 100
)

// NewHandler creates a new food handler
func NewHandler(db *database.DB, ai *gemini.Client, log *logger.Logger) *Handler {
	return &Handler{db: db, ai: ai, log: log}
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

// List handles GET /food/search?q=...&limit=...&offset=...
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePagination(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("q"))
	products, err := fetchProducts(r.Context(), h.db, search, limit, offset)
	if err != nil {
		h.log.Error("Failed to fetch products: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	if h.ai != nil {
		for i := range products {
			if needsEnrichment(products[i]) {
				h.enrichProduct(r.Context(), &products[i])
			}
		}
	}

	h.writeJSON(w, http.StatusOK, products)
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

	rows, err := h.db.Pool.Query(r.Context(),
		`SELECT id, product_name, norm_environmental_score, nutriscore_score,
		 labels_en, allergens_en, traces_en, image_url, image_small_url, shelf_life, category
		 FROM foods WHERE id = $1`, id)
	if err != nil {
		h.log.Error("Failed to get product: %v", err)
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

func parsePagination(r *http.Request) (int, int, error) {
	query := r.URL.Query()
	limit := defaultLimit
	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, 0, err
		}
		limit = parsed
	}
	if limit < 0 {
		limit = 0
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := 0
	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, 0, err
		}
		offset = parsed
	}
	if offset < 0 {
		offset = 0
	}

	return limit, offset, nil
}

func needsEnrichment(product Product) bool {
	categoryMissing := product.Category == nil || strings.TrimSpace(*product.Category) == ""
	shelfLifeMissing := product.ShelfLife == nil || *product.ShelfLife <= 0
	return categoryMissing || shelfLifeMissing
}

func (h *Handler) enrichProduct(ctx context.Context, product *Product) {
	if product == nil || product.ProductName == "" {
		return
	}

	categories, err := h.ai.CategorizeFood(ctx, []string{product.ProductName})
	if err != nil {
		h.log.Error("gemini categorize failed: %v", err)
		return
	}

	category, ok := categories[product.ProductName]
	if !ok || strings.TrimSpace(category) == "" {
		return
	}

	category = strings.TrimSpace(category)
	product.Category = &category

	go func(id int64, category *string) {
		updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := updateProductMetadata(updateCtx, h.db, id, category, nil); err != nil {
			h.log.Error("failed to update product metadata: %v", err)
		}
	}(product.ID, product.Category)
}
