package agents

import (
	"context"
	"time"

	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

// PantryOnlyAgent generates recipes using ONLY items from the user's pantry
// No additional ingredients will be suggested
type PantryOnlyAgent struct {
	client *gemini.Client
	log    *logger.Logger
}

// NewPantryOnlyAgent creates a new PantryOnlyAgent
func NewPantryOnlyAgent(client *gemini.Client, log *logger.Logger) *PantryOnlyAgent {
	return &PantryOnlyAgent{
		client: client,
		log:    log,
	}
}

// Name returns the agent's identifier
func (a *PantryOnlyAgent) Name() string {
	return "pantry_only"
}

// GenerateRecipes generates recipes using only pantry items
func (a *PantryOnlyAgent) GenerateRecipes(ctx context.Context, req RecipeRequest) (*RecipeResponse, error) {
	a.log.Info("PantryOnlyAgent: Generating %d recipes from %d pantry items", req.RecipeCount, len(req.PantryItems))

	// Convert to Gemini request format
	geminiReq := gemini.GenerateRecipesRequest{
		PantryItems:     convertToGeminiPantryItems(req.PantryItems),
		Preferences:     convertToGeminiPreferences(req),
		RecipeCount:     req.RecipeCount,
		PantryOnly:      true,
		MaxMissingItems: 0, // Strictly pantry-only
		UserPrompt:      req.UserPrompt,
	}

	// Generate recipes
	recipes, err := a.client.GenerateRecipes(ctx, geminiReq)
	if err != nil {
		a.log.Error("PantryOnlyAgent: Failed to generate recipes: %v", err)
		return nil, err
	}

	a.log.Info("PantryOnlyAgent: Generated %d recipes", len(recipes))

	return &RecipeResponse{
		Recipes:     convertFromGeminiRecipes(recipes, "pantry_only"),
		AgentUsed:   a.Name(),
		GeneratedAt: time.Now(),
	}, nil
}
