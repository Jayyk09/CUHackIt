package food

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

type foodHandler struct {
	db  *database.DB
	log *logger.Logger
	cfg *config.Config
	ai  *gemini.Client
}

const (
	defaultLimit = 12
	maxLimit     = 100
)

func (h *foodHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePagination(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	products, err := fetchProducts(r.Context(), h.db, search, limit, offset)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	if h.ai != nil {
		for i := range products {
			if needsEnrichment(products[i]) {
				h.enrichProduct(r.Context(), &products[i])
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		h.log.Error(err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
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

type geminiFoodItem struct {
	FoodName string `json:"food_name"`
}

type geminiCategorization struct {
	FoodName  string `json:"food_name"`
	Category  string `json:"category"`
	ShelfLife int    `json:"shelf_life"`
}

func (h *foodHandler) enrichProduct(ctx context.Context, product *Product) {
	if product == nil || product.ProductName == "" {
		return
	}

	requestPayload, err := json.Marshal([]geminiFoodItem{{FoodName: product.ProductName}})
	if err != nil {
		h.log.Warn("failed to build gemini payload: %v", err)
		return
	}

	response, err := h.ai.Categorize(ctx, string(requestPayload))
	if err != nil {
		h.log.Warn("gemini categorize failed: %v", err)
		return
	}

	cleaned := cleanGeminiResponse(response)
	var categorized []geminiCategorization
	if err := json.Unmarshal([]byte(cleaned), &categorized); err != nil {
		h.log.Warn("failed to parse gemini response: %v", err)
		return
	}

	if len(categorized) == 0 {
		return
	}

	category := strings.TrimSpace(categorized[0].Category)
	if category != "" {
		product.Category = &category
	}

	if categorized[0].ShelfLife > 0 {
		shelfLife := categorized[0].ShelfLife
		product.ShelfLife = &shelfLife
	}

	if product.Category == nil && product.ShelfLife == nil {
		return
	}

	go func(id int64, category *string, shelfLife *int) {
		updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := updateProductMetadata(updateCtx, h.db, id, category, shelfLife); err != nil {
			h.log.Warn("failed to update product metadata: %v", err)
		}
	}(product.ID, product.Category, product.ShelfLife)
}

func cleanGeminiResponse(response string) string {
	cleaned := strings.TrimSpace(response)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	return strings.TrimSpace(cleaned)
}
