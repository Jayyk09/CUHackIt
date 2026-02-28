package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

// SpoilingAgent generates recipes that prioritize ingredients close to expiring.
// It uses a custom prompt that tells the AI to build meals around the most
// urgent items first.
type SpoilingAgent struct {
	client *gemini.Client
	log    *logger.Logger
}

// NewSpoilingAgent creates a new SpoilingAgent
func NewSpoilingAgent(client *gemini.Client, log *logger.Logger) *SpoilingAgent {
	return &SpoilingAgent{
		client: client,
		log:    log,
	}
}

// Name returns the agent's identifier
func (a *SpoilingAgent) Name() string {
	return "spoiling"
}

// GenerateRecipes generates recipes prioritizing expiring ingredients
func (a *SpoilingAgent) GenerateRecipes(ctx context.Context, req RecipeRequest) (*RecipeResponse, error) {
	// Sort pantry items by days until expiry (ascending) so the prompt
	// presents the most urgent items first.
	sorted := make([]PantryItem, len(req.PantryItems))
	copy(sorted, req.PantryItems)
	sort.Slice(sorted, func(i, j int) bool {
		di, dj := daysUntilExpiry(sorted[i]), daysUntilExpiry(sorted[j])
		return di < dj
	})

	// Build the list of urgent items (expiring within 3 days)
	var urgentNames []string
	for _, item := range sorted {
		if daysUntilExpiry(item) <= 3 {
			urgentNames = append(urgentNames, item.Name)
		}
	}

	a.log.Info("SpoilingAgent: Generating %d recipes from %d pantry items (%d expiring soon)",
		req.RecipeCount, len(req.PantryItems), len(urgentNames))

	geminiItems := convertToGeminiPantryItems(sorted)
	prefs := convertToGeminiPreferences(req)

	pantryJSON, _ := json.MarshalIndent(geminiItems, "", "  ")
	prefsJSON, _ := json.MarshalIndent(prefs, "", "  ")
	urgentJSON, _ := json.Marshal(urgentNames)

	prompt := buildSpoilingPrompt(string(pantryJSON), string(prefsJSON), string(urgentJSON), req.RecipeCount, prefs)

	recipes, err := a.client.GenerateRecipes(ctx, gemini.GenerateRecipesRequest{
		PantryItems:     geminiItems,
		Preferences:     prefs,
		RecipeCount:     req.RecipeCount,
		PantryOnly:      true,
		MaxMissingItems: 0,
		UserPrompt:      prompt,
	})
	if err != nil {
		a.log.Error("SpoilingAgent: Failed to generate recipes: %v", err)
		return nil, err
	}

	a.log.Info("SpoilingAgent: Generated %d recipes", len(recipes))

	return &RecipeResponse{
		Recipes:     convertFromGeminiRecipes(recipes, "spoiling"),
		AgentUsed:   a.Name(),
		GeneratedAt: time.Now(),
	}, nil
}

// daysUntilExpiry returns the days until a pantry item expires.
// Items without an expiration date return a large number.
func daysUntilExpiry(item PantryItem) int {
	if item.ExpirationDate == nil {
		return 999
	}
	d := int(time.Until(*item.ExpirationDate).Hours() / 24)
	if d < 0 {
		return 0
	}
	return d
}

// buildSpoilingPrompt constructs a prompt that heavily prioritizes expiring items
func buildSpoilingPrompt(pantryJSON, prefsJSON, urgentJSON string, count int, prefs gemini.UserPreferences) string {
	allergenWarning := ""
	if len(prefs.Allergens) > 0 {
		allergenWarning = fmt.Sprintf("\nCRITICAL: User has the following allergens. NEVER include these or any derivatives: %v", prefs.Allergens)
	}

	return fmt.Sprintf(`You are a professional chef specializing in REDUCING FOOD WASTE.
Your #1 goal is to create delicious recipes that USE UP ingredients that are about to expire.

## URGENT - These ingredients are expiring within 3 days and MUST be used:
%s

Every recipe you generate MUST prominently feature at least one of the urgent ingredients above.
Build each recipe AROUND these expiring items. They should be the star of the dish, not a garnish.

## Full Pantry (sorted by urgency — items expiring soonest listed first):
%s

## User Preferences:
%s
%s

## Rules:
1. EVERY recipe MUST use at least one of the urgent expiring ingredients as a PRIMARY ingredient
2. Use ONLY ingredients from the pantry — do not suggest any additional purchases
3. Prioritize recipes that use MULTIPLE expiring items together
4. Respect the user's allergens, dietary preferences, and cooking skill level
5. Make the recipes practical and quick — the user needs to cook these soon
6. Mark all ingredients with from_pantry: true

## Response Format:
Return a JSON object with a "recipes" array containing exactly %d recipes. Each recipe must have:
- title: Creative, appetizing name
- description: 1-2 sentence description mentioning which expiring items it uses up
- cuisine: The cuisine type
- prep_time_minutes: Realistic prep time
- cook_time_minutes: Realistic cook time
- servings: Number of servings (2-4 typical)
- difficulty: "easy", "medium", or "hard"
- ingredients: Array of {name, amount, unit, from_pantry}
- instructions: Array of step-by-step instructions
- missing_items: [] (always empty — pantry only)
- calories_per_serving: Estimated calories
- protein_g, carbs_g, fat_g: Estimated macros
- tags: Array of relevant tags — always include "use-it-up"

Generate recipes now:`,
		urgentJSON,
		pantryJSON,
		prefsJSON,
		allergenWarning,
		count,
	)
}
