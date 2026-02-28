package food

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

// RegisterRoutes registers all food routes
func RegisterRoutes(r *http.ServeMux, db *database.DB, log *logger.Logger) {
	h := NewHandler(db, log)

	r.HandleFunc("GET /food/search", h.Search)
	r.HandleFunc("GET /food/{id}", h.GetProduct)
	r.HandleFunc("PATCH /food/{id}/metadata", h.UpdateMetadata)
}
