package food

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

// RegisterRoutes registers all food routes
func RegisterRoutes(r *http.ServeMux, log *logger.Logger) {
	h := NewHandler(log)

	// Search
	r.HandleFunc("GET /food/search", h.Search)

	// Get product by ID
	r.HandleFunc("GET /food/{id}", h.GetProduct)

	// List by category
	r.HandleFunc("GET /food/category/{category}", h.ListByCategory)

	// List all categories
	r.HandleFunc("GET /food/categories", h.ListCategories)
}
