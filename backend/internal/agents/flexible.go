package agents

import (
	"context"
	"time"

	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

const DefaultMaxMissingItems = 3

// FlexibleRecipeAgent generates recipes using pantry items plus a limited number of additional ingredients
type FlexibleRecipeAgent struct {
	client          *gemini.Client
	log             *logger.Logger
	maxMissingItems int
}

// NewFlexibleRecipeAgent creates a new FlexibleRecipeAgent
func NewFlexibleRecipeAgent(client *gemini.Client, log *logger.Logger) *FlexibleRecipeAgent {
	return &FlexibleRecipeAgent{
		client:          client,
		log:             log,
		maxMissingItems: DefaultMaxMissingItems,
	}
}

// WithMaxMissingItems sets the maximum number of additional ingredients
func (a *FlexibleRecipeAgent) WithMaxMissingItems(max int) *FlexibleRecipeAgent {
	a.maxMissingItems = max
	return a
}

// Name returns the agent's identifier
func (a *FlexibleRecipeAgent) Name() string {
	return "flexible"
}

// GenerateRecipes generates recipes using pantry items plus up to N additional ingredients
func (a *FlexibleRecipeAgent) GenerateRecipes(ctx context.Context, req RecipeRequest) (*RecipeResponse, error) {
	a.log.Info("FlexibleRecipeAgent: Generating %d recipes from %d pantry items (max %d missing allowed)",
		req.RecipeCount, len(req.PantryItems), a.maxMissingItems)

	// Convert to Gemini request format
	geminiReq := gemini.GenerateRecipesRequest{
		PantryItems:     convertToGeminiPantryItems(req.PantryItems),
		Preferences:     convertToGeminiPreferences(req),
		RecipeCount:     req.RecipeCount,
		PantryOnly:      false,
		MaxMissingItems: a.maxMissingItems,
		UserPrompt:      req.UserPrompt,
	}

	// Generate recipes
	recipes, err := a.client.GenerateRecipes(ctx, geminiReq)
	if err != nil {
		a.log.Error("FlexibleRecipeAgent: Failed to generate recipes: %v", err)
		return nil, err
	}

	a.log.Info("FlexibleRecipeAgent: Generated %d recipes", len(recipes))

	return &RecipeResponse{
		Recipes:     convertFromGeminiRecipes(recipes, "flexible"),
		AgentUsed:   a.Name(),
		GeneratedAt: time.Now(),
	}, nil
}
