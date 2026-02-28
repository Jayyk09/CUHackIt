package agents

import (
	"strings"

	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

// AllergenFilter filters recipes based on user allergens
// This is a post-processing filter applied to generated recipes
type AllergenFilter struct {
	log *logger.Logger
}

// NewAllergenFilter creates a new AllergenFilter
func NewAllergenFilter(log *logger.Logger) *AllergenFilter {
	return &AllergenFilter{log: log}
}

// commonAllergenDerivatives maps allergens to their derivatives/related ingredients
var commonAllergenDerivatives = map[string][]string{
	"peanuts": {
		"peanut", "peanut butter", "peanut oil", "groundnut", "arachis oil",
	},
	"tree nuts": {
		"almond", "cashew", "walnut", "pecan", "pistachio", "hazelnut", "macadamia",
		"brazil nut", "pine nut", "chestnut", "almond milk", "almond butter",
		"cashew milk", "walnut oil",
	},
	"milk": {
		"milk", "dairy", "cheese", "butter", "cream", "yogurt", "ice cream",
		"whey", "casein", "lactose", "ghee", "sour cream", "half and half",
		"condensed milk", "evaporated milk", "cream cheese", "cottage cheese",
		"ricotta", "mozzarella", "parmesan", "cheddar", "brie", "feta",
	},
	"dairy": {
		"milk", "dairy", "cheese", "butter", "cream", "yogurt", "ice cream",
		"whey", "casein", "lactose", "ghee", "sour cream", "half and half",
		"condensed milk", "evaporated milk", "cream cheese", "cottage cheese",
		"ricotta", "mozzarella", "parmesan", "cheddar", "brie", "feta",
	},
	"eggs": {
		"egg", "eggs", "egg white", "egg yolk", "mayonnaise", "meringue",
		"albumin", "globulin", "lysozyme", "ovalbumin",
	},
	"wheat": {
		"wheat", "flour", "bread", "pasta", "noodle", "cracker", "cookie",
		"cake", "pastry", "couscous", "bulgur", "semolina", "durum",
		"farina", "seitan", "breadcrumb",
	},
	"gluten": {
		"wheat", "barley", "rye", "flour", "bread", "pasta", "noodle",
		"cracker", "cookie", "cake", "pastry", "couscous", "bulgur",
		"semolina", "seitan", "soy sauce", "malt",
	},
	"soy": {
		"soy", "soya", "soybean", "tofu", "tempeh", "edamame", "miso",
		"soy sauce", "soy milk", "soy lecithin", "tamari",
	},
	"fish": {
		"fish", "salmon", "tuna", "cod", "tilapia", "halibut", "sardine",
		"anchovy", "fish sauce", "worcestershire", "caesar dressing",
	},
	"shellfish": {
		"shellfish", "shrimp", "crab", "lobster", "clam", "mussel", "oyster",
		"scallop", "crawfish", "crayfish", "prawn",
	},
	"sesame": {
		"sesame", "sesame seed", "sesame oil", "tahini", "hummus", "halvah",
	},
}

// FilterRecipes filters out recipes that contain allergens
func (f *AllergenFilter) FilterRecipes(recipes []Recipe, allergens []string) []Recipe {
	if len(allergens) == 0 {
		return recipes
	}

	f.log.Info("AllergenFilter: Filtering %d recipes against %d allergens", len(recipes), len(allergens))

	// Build a set of all allergen-related terms
	allergenTerms := make(map[string]bool)
	for _, allergen := range allergens {
		allergenLower := strings.ToLower(strings.TrimSpace(allergen))
		allergenTerms[allergenLower] = true

		// Add derivatives
		if derivatives, ok := commonAllergenDerivatives[allergenLower]; ok {
			for _, derivative := range derivatives {
				allergenTerms[strings.ToLower(derivative)] = true
			}
		}
	}

	var safeRecipes []Recipe
	for _, recipe := range recipes {
		if f.isRecipeSafe(recipe, allergenTerms) {
			safeRecipes = append(safeRecipes, recipe)
		} else {
			f.log.Debug("AllergenFilter: Filtered out recipe '%s' due to allergen match", recipe.Title)
		}
	}

	f.log.Info("AllergenFilter: %d of %d recipes passed allergen check", len(safeRecipes), len(recipes))

	return safeRecipes
}

// isRecipeSafe checks if a recipe contains any allergen-related ingredients
func (f *AllergenFilter) isRecipeSafe(recipe Recipe, allergenTerms map[string]bool) bool {
	// Check all ingredients
	for _, ingredient := range recipe.Ingredients {
		if f.containsAllergen(ingredient.Name, allergenTerms) {
			return false
		}
	}

	// Check missing ingredients too
	for _, ingredient := range recipe.MissingIngredients {
		if f.containsAllergen(ingredient.Name, allergenTerms) {
			return false
		}
	}

	return true
}

// containsAllergen checks if an ingredient name contains any allergen terms
func (f *AllergenFilter) containsAllergen(ingredientName string, allergenTerms map[string]bool) bool {
	ingredientLower := strings.ToLower(ingredientName)

	// Check direct match
	if allergenTerms[ingredientLower] {
		return true
	}

	// Check if ingredient contains any allergen term
	for term := range allergenTerms {
		if strings.Contains(ingredientLower, term) {
			return true
		}
	}

	return false
}

// ValidateRecipeIngredients validates that a recipe doesn't contain any allergens
// Returns a list of problematic ingredients if any are found
func (f *AllergenFilter) ValidateRecipeIngredients(recipe Recipe, allergens []string) []string {
	if len(allergens) == 0 {
		return nil
	}

	// Build allergen terms
	allergenTerms := make(map[string]bool)
	for _, allergen := range allergens {
		allergenLower := strings.ToLower(strings.TrimSpace(allergen))
		allergenTerms[allergenLower] = true
		if derivatives, ok := commonAllergenDerivatives[allergenLower]; ok {
			for _, derivative := range derivatives {
				allergenTerms[strings.ToLower(derivative)] = true
			}
		}
	}

	var problematicIngredients []string

	for _, ingredient := range recipe.Ingredients {
		if f.containsAllergen(ingredient.Name, allergenTerms) {
			problematicIngredients = append(problematicIngredients, ingredient.Name)
		}
	}

	for _, ingredient := range recipe.MissingIngredients {
		if f.containsAllergen(ingredient.Name, allergenTerms) {
			problematicIngredients = append(problematicIngredients, ingredient.Name)
		}
	}

	return problematicIngredients
}
