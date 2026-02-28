package agents

import (
	"context"
	"errors"
	"time"

	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
)

var (
	ErrNoRecipesGenerated = errors.New("no recipes could be generated")
	ErrAllRecipesFiltered = errors.New("all recipes were filtered due to allergens")
	ErrInvalidRequest     = errors.New("invalid recipe request")
)

// OrchestratorMode determines which agents to use
type OrchestratorMode string

const (
	ModePantryOnly OrchestratorMode = "pantry_only"
	ModeFlexible   OrchestratorMode = "flexible"
	ModeBoth       OrchestratorMode = "both"
)

// Orchestrator coordinates multiple agents to generate recipes
type Orchestrator struct {
	pantryAgent   *PantryOnlyAgent
	flexibleAgent *FlexibleRecipeAgent
	allergenFilter *AllergenFilter
	log           *logger.Logger
}

// NewOrchestrator creates a new Orchestrator with all agents
func NewOrchestrator(geminiClient *gemini.Client, log *logger.Logger) *Orchestrator {
	return &Orchestrator{
		pantryAgent:    NewPantryOnlyAgent(geminiClient, log),
		flexibleAgent:  NewFlexibleRecipeAgent(geminiClient, log),
		allergenFilter: NewAllergenFilter(log),
		log:            log,
	}
}

// GenerateRequest contains the request for the orchestrator
type GenerateRequest struct {
	RecipeRequest
	Mode OrchestratorMode `json:"mode"`
}

// GenerateResult contains the combined results from all agents
type GenerateResult struct {
	PantryOnlyRecipes []Recipe  `json:"pantry_only_recipes,omitempty"`
	FlexibleRecipes   []Recipe  `json:"flexible_recipes,omitempty"`
	AllRecipes        []Recipe  `json:"all_recipes"`
	GeneratedAt       time.Time `json:"generated_at"`
	TotalCount        int       `json:"total_count"`
	FilteredCount     int       `json:"filtered_count"` // How many were removed due to allergens
}

// Generate orchestrates recipe generation across agents
func (o *Orchestrator) Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error) {
	if len(req.PantryItems) == 0 {
		return nil, ErrInvalidRequest
	}

	if req.RecipeCount <= 0 {
		req.RecipeCount = 2 // Default to 2 recipes
	}
	if req.RecipeCount > 3 {
		req.RecipeCount = 3 // Max 3 recipes per agent
	}

	o.log.Info("Orchestrator: Starting recipe generation (mode=%s, pantry_items=%d, recipe_count=%d)",
		req.Mode, len(req.PantryItems), req.RecipeCount)

	result := &GenerateResult{
		GeneratedAt: time.Now(),
	}

	var allRecipes []Recipe
	var totalGenerated int

	switch req.Mode {
	case ModePantryOnly:
		recipes, err := o.generatePantryOnly(ctx, req.RecipeRequest)
		if err != nil {
			return nil, err
		}
		result.PantryOnlyRecipes = recipes
		allRecipes = append(allRecipes, recipes...)
		totalGenerated = len(recipes)

	case ModeFlexible:
		recipes, err := o.generateFlexible(ctx, req.RecipeRequest)
		if err != nil {
			return nil, err
		}
		result.FlexibleRecipes = recipes
		allRecipes = append(allRecipes, recipes...)
		totalGenerated = len(recipes)

	case ModeBoth:
		// Generate from both agents concurrently
		pantryRecipes, flexibleRecipes, err := o.generateBoth(ctx, req.RecipeRequest)
		if err != nil {
			return nil, err
		}
		result.PantryOnlyRecipes = pantryRecipes
		result.FlexibleRecipes = flexibleRecipes
		allRecipes = append(allRecipes, pantryRecipes...)
		allRecipes = append(allRecipes, flexibleRecipes...)
		totalGenerated = len(pantryRecipes) + len(flexibleRecipes)

	default:
		// Default to pantry-only
		recipes, err := o.generatePantryOnly(ctx, req.RecipeRequest)
		if err != nil {
			return nil, err
		}
		result.PantryOnlyRecipes = recipes
		allRecipes = append(allRecipes, recipes...)
		totalGenerated = len(recipes)
	}

	// Apply allergen filter
	if len(req.Allergens) > 0 {
		filteredRecipes := o.allergenFilter.FilterRecipes(allRecipes, req.Allergens)
		result.FilteredCount = totalGenerated - len(filteredRecipes)
		allRecipes = filteredRecipes

		// Also filter the categorized lists
		if result.PantryOnlyRecipes != nil {
			result.PantryOnlyRecipes = o.allergenFilter.FilterRecipes(result.PantryOnlyRecipes, req.Allergens)
		}
		if result.FlexibleRecipes != nil {
			result.FlexibleRecipes = o.allergenFilter.FilterRecipes(result.FlexibleRecipes, req.Allergens)
		}
	}

	result.AllRecipes = allRecipes
	result.TotalCount = len(allRecipes)

	if result.TotalCount == 0 {
		if result.FilteredCount > 0 {
			return result, ErrAllRecipesFiltered
		}
		return result, ErrNoRecipesGenerated
	}

	o.log.Info("Orchestrator: Generated %d recipes (%d filtered)", result.TotalCount, result.FilteredCount)

	return result, nil
}

// generatePantryOnly generates recipes using only pantry items
func (o *Orchestrator) generatePantryOnly(ctx context.Context, req RecipeRequest) ([]Recipe, error) {
	resp, err := o.pantryAgent.GenerateRecipes(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Recipes, nil
}

// generateFlexible generates recipes with additional ingredients allowed
func (o *Orchestrator) generateFlexible(ctx context.Context, req RecipeRequest) ([]Recipe, error) {
	resp, err := o.flexibleAgent.GenerateRecipes(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Recipes, nil
}

// generateBoth generates recipes from both agents concurrently
func (o *Orchestrator) generateBoth(ctx context.Context, req RecipeRequest) ([]Recipe, []Recipe, error) {
	type result struct {
		recipes []Recipe
		err     error
		source  string
	}

	resultChan := make(chan result, 2)

	// Generate pantry-only recipes
	go func() {
		recipes, err := o.generatePantryOnly(ctx, req)
		resultChan <- result{recipes: recipes, err: err, source: "pantry_only"}
	}()

	// Generate flexible recipes
	go func() {
		recipes, err := o.generateFlexible(ctx, req)
		resultChan <- result{recipes: recipes, err: err, source: "flexible"}
	}()

	var pantryRecipes, flexibleRecipes []Recipe
	var lastErr error

	// Collect results
	for i := 0; i < 2; i++ {
		r := <-resultChan
		if r.err != nil {
			o.log.Error("Orchestrator: %s agent failed: %v", r.source, r.err)
			lastErr = r.err
			continue
		}
		if r.source == "pantry_only" {
			pantryRecipes = r.recipes
		} else {
			flexibleRecipes = r.recipes
		}
	}

	// Return error only if both failed
	if pantryRecipes == nil && flexibleRecipes == nil && lastErr != nil {
		return nil, nil, lastErr
	}

	return pantryRecipes, flexibleRecipes, nil
}

// QuickGenerate is a convenience method for quick recipe generation
func (o *Orchestrator) QuickGenerate(ctx context.Context, pantryItems []PantryItem, allergens []string) (*GenerateResult, error) {
	return o.Generate(ctx, GenerateRequest{
		RecipeRequest: RecipeRequest{
			PantryItems: pantryItems,
			Allergens:   allergens,
			RecipeCount: 2,
		},
		Mode: ModePantryOnly,
	})
}
