package recipes

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecipeNotFound = errors.New("recipe not found")
	ErrUnauthorized   = errors.New("unauthorized to access this recipe")
	ErrInvalidInput   = errors.New("invalid input")
)

// RecipeDifficulty represents the difficulty level
type RecipeDifficulty string

const (
	DifficultyEasy   RecipeDifficulty = "easy"
	DifficultyMedium RecipeDifficulty = "medium"
	DifficultyHard   RecipeDifficulty = "hard"
)

// RecipeSource represents which agent generated the recipe
type RecipeSource string

const (
	SourcePantryOnly  RecipeSource = "pantry_only"
	SourceFlexible    RecipeSource = "flexible"
	SourceSpoiling    RecipeSource = "spoiling"
	SourceUserCreated RecipeSource = "user_created"
)

// Recipe represents a saved recipe in the database
type Recipe struct {
	ID                 uuid.UUID        `json:"id"`
	UserID             string           `json:"user_id"`
	Title              string           `json:"title"`
	Description        string           `json:"description,omitempty"`
	Cuisine            string           `json:"cuisine,omitempty"`
	PrepTimeMinutes    *int             `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes    *int             `json:"cook_time_minutes,omitempty"`
	TotalTimeMinutes   int              `json:"total_time_minutes,omitempty"`
	Servings           int              `json:"servings"`
	Difficulty         RecipeDifficulty `json:"difficulty"`
	Ingredients        json.RawMessage  `json:"ingredients"`
	MissingIngredients json.RawMessage  `json:"missing_ingredients,omitempty"`
	Instructions       json.RawMessage  `json:"instructions"`
	CaloriesPerServing *float64         `json:"calories_per_serving,omitempty"`
	ProteinG           *float64         `json:"protein_g,omitempty"`
	CarbsG             *float64         `json:"carbs_g,omitempty"`
	FatG               *float64         `json:"fat_g,omitempty"`
	Source             RecipeSource     `json:"source"`
	AIModel            string           `json:"ai_model,omitempty"`
	IsFavorite         bool             `json:"is_favorite"`
	TimesCooked        int              `json:"times_cooked"`
	LastCookedAt       *time.Time       `json:"last_cooked_at,omitempty"`
	Rating             *int             `json:"rating,omitempty"`
	Notes              *string          `json:"notes,omitempty"`
	Tags               []string         `json:"tags"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

// Ingredient represents an ingredient in a recipe
type Ingredient struct {
	Name       string `json:"name"`
	Amount     string `json:"amount"`
	Unit       string `json:"unit,omitempty"`
	FromPantry bool   `json:"from_pantry"`
}

// CreateRecipeInput is the input for creating a new recipe
type CreateRecipeInput struct {
	Title              string           `json:"title"`
	Description        string           `json:"description,omitempty"`
	Cuisine            string           `json:"cuisine,omitempty"`
	PrepTimeMinutes    *int             `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes    *int             `json:"cook_time_minutes,omitempty"`
	Servings           int              `json:"servings"`
	Difficulty         RecipeDifficulty `json:"difficulty"`
	Ingredients        []Ingredient     `json:"ingredients"`
	MissingIngredients []Ingredient     `json:"missing_ingredients,omitempty"`
	Instructions       []string         `json:"instructions"`
	CaloriesPerServing *float64         `json:"calories_per_serving,omitempty"`
	ProteinG           *float64         `json:"protein_g,omitempty"`
	CarbsG             *float64         `json:"carbs_g,omitempty"`
	FatG               *float64         `json:"fat_g,omitempty"`
	Source             RecipeSource     `json:"source"`
	AIModel            string           `json:"ai_model,omitempty"`
	Tags               []string         `json:"tags,omitempty"`
}

// UpdateRecipeInput is the input for updating a recipe
type UpdateRecipeInput struct {
	IsFavorite *bool    `json:"is_favorite,omitempty"`
	Rating     *int     `json:"rating,omitempty"`
	Notes      *string  `json:"notes,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

// Repository handles database operations for recipes
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new recipe repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create saves a new recipe to the database
func (r *Repository) Create(ctx context.Context, userID string, input CreateRecipeInput) (*Recipe, error) {
	if input.Title == "" {
		return nil, ErrInvalidInput
	}

	ingredientsJSON, err := json.Marshal(input.Ingredients)
	if err != nil {
		return nil, err
	}

	missingIngredientsJSON, err := json.Marshal(input.MissingIngredients)
	if err != nil {
		return nil, err
	}

	instructionsJSON, err := json.Marshal(input.Instructions)
	if err != nil {
		return nil, err
	}

	// Default values
	if input.Servings == 0 {
		input.Servings = 2
	}
	if input.Difficulty == "" {
		input.Difficulty = DifficultyMedium
	}
	if input.Source == "" {
		input.Source = SourcePantryOnly
	}

	var recipe Recipe
	err = r.pool.QueryRow(ctx, `
		INSERT INTO recipes (
			user_id, title, description, cuisine,
			prep_time_minutes, cook_time_minutes, servings, difficulty,
			ingredients, missing_ingredients, instructions,
			calories_per_serving, protein_g, carbs_g, fat_g,
			source, ai_model, tags
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, user_id, title, description, cuisine,
		          prep_time_minutes, cook_time_minutes, total_time_minutes, servings, difficulty,
		          ingredients, missing_ingredients, instructions,
		          calories_per_serving, protein_g, carbs_g, fat_g,
		          source, ai_model, is_favorite, times_cooked, last_cooked_at,
		          rating, notes, tags, created_at, updated_at
	`, userID, input.Title, input.Description, input.Cuisine,
		input.PrepTimeMinutes, input.CookTimeMinutes, input.Servings, input.Difficulty,
		ingredientsJSON, missingIngredientsJSON, instructionsJSON,
		input.CaloriesPerServing, input.ProteinG, input.CarbsG, input.FatG,
		input.Source, input.AIModel, input.Tags).Scan(
		&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.Description, &recipe.Cuisine,
		&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.Servings, &recipe.Difficulty,
		&recipe.Ingredients, &recipe.MissingIngredients, &recipe.Instructions,
		&recipe.CaloriesPerServing, &recipe.ProteinG, &recipe.CarbsG, &recipe.FatG,
		&recipe.Source, &recipe.AIModel, &recipe.IsFavorite, &recipe.TimesCooked, &recipe.LastCookedAt,
		&recipe.Rating, &recipe.Notes, &recipe.Tags, &recipe.CreatedAt, &recipe.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

// GetByID retrieves a recipe by ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Recipe, error) {
	var recipe Recipe
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, title, description, cuisine,
		       prep_time_minutes, cook_time_minutes, total_time_minutes, servings, difficulty,
		       ingredients, missing_ingredients, instructions,
		       calories_per_serving, protein_g, carbs_g, fat_g,
		       source, ai_model, is_favorite, times_cooked, last_cooked_at,
		       rating, notes, tags, created_at, updated_at
		FROM recipes
		WHERE id = $1
	`, id).Scan(
		&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.Description, &recipe.Cuisine,
		&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.Servings, &recipe.Difficulty,
		&recipe.Ingredients, &recipe.MissingIngredients, &recipe.Instructions,
		&recipe.CaloriesPerServing, &recipe.ProteinG, &recipe.CarbsG, &recipe.FatG,
		&recipe.Source, &recipe.AIModel, &recipe.IsFavorite, &recipe.TimesCooked, &recipe.LastCookedAt,
		&recipe.Rating, &recipe.Notes, &recipe.Tags, &recipe.CreatedAt, &recipe.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecipeNotFound
		}
		return nil, err
	}

	return &recipe, nil
}

// ListByUserID retrieves all recipes for a user
func (r *Repository) ListByUserID(ctx context.Context, userID string) ([]Recipe, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, title, description, cuisine,
		       prep_time_minutes, cook_time_minutes, total_time_minutes, servings, difficulty,
		       ingredients, missing_ingredients, instructions,
		       calories_per_serving, protein_g, carbs_g, fat_g,
		       source, ai_model, is_favorite, times_cooked, last_cooked_at,
		       rating, notes, tags, created_at, updated_at
		FROM recipes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(
			&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.Description, &recipe.Cuisine,
			&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.Servings, &recipe.Difficulty,
			&recipe.Ingredients, &recipe.MissingIngredients, &recipe.Instructions,
			&recipe.CaloriesPerServing, &recipe.ProteinG, &recipe.CarbsG, &recipe.FatG,
			&recipe.Source, &recipe.AIModel, &recipe.IsFavorite, &recipe.TimesCooked, &recipe.LastCookedAt,
			&recipe.Rating, &recipe.Notes, &recipe.Tags, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// ListFavorites retrieves favorite recipes for a user
func (r *Repository) ListFavorites(ctx context.Context, userID string) ([]Recipe, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, title, description, cuisine,
		       prep_time_minutes, cook_time_minutes, total_time_minutes, servings, difficulty,
		       ingredients, missing_ingredients, instructions,
		       calories_per_serving, protein_g, carbs_g, fat_g,
		       source, ai_model, is_favorite, times_cooked, last_cooked_at,
		       rating, notes, tags, created_at, updated_at
		FROM recipes
		WHERE user_id = $1 AND is_favorite = TRUE
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(
			&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.Description, &recipe.Cuisine,
			&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.Servings, &recipe.Difficulty,
			&recipe.Ingredients, &recipe.MissingIngredients, &recipe.Instructions,
			&recipe.CaloriesPerServing, &recipe.ProteinG, &recipe.CarbsG, &recipe.FatG,
			&recipe.Source, &recipe.AIModel, &recipe.IsFavorite, &recipe.TimesCooked, &recipe.LastCookedAt,
			&recipe.Rating, &recipe.Notes, &recipe.Tags, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// Update updates a recipe
func (r *Repository) Update(ctx context.Context, id uuid.UUID, input UpdateRecipeInput) (*Recipe, error) {
	var recipe Recipe
	err := r.pool.QueryRow(ctx, `
		UPDATE recipes
		SET 
			is_favorite = COALESCE($2, is_favorite),
			rating = COALESCE($3, rating),
			notes = COALESCE($4, notes),
			tags = COALESCE($5, tags)
		WHERE id = $1
		RETURNING id, user_id, title, description, cuisine,
		          prep_time_minutes, cook_time_minutes, total_time_minutes, servings, difficulty,
		          ingredients, missing_ingredients, instructions,
		          calories_per_serving, protein_g, carbs_g, fat_g,
		          source, ai_model, is_favorite, times_cooked, last_cooked_at,
		          rating, notes, tags, created_at, updated_at
	`, id, input.IsFavorite, input.Rating, input.Notes, input.Tags).Scan(
		&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.Description, &recipe.Cuisine,
		&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.Servings, &recipe.Difficulty,
		&recipe.Ingredients, &recipe.MissingIngredients, &recipe.Instructions,
		&recipe.CaloriesPerServing, &recipe.ProteinG, &recipe.CarbsG, &recipe.FatG,
		&recipe.Source, &recipe.AIModel, &recipe.IsFavorite, &recipe.TimesCooked, &recipe.LastCookedAt,
		&recipe.Rating, &recipe.Notes, &recipe.Tags, &recipe.CreatedAt, &recipe.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecipeNotFound
		}
		return nil, err
	}

	return &recipe, nil
}

// MarkAsCooked increments the times_cooked counter and updates last_cooked_at
func (r *Repository) MarkAsCooked(ctx context.Context, id uuid.UUID) (*Recipe, error) {
	var recipe Recipe
	err := r.pool.QueryRow(ctx, `
		UPDATE recipes
		SET times_cooked = times_cooked + 1,
		    last_cooked_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, title, description, cuisine,
		          prep_time_minutes, cook_time_minutes, total_time_minutes, servings, difficulty,
		          ingredients, missing_ingredients, instructions,
		          calories_per_serving, protein_g, carbs_g, fat_g,
		          source, ai_model, is_favorite, times_cooked, last_cooked_at,
		          rating, notes, tags, created_at, updated_at
	`, id).Scan(
		&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.Description, &recipe.Cuisine,
		&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.Servings, &recipe.Difficulty,
		&recipe.Ingredients, &recipe.MissingIngredients, &recipe.Instructions,
		&recipe.CaloriesPerServing, &recipe.ProteinG, &recipe.CarbsG, &recipe.FatG,
		&recipe.Source, &recipe.AIModel, &recipe.IsFavorite, &recipe.TimesCooked, &recipe.LastCookedAt,
		&recipe.Rating, &recipe.Notes, &recipe.Tags, &recipe.CreatedAt, &recipe.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecipeNotFound
		}
		return nil, err
	}

	return &recipe, nil
}

// Delete removes a recipe
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM recipes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrRecipeNotFound
	}
	return nil
}
