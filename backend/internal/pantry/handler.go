package pantry

import (
	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

type pantryHandler struct {
	db  *database.DB
	log *logger.Logger
	cfg *config.Config
}
