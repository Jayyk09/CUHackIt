package food

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

// Handler handles HTTP requests for food products
type Handler struct {
	log *logger.Logger
}

// NewHandler creates a new food handler
func NewHandler(log *logger.Logger) *Handler {
	return &Handler{log: log}
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

// SearchResponse is the response for search queries
type SearchResponse struct {
	Query    string        `json:"query"`
	Count    int           `json:"count"`
	Products []FoodProduct `json:"products"`
}

// Search handles GET /food/search?q=...&limit=...
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		h.writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	if len(query) < 3 {
		h.writeError(w, http.StatusBadRequest, "query must be at least 3 characters")
		return
	}

	// Parse limit (default 20, max 50)
	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= 50 {
				limit = parsedLimit
			}
		}
	}

	products := SearchProducts(query, limit)

	h.writeJSON(w, http.StatusOK, SearchResponse{
		Query:    query,
		Count:    len(products),
		Products: products,
	})
}

// GetProduct handles GET /food/{id}
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.writeError(w, http.StatusBadRequest, "product id is required")
		return
	}

	product := GetProductByID(id)
	if product == nil {
		h.writeError(w, http.StatusNotFound, "product not found")
		return
	}

	h.writeJSON(w, http.StatusOK, product)
}

// ListByCategory handles GET /food/category/{category}
func (h *Handler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	categoryStr := r.PathValue("category")
	if categoryStr == "" {
		h.writeError(w, http.StatusBadRequest, "category is required")
		return
	}

	category := FoodCategory(categoryStr)

	// Parse limit (default 20, max 50)
	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= 50 {
				limit = parsedLimit
			}
		}
	}

	products := GetProductsByCategory(category, limit)

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"category": category,
		"count":    len(products),
		"products": products,
	})
}

// ListCategories handles GET /food/categories
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories := GetAllCategories()

	// Convert to a slice for better JSON output
	type CategoryInfo struct {
		Name  FoodCategory `json:"name"`
		Count int          `json:"count"`
	}

	var result []CategoryInfo
	for cat, count := range categories {
		result = append(result, CategoryInfo{Name: cat, Count: count})
	}

	h.writeJSON(w, http.StatusOK, result)
}
