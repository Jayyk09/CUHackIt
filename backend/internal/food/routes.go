package food

import (
	"net/http"

	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

// RegisterRoutes registers all food routes
func RegisterRoutes(r *http.ServeMux, db *database.DB, ai *gemini.Client, log *logger.Logger) {
	h := NewHandler(db, ai, log)

	r.HandleFunc("GET /food/search", h.List)
	r.HandleFunc("GET /food/{id}", h.GetProduct)
	r.HandleFunc("PATCH /food/{id}/metadata", h.UpdateMetadata)
}
