package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	ErrNoAPIKey       = errors.New("gemini API key not configured")
	ErrGenerationFail = errors.New("failed to generate content")
	ErrInvalidResponse = errors.New("invalid response from Gemini")
)

// Client is a wrapper around the Gemini API client
type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
	log    *logger.Logger
}

// Recipe represents a generated recipe
type Recipe struct {
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	Cuisine          string       `json:"cuisine,omitempty"`
	PrepTimeMinutes  int          `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes  int          `json:"cook_time_minutes,omitempty"`
	Servings         int          `json:"servings,omitempty"`
	Difficulty       string       `json:"difficulty,omitempty"`
	Ingredients      []Ingredient `json:"ingredients"`
	Instructions     []string     `json:"instructions"`
	MissingItems     []Ingredient `json:"missing_items,omitempty"`
	CaloriesPerServing float64    `json:"calories_per_serving,omitempty"`
	ProteinG         float64      `json:"protein_g,omitempty"`
	CarbsG           float64      `json:"carbs_g,omitempty"`
	FatG             float64      `json:"fat_g,omitempty"`
	Tags             []string     `json:"tags,omitempty"`
}

// Ingredient represents an ingredient in a recipe
type Ingredient struct {
	Name       string  `json:"name"`
	Amount     string  `json:"amount"`
	Unit       string  `json:"unit,omitempty"`
	FromPantry bool    `json:"from_pantry"`
}

// PantryItem represents an item in the user's pantry for recipe generation
type PantryItem struct {
	Name           string  `json:"name"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	Category       string  `json:"category"`
	IsExpiringSoon bool    `json:"is_expiring_soon"`
	DaysUntilExpiry int    `json:"days_until_expiry,omitempty"`
}

// UserPreferences contains user dietary preferences for recipe generation
type UserPreferences struct {
	Allergens          []string `json:"allergens,omitempty"`
	DietaryPreferences []string `json:"dietary_preferences,omitempty"`
	NutritionalGoals   []string `json:"nutritional_goals,omitempty"`
	CookingSkill       string   `json:"cooking_skill,omitempty"`
	CuisinePreferences []string `json:"cuisine_preferences,omitempty"`
}

// GenerateRecipesRequest is the input for recipe generation
type GenerateRecipesRequest struct {
	PantryItems    []PantryItem    `json:"pantry_items"`
	Preferences    UserPreferences `json:"preferences"`
	RecipeCount    int             `json:"recipe_count"`
	PantryOnly     bool            `json:"pantry_only"`      // If true, only use pantry items
	MaxMissingItems int            `json:"max_missing_items"` // Max additional items (0 for pantry-only)
}

// NewClient creates a new Gemini client
func NewClient(ctx context.Context, apiKey string, modelName string, log *logger.Logger) (*Client, error) {
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	if modelName == "" {
		modelName = "gemini-1.5-flash"
	}

	model := client.GenerativeModel(modelName)
	
	// Configure the model for JSON output
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.ResponseMIMEType = "application/json"

	return &Client{
		client: client,
		model:  model,
		log:    log,
	}, nil
}

// Close closes the Gemini client
func (c *Client) Close() error {
	return c.client.Close()
}

// GenerateRecipes generates recipes based on pantry items and preferences
func (c *Client) GenerateRecipes(ctx context.Context, req GenerateRecipesRequest) ([]Recipe, error) {
	prompt := buildRecipePrompt(req)

	c.log.Debug("Generating recipes with prompt length: %d", len(prompt))

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		c.log.Error("Gemini generation error: %v", err)
		return nil, fmt.Errorf("%w: %v", ErrGenerationFail, err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, ErrInvalidResponse
	}

	// Extract text from response
	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, ErrInvalidResponse
	}

	// Parse JSON response
	var result struct {
		Recipes []Recipe `json:"recipes"`
	}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		// Try parsing as array directly
		var recipes []Recipe
		if err2 := json.Unmarshal([]byte(text), &recipes); err2 != nil {
			c.log.Error("Failed to parse recipes JSON: %v", err)
			return nil, fmt.Errorf("%w: %v", ErrInvalidResponse, err)
		}
		return recipes, nil
	}

	return result.Recipes, nil
}

// CategorizeFood categorizes food items using AI
func (c *Client) CategorizeFood(ctx context.Context, items []string) (map[string]string, error) {
	input, _ := json.Marshal(items)
	prompt := CategorizerPrompt + "\n\n" + string(input)

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGenerationFail, err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, ErrInvalidResponse
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, ErrInvalidResponse
	}

	var result []struct {
		FoodName  string `json:"food_name"`
		Category  string `json:"category"`
		ShelfLife int    `json:"shelf_life"`
	}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidResponse, err)
	}

	categories := make(map[string]string)
	for _, item := range result {
		categories[item.FoodName] = item.Category
	}

	return categories, nil
}

// buildRecipePrompt constructs the prompt for recipe generation
func buildRecipePrompt(req GenerateRecipesRequest) string {
	pantryJSON, _ := json.MarshalIndent(req.PantryItems, "", "  ")
	prefsJSON, _ := json.MarshalIndent(req.Preferences, "", "  ")

	recipeType := "flexible"
	maxMissing := req.MaxMissingItems
	if req.PantryOnly || maxMissing == 0 {
		recipeType = "pantry-only"
		maxMissing = 0
	}

	prompt := fmt.Sprintf(`You are a professional chef and recipe recommendation engine. Generate %d recipes based on the user's pantry items and preferences.

## Recipe Type: %s
%s

## User's Pantry Items:
%s

## User's Preferences:
%s

## Instructions:
1. Prioritize items that are expiring soon (is_expiring_soon: true)
2. Create recipes that match the user's cooking skill level
3. Respect ALL allergens - never include any ingredient the user is allergic to
4. Consider dietary preferences and nutritional goals
5. If cuisine preferences are specified, favor those cuisines
%s

## Response Format:
Return a JSON object with a "recipes" array containing exactly %d recipes. Each recipe must have:
- title: Creative, appetizing name
- description: 1-2 sentence description
- cuisine: The cuisine type (Italian, Asian, Mexican, etc.)
- prep_time_minutes: Realistic prep time
- cook_time_minutes: Realistic cook time  
- servings: Number of servings (2-4 typical)
- difficulty: "easy", "medium", or "hard"
- ingredients: Array of {name, amount, unit, from_pantry}
- instructions: Array of step-by-step instructions (clear, numbered)
- missing_items: Array of {name, amount, unit} for items not in pantry (empty for pantry-only)
- calories_per_serving: Estimated calories
- protein_g, carbs_g, fat_g: Estimated macros
- tags: Array of relevant tags (e.g., "quick", "healthy", "comfort-food")

Generate recipes now:`,
		req.RecipeCount,
		recipeType,
		func() string {
			if recipeType == "pantry-only" {
				return "Use ONLY ingredients from the pantry. Do not suggest any additional ingredients."
			}
			return fmt.Sprintf("You may suggest up to %d additional ingredients not in the pantry.", maxMissing)
		}(),
		string(pantryJSON),
		string(prefsJSON),
		func() string {
			if len(req.Preferences.Allergens) > 0 {
				return fmt.Sprintf("\n⚠️ CRITICAL: User has the following allergens. NEVER include these or any derivatives: %v", req.Preferences.Allergens)
			}
			return ""
		}(),
		req.RecipeCount,
	)

	return prompt
}

// GenerateText is a general-purpose text generation method
func (c *Client) GenerateText(ctx context.Context, prompt string) (string, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrGenerationFail, err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", ErrInvalidResponse
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", ErrInvalidResponse
	}

	return string(text), nil
}
