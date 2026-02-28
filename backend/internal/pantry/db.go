package pantry

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrItemNotFound = errors.New("pantry item not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized to access this item")
	ErrUserNotFound = errors.New("user not found")
)

// PantryItemWithFood represents a pantry_items row joined with the foods table.
// This is what GET endpoints return.
type PantryItemWithFood struct {
	// pantry_items columns
	ID       int       `json:"id"`
	UserID   string    `json:"user_id"`
	FoodID   int64     `json:"food_id"`
	Quantity int       `json:"quantity"`
	IsFrozen bool      `json:"is_frozen"`
	AddedAt  time.Time `json:"added_at"`

	// foods columns (joined)
	ProductName            string   `json:"product_name"`
	EnvironmentalScore     *float64 `json:"environmental_score,omitempty"`
	NutriscoreScore        *float64 `json:"nutriscore_score,omitempty"`
	LabelsEn               []string `json:"labels_en"`
	AllergensEn            []string `json:"allergens_en"`
	TracesEn               []string `json:"traces_en"`
	ImageURL               *string  `json:"image_url,omitempty"`
	ImageSmallURL          *string  `json:"image_small_url,omitempty"`
	NormEnvironmentalScore *float64 `json:"norm_environmental_score,omitempty"`
	NormNutriscoreScore    *float64 `json:"norm_nutriscore_score,omitempty"`
	ShelfLife              *int     `json:"shelf_life,omitempty"`
	Category               *string  `json:"category,omitempty"`
}

// AddToPantryInput is the input for the POST /pantry endpoint
type AddToPantryInput struct {
	Auth0ID  string `json:"auth0_id"`
	FoodID   int64  `json:"food_id"`
	Quantity int    `json:"quantity"`
	IsFrozen bool   `json:"is_frozen"`
}

// SimplePantryEntry represents a raw row in the pantry_items table (no join)
type SimplePantryEntry struct {
	ID       int       `json:"id"`
	UserID   string    `json:"user_id"`
	FoodID   int64     `json:"food_id"`
	Quantity int       `json:"quantity"`
	IsFrozen bool      `json:"is_frozen"`
	AddedAt  time.Time `json:"added_at"`
}

// Repository handles database operations for pantry items
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new pantry repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// columns selected for the joined query
const pantryJoinSelect = `
	p.id, p.user_id, p.food_id, p.quantity, p.is_frozen, p.added_at,
	f.product_name,
	f.environmental_score,
	f.nutriscore_score,
	f.labels_en,
	f.allergens_en,
	f.traces_en,
	f.image_url,
	f.image_small_url,
	f.norm_environmental_score,
	f.norm_nutriscore,
	f.shelf_life,
	f.category
`

// scanPantryItemWithFood scans a row from the joined query into a PantryItemWithFood.
func scanPantryItemWithFood(scanner interface{ Scan(dest ...any) error }) (*PantryItemWithFood, error) {
	var item PantryItemWithFood
	err := scanner.Scan(
		&item.ID, &item.UserID, &item.FoodID, &item.Quantity, &item.IsFrozen, &item.AddedAt,
		&item.ProductName,
		&item.EnvironmentalScore,
		&item.NutriscoreScore,
		&item.LabelsEn,
		&item.AllergensEn,
		&item.TracesEn,
		&item.ImageURL,
		&item.ImageSmallURL,
		&item.NormEnvironmentalScore,
		&item.NormNutriscoreScore,
		&item.ShelfLife,
		&item.Category,
	)
	return &item, err
}

// ListByUserID retrieves all pantry items for a user, joined with food details.
func (r *Repository) ListByUserID(ctx context.Context, userID string) ([]PantryItemWithFood, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+pantryJoinSelect+`
		FROM pantry_items p
		JOIN foods f ON f.id = p.food_id
		WHERE p.user_id = $1
		ORDER BY f.product_name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PantryItemWithFood
	for rows.Next() {
		item, err := scanPantryItemWithFood(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}

// ListByCategory retrieves pantry items for a user filtered by food category.
func (r *Repository) ListByCategory(ctx context.Context, userID string, category string) ([]PantryItemWithFood, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+pantryJoinSelect+`
		FROM pantry_items p
		JOIN foods f ON f.id = p.food_id
		WHERE p.user_id = $1 AND f.category = $2
		ORDER BY f.product_name
	`, userID, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PantryItemWithFood
	for rows.Next() {
		item, err := scanPantryItemWithFood(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}

// GetByID retrieves a single pantry item by its pantry_items.id, joined with food details.
func (r *Repository) GetByID(ctx context.Context, id int) (*PantryItemWithFood, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT `+pantryJoinSelect+`
		FROM pantry_items p
		JOIN foods f ON f.id = p.food_id
		WHERE p.id = $1
	`, id)

	item, err := scanPantryItemWithFood(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return item, nil
}

// GetCategorySummary returns count of pantry items per food category for a user.
func (r *Repository) GetCategorySummary(ctx context.Context, userID string) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT f.category, COUNT(*) as count
		FROM pantry_items p
		JOIN foods f ON f.id = p.food_id
		WHERE p.user_id = $1 AND f.category IS NOT NULL
		GROUP BY f.category
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := make(map[string]int)
	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, err
		}
		summary[category] = count
	}

	return summary, nil
}

// AddToPantry looks up the user by auth0_id and inserts into the pantry table.
func (r *Repository) AddToPantry(ctx context.Context, input AddToPantryInput) (*SimplePantryEntry, error) {
	if input.Auth0ID == "" || input.FoodID <= 0 {
		return nil, ErrInvalidInput
	}

	if input.Quantity <= 0 {
		input.Quantity = 1
	}

	// Look up user_id from auth0_id
	var userID string
	err := r.pool.QueryRow(ctx, `SELECT id FROM users WHERE auth0_id = $1`, input.Auth0ID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Insert into pantry table
	var entry SimplePantryEntry
	err = r.pool.QueryRow(ctx, `
		INSERT INTO pantry_items (user_id, food_id, quantity, is_frozen)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, food_id, quantity, is_frozen, added_at
	`, userID, input.FoodID, input.Quantity, input.IsFrozen).Scan(
		&entry.ID, &entry.UserID, &entry.FoodID, &entry.Quantity, &entry.IsFrozen, &entry.AddedAt,
	)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// Delete removes a pantry item by its id.
func (r *Repository) Delete(ctx context.Context, id int) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM pantry_items WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrItemNotFound
	}
	return nil
}
