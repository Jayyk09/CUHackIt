package pantry

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrItemNotFound = errors.New("pantry item not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized to access this item")
	ErrUserNotFound = errors.New("user not found")
)

// FoodCategory represents the category of a food item
type FoodCategory string

const (
	CategoryProduce   FoodCategory = "PRODUCE"
	CategoryDairy     FoodCategory = "DAIRY"
	CategoryMeat      FoodCategory = "MEAT"
	CategorySeafood   FoodCategory = "SEAFOOD"
	CategoryPantry    FoodCategory = "PANTRY"
	CategoryFrozen    FoodCategory = "FROZEN"
	CategoryBakery    FoodCategory = "BAKERY"
	CategorySnacks    FoodCategory = "SNACKS"
	CategoryBeverage  FoodCategory = "BEVERAGE"
	CategoryDeli      FoodCategory = "DELI"
	CategorySpecialty FoodCategory = "SPECIALTY"
)

// PantryItem represents an item in a user's pantry
type PantryItem struct {
	ID                 uuid.UUID    `json:"id"`
	UserID             string       `json:"user_id"`
	Name               string       `json:"name"`
	Brand              *string      `json:"brand,omitempty"`
	Barcode            *string      `json:"barcode,omitempty"`
	OpenFoodFactsID    *string      `json:"open_food_facts_id,omitempty"`
	Category           FoodCategory `json:"category"`
	Quantity           float64      `json:"quantity"`
	Unit               string       `json:"unit"`
	PurchaseDate       *time.Time   `json:"purchase_date,omitempty"`
	ExpirationDate     *time.Time   `json:"expiration_date,omitempty"`
	ShelfLifeDays      *int         `json:"shelf_life_days,omitempty"`
	CaloriesPerServing *float64     `json:"calories_per_serving,omitempty"`
	ProteinG           *float64     `json:"protein_g,omitempty"`
	CarbsG             *float64     `json:"carbs_g,omitempty"`
	FatG               *float64     `json:"fat_g,omitempty"`
	FiberG             *float64     `json:"fiber_g,omitempty"`
	SugarG             *float64     `json:"sugar_g,omitempty"`
	SodiumMg           *float64     `json:"sodium_mg,omitempty"`
	ServingSize        *string      `json:"serving_size,omitempty"`
	ImageURL           *string      `json:"image_url,omitempty"`
	IsExpired          bool         `json:"is_expired"`
	IsExpiringSoon     bool         `json:"is_expiring_soon"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}

// CreateItemInput is the input for creating a new pantry item
type CreateItemInput struct {
	Name               string       `json:"name"`
	Brand              *string      `json:"brand,omitempty"`
	Barcode            *string      `json:"barcode,omitempty"`
	OpenFoodFactsID    *string      `json:"open_food_facts_id,omitempty"`
	Category           FoodCategory `json:"category"`
	Quantity           float64      `json:"quantity"`
	Unit               string       `json:"unit"`
	ExpirationDate     *time.Time   `json:"expiration_date,omitempty"`
	ShelfLifeDays      *int         `json:"shelf_life_days,omitempty"`
	CaloriesPerServing *float64     `json:"calories_per_serving,omitempty"`
	ProteinG           *float64     `json:"protein_g,omitempty"`
	CarbsG             *float64     `json:"carbs_g,omitempty"`
	FatG               *float64     `json:"fat_g,omitempty"`
	FiberG             *float64     `json:"fiber_g,omitempty"`
	SugarG             *float64     `json:"sugar_g,omitempty"`
	SodiumMg           *float64     `json:"sodium_mg,omitempty"`
	ServingSize        *string      `json:"serving_size,omitempty"`
	ImageURL           *string      `json:"image_url,omitempty"`
}

// UpdateItemInput is the input for updating a pantry item
type UpdateItemInput struct {
	Name           *string       `json:"name,omitempty"`
	Brand          *string       `json:"brand,omitempty"`
	Category       *FoodCategory `json:"category,omitempty"`
	Quantity       *float64      `json:"quantity,omitempty"`
	Unit           *string       `json:"unit,omitempty"`
	ExpirationDate *time.Time    `json:"expiration_date,omitempty"`
}

// Repository handles database operations for pantry items
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new pantry repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create creates a new pantry item
func (r *Repository) Create(ctx context.Context, userID string, input CreateItemInput) (*PantryItem, error) {
	if input.Name == "" || input.Category == "" {
		return nil, ErrInvalidInput
	}

	// Default values
	if input.Quantity == 0 {
		input.Quantity = 1
	}
	if input.Unit == "" {
		input.Unit = "item"
	}

	var item PantryItem
	err := r.pool.QueryRow(ctx, `
		INSERT INTO pantry_items (
			user_id, name, brand, barcode, open_food_facts_id, category,
			quantity, unit, expiration_date, shelf_life_days,
			calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg,
			serving_size, image_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id, user_id, name, brand, barcode, open_food_facts_id, category,
		          quantity, unit, purchase_date, expiration_date, shelf_life_days,
		          calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg,
		          serving_size, image_url, is_expired, is_expiring_soon, created_at, updated_at
	`, userID, input.Name, input.Brand, input.Barcode, input.OpenFoodFactsID, input.Category,
		input.Quantity, input.Unit, input.ExpirationDate, input.ShelfLifeDays,
		input.CaloriesPerServing, input.ProteinG, input.CarbsG, input.FatG, input.FiberG, input.SugarG, input.SodiumMg,
		input.ServingSize, input.ImageURL).Scan(
		&item.ID, &item.UserID, &item.Name, &item.Brand, &item.Barcode, &item.OpenFoodFactsID, &item.Category,
		&item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpirationDate, &item.ShelfLifeDays,
		&item.CaloriesPerServing, &item.ProteinG, &item.CarbsG, &item.FatG, &item.FiberG, &item.SugarG, &item.SodiumMg,
		&item.ServingSize, &item.ImageURL, &item.IsExpired, &item.IsExpiringSoon, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// GetByID retrieves a pantry item by its ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*PantryItem, error) {
	var item PantryItem
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, name, brand, barcode, open_food_facts_id, category,
		       quantity, unit, purchase_date, expiration_date, shelf_life_days,
		       calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg,
		       serving_size, image_url, is_expired, is_expiring_soon, created_at, updated_at
		FROM pantry_items
		WHERE id = $1
	`, id).Scan(
		&item.ID, &item.UserID, &item.Name, &item.Brand, &item.Barcode, &item.OpenFoodFactsID, &item.Category,
		&item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpirationDate, &item.ShelfLifeDays,
		&item.CaloriesPerServing, &item.ProteinG, &item.CarbsG, &item.FatG, &item.FiberG, &item.SugarG, &item.SodiumMg,
		&item.ServingSize, &item.ImageURL, &item.IsExpired, &item.IsExpiringSoon, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	return &item, nil
}

// ListByUserID retrieves all pantry items for a user
func (r *Repository) ListByUserID(ctx context.Context, userID string) ([]PantryItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, brand, barcode, open_food_facts_id, category,
		       quantity, unit, purchase_date, expiration_date, shelf_life_days,
		       calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg,
		       serving_size, image_url, is_expired, is_expiring_soon, created_at, updated_at
		FROM pantry_items
		WHERE user_id = $1
		ORDER BY category, name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PantryItem
	for rows.Next() {
		var item PantryItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.Name, &item.Brand, &item.Barcode, &item.OpenFoodFactsID, &item.Category,
			&item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpirationDate, &item.ShelfLifeDays,
			&item.CaloriesPerServing, &item.ProteinG, &item.CarbsG, &item.FatG, &item.FiberG, &item.SugarG, &item.SodiumMg,
			&item.ServingSize, &item.ImageURL, &item.IsExpired, &item.IsExpiringSoon, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// ListByCategory retrieves pantry items for a user filtered by category
func (r *Repository) ListByCategory(ctx context.Context, userID string, category FoodCategory) ([]PantryItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, brand, barcode, open_food_facts_id, category,
		       quantity, unit, purchase_date, expiration_date, shelf_life_days,
		       calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg,
		       serving_size, image_url, is_expired, is_expiring_soon, created_at, updated_at
		FROM pantry_items
		WHERE user_id = $1 AND category = $2
		ORDER BY name
	`, userID, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PantryItem
	for rows.Next() {
		var item PantryItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.Name, &item.Brand, &item.Barcode, &item.OpenFoodFactsID, &item.Category,
			&item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpirationDate, &item.ShelfLifeDays,
			&item.CaloriesPerServing, &item.ProteinG, &item.CarbsG, &item.FatG, &item.FiberG, &item.SugarG, &item.SodiumMg,
			&item.ServingSize, &item.ImageURL, &item.IsExpired, &item.IsExpiringSoon, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// ListExpiringSoon retrieves pantry items that are expiring within 3 days
func (r *Repository) ListExpiringSoon(ctx context.Context, userID string) ([]PantryItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, brand, barcode, open_food_facts_id, category,
		       quantity, unit, purchase_date, expiration_date, shelf_life_days,
		       calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg,
		       serving_size, image_url, is_expired, is_expiring_soon, created_at, updated_at
		FROM pantry_items
		WHERE user_id = $1 AND is_expiring_soon = TRUE
		ORDER BY expiration_date
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PantryItem
	for rows.Next() {
		var item PantryItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.Name, &item.Brand, &item.Barcode, &item.OpenFoodFactsID, &item.Category,
			&item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpirationDate, &item.ShelfLifeDays,
			&item.CaloriesPerServing, &item.ProteinG, &item.CarbsG, &item.FatG, &item.FiberG, &item.SugarG, &item.SodiumMg,
			&item.ServingSize, &item.ImageURL, &item.IsExpired, &item.IsExpiringSoon, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Update updates a pantry item
func (r *Repository) Update(ctx context.Context, id uuid.UUID, input UpdateItemInput) (*PantryItem, error) {
	var item PantryItem
	err := r.pool.QueryRow(ctx, `
		UPDATE pantry_items
		SET 
			name = COALESCE($2, name),
			brand = COALESCE($3, brand),
			category = COALESCE($4, category),
			quantity = COALESCE($5, quantity),
			unit = COALESCE($6, unit),
			expiration_date = COALESCE($7, expiration_date)
		WHERE id = $1
		RETURNING id, user_id, name, brand, barcode, open_food_facts_id, category,
		          quantity, unit, purchase_date, expiration_date, shelf_life_days,
		          calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg,
		          serving_size, image_url, is_expired, is_expiring_soon, created_at, updated_at
	`, id, input.Name, input.Brand, input.Category, input.Quantity, input.Unit, input.ExpirationDate).Scan(
		&item.ID, &item.UserID, &item.Name, &item.Brand, &item.Barcode, &item.OpenFoodFactsID, &item.Category,
		&item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpirationDate, &item.ShelfLifeDays,
		&item.CaloriesPerServing, &item.ProteinG, &item.CarbsG, &item.FatG, &item.FiberG, &item.SugarG, &item.SodiumMg,
		&item.ServingSize, &item.ImageURL, &item.IsExpired, &item.IsExpiringSoon, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	return &item, nil
}

// Delete removes a pantry item
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM pantry_items WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrItemNotFound
	}
	return nil
}

// AddToPantryInput is the input for the POST /pantry endpoint
type AddToPantryInput struct {
	Auth0ID  string `json:"auth0_id"`
	FoodID   int64  `json:"food_id"`
	Quantity int    `json:"quantity"`
	IsFrozen bool   `json:"is_frozen"`
}

// SimplePantryEntry represents a row in the pantry table
type SimplePantryEntry struct {
	ID       int       `json:"id"`
	UserID   string    `json:"user_id"`
	FoodID   int64     `json:"food_id"`
	Quantity int       `json:"quantity"`
	IsFrozen bool      `json:"is_frozen"`
	AddedAt  time.Time `json:"added_at"`
}

// AddToPantry looks up the user by auth0_id and inserts into the pantry table
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

// GetCategorySummary returns count of items per category for a user
func (r *Repository) GetCategorySummary(ctx context.Context, userID string) (map[FoodCategory]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT category, COUNT(*) as count
		FROM pantry_items
		WHERE user_id = $1
		GROUP BY category
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := make(map[FoodCategory]int)
	for rows.Next() {
		var category FoodCategory
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, err
		}
		summary[category] = count
	}

	return summary, nil
}
