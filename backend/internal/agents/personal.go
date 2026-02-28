package agents

import (
	"context"
	"time"

	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

// PersonalRecipeAgent generates recipes purely from the user's stored profile
// (allergens, dietary prefs, nutritional goals, cooking skill, cuisines).
// No pantry is required — all ingredients are "to buy".
type PersonalRecipeAgent struct {
	client *gemini.Client
	log    *logger.Logger
}

// NewPersonalRecipeAgent creates a new PersonalRecipeAgent.
func NewPersonalRecipeAgent(client *gemini.Client, log *logger.Logger) *PersonalRecipeAgent {
	return &PersonalRecipeAgent{
		client: client,
		log:    log,
	}
}

// Name returns the agent's identifier.
func (a *PersonalRecipeAgent) Name() string {
	return "personal"
}

// GenerateRecipes generates profile-first recipes with allergens as hard blocks.
func (a *PersonalRecipeAgent) GenerateRecipes(ctx context.Context, req RecipeRequest) (*RecipeResponse, error) {
	a.log.Info("PersonalRecipeAgent: Generating %d recipes from user profile", req.RecipeCount)

	geminiReq := gemini.GenerateRecipesRequest{
		PantryItems:     []gemini.PantryItem{}, // no pantry
		Preferences:     convertToGeminiPreferences(req),
		RecipeCount:     req.RecipeCount,
		PantryOnly:      false,
		MaxMissingItems: 10, // everything is "to buy"
		UserPrompt:      req.UserPrompt,
	}

	recipes, err := a.client.GeneratePersonalRecipes(ctx, geminiReq)
	if err != nil {
		a.log.Error("PersonalRecipeAgent: Failed to generate recipes: %v", err)
		return nil, err
	}

	a.log.Info("PersonalRecipeAgent: Generated %d recipes", len(recipes))

	return &RecipeResponse{
		Recipes:     convertFromGeminiRecipes(recipes, "personal"),
		AgentUsed:   a.Name(),
		GeneratedAt: time.Now(),
	}, nil
}
