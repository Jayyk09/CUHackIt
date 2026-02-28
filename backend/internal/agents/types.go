package agents

import (
	"context"
	"time"

	"github.com/Jayyk09/CUHackIt/services/gemini"
)

// Agent is the interface that all recipe generation agents must implement
type Agent interface {
	// Name returns the agent's identifier
	Name() string

	// GenerateRecipes generates recipes based on the request
	GenerateRecipes(ctx context.Context, req RecipeRequest) (*RecipeResponse, error)
}

// RecipeRequest contains all the information needed to generate recipes
type RecipeRequest struct {
	// User's pantry items
	PantryItems []PantryItem `json:"pantry_items"`

	// User's preferences and restrictions
	Allergens          []string `json:"allergens"`
	DietaryPreferences []string `json:"dietary_preferences"`
	NutritionalGoals   []string `json:"nutritional_goals"`
	CookingSkill       string   `json:"cooking_skill"`
	CuisinePreferences []string `json:"cuisine_preferences"`

	// Generation options
	RecipeCount int `json:"recipe_count"` // How many recipes to generate (1-3)
}

// PantryItem represents an item in the user's pantry
type PantryItem struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Quantity       float64   `json:"quantity"`
	Unit           string    `json:"unit"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
	IsExpiringSoon bool      `json:"is_expiring_soon"`
	IsExpired      bool      `json:"is_expired"`
}

// RecipeResponse contains the generated recipes and metadata
type RecipeResponse struct {
	Recipes    []Recipe `json:"recipes"`
	AgentUsed  string   `json:"agent_used"`
	GeneratedAt time.Time `json:"generated_at"`
}

// Recipe represents a generated recipe
type Recipe struct {
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	Cuisine            string       `json:"cuisine,omitempty"`
	PrepTimeMinutes    int          `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes    int          `json:"cook_time_minutes,omitempty"`
	TotalTimeMinutes   int          `json:"total_time_minutes,omitempty"`
	Servings           int          `json:"servings,omitempty"`
	Difficulty         string       `json:"difficulty,omitempty"`
	Ingredients        []Ingredient `json:"ingredients"`
	Instructions       []string     `json:"instructions"`
	MissingIngredients []Ingredient `json:"missing_ingredients,omitempty"`
	CaloriesPerServing float64      `json:"calories_per_serving,omitempty"`
	ProteinG           float64      `json:"protein_g,omitempty"`
	CarbsG             float64      `json:"carbs_g,omitempty"`
	FatG               float64      `json:"fat_g,omitempty"`
	Tags               []string     `json:"tags,omitempty"`
	Source             string       `json:"source"` // "pantry_only" or "flexible"
}

// Ingredient represents an ingredient in a recipe
type Ingredient struct {
	Name       string `json:"name"`
	Amount     string `json:"amount"`
	Unit       string `json:"unit,omitempty"`
	FromPantry bool   `json:"from_pantry"`
}

// convertToGeminiPantryItems converts agent PantryItems to Gemini PantryItems
func convertToGeminiPantryItems(items []PantryItem) []gemini.PantryItem {
	result := make([]gemini.PantryItem, len(items))
	for i, item := range items {
		daysUntilExpiry := 999 // Default for non-expiring items
		if item.ExpirationDate != nil {
			daysUntilExpiry = int(time.Until(*item.ExpirationDate).Hours() / 24)
			if daysUntilExpiry < 0 {
				daysUntilExpiry = 0
			}
		}
		result[i] = gemini.PantryItem{
			Name:           item.Name,
			Quantity:       item.Quantity,
			Unit:           item.Unit,
			Category:       item.Category,
			IsExpiringSoon: item.IsExpiringSoon,
			DaysUntilExpiry: daysUntilExpiry,
		}
	}
	return result
}

// convertToGeminiPreferences converts agent preferences to Gemini UserPreferences
func convertToGeminiPreferences(req RecipeRequest) gemini.UserPreferences {
	return gemini.UserPreferences{
		Allergens:          req.Allergens,
		DietaryPreferences: req.DietaryPreferences,
		NutritionalGoals:   req.NutritionalGoals,
		CookingSkill:       req.CookingSkill,
		CuisinePreferences: req.CuisinePreferences,
	}
}

// convertFromGeminiRecipes converts Gemini Recipes to agent Recipes
func convertFromGeminiRecipes(recipes []gemini.Recipe, source string) []Recipe {
	result := make([]Recipe, len(recipes))
	for i, r := range recipes {
		ingredients := make([]Ingredient, len(r.Ingredients))
		for j, ing := range r.Ingredients {
			ingredients[j] = Ingredient{
				Name:       ing.Name,
				Amount:     ing.Amount,
				Unit:       ing.Unit,
				FromPantry: ing.FromPantry,
			}
		}

		missingIngredients := make([]Ingredient, len(r.MissingItems))
		for j, ing := range r.MissingItems {
			missingIngredients[j] = Ingredient{
				Name:       ing.Name,
				Amount:     ing.Amount,
				Unit:       ing.Unit,
				FromPantry: false,
			}
		}

		result[i] = Recipe{
			Title:              r.Title,
			Description:        r.Description,
			Cuisine:            r.Cuisine,
			PrepTimeMinutes:    r.PrepTimeMinutes,
			CookTimeMinutes:    r.CookTimeMinutes,
			TotalTimeMinutes:   r.PrepTimeMinutes + r.CookTimeMinutes,
			Servings:           r.Servings,
			Difficulty:         r.Difficulty,
			Ingredients:        ingredients,
			Instructions:       r.Instructions,
			MissingIngredients: missingIngredients,
			CaloriesPerServing: r.CaloriesPerServing,
			ProteinG:           r.ProteinG,
			CarbsG:             r.CarbsG,
			FatG:               r.FatG,
			Tags:               r.Tags,
			Source:             source,
		}
	}
	return result
}
