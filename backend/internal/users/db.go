package users

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidInput      = errors.New("invalid input")
)

// User represents a user in the system with profile information
type User struct {
	ID                   string    `json:"id"`
	Auth0ID              string    `json:"auth0_id"`
	Email                string    `json:"email"`
	Name                 string    `json:"name,omitempty"`
	Allergens            []string  `json:"allergens"`
	DietaryPreferences   []string  `json:"dietary_preferences"`
	NutritionalGoals     []string  `json:"nutritional_goals"`
	CookingSkill         string    `json:"cooking_skill"`
	CuisinePreferences   []string  `json:"cuisine_preferences"`
	OnboardingCompleted  bool      `json:"onboarding_completed"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// CreateUserInput is the input for creating a new user
type CreateUserInput struct {
	Auth0ID string `json:"auth0_id"`
	Email   string `json:"email"`
	Name    string `json:"name,omitempty"`
}

// UpdateProfileInput is the input for updating user profile (onboarding)
type UpdateProfileInput struct {
	Name               *string  `json:"name,omitempty"`
	Allergens          []string `json:"allergens,omitempty"`
	DietaryPreferences []string `json:"dietary_preferences,omitempty"`
	NutritionalGoals   []string `json:"nutritional_goals,omitempty"`
	CookingSkill       *string  `json:"cooking_skill,omitempty"`
	CuisinePreferences []string `json:"cuisine_preferences,omitempty"`
}

// Repository handles database operations for users
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new user repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create creates a new user
func (r *Repository) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	if input.Auth0ID == "" || input.Email == "" {
		return nil, ErrInvalidInput
	}

	var user User
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (auth0_id, email, name)
		VALUES ($1, $2, $3)
		RETURNING id, auth0_id, email, name, allergens, dietary_preferences, 
		          nutritional_goals, cooking_skill, cuisine_preferences, 
		          onboarding_completed, created_at, updated_at
	`, input.Auth0ID, input.Email, input.Name).Scan(
		&user.ID,
		&user.Auth0ID,
		&user.Email,
		&user.Name,
		&user.Allergens,
		&user.DietaryPreferences,
		&user.NutritionalGoals,
		&user.CookingSkill,
		&user.CuisinePreferences,
		&user.OnboardingCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user by their ID
func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `
		SELECT id, auth0_id, email, name, allergens, dietary_preferences,
		       nutritional_goals, cooking_skill, cuisine_preferences,
		       onboarding_completed, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Auth0ID,
		&user.Email,
		&user.Name,
		&user.Allergens,
		&user.DietaryPreferences,
		&user.NutritionalGoals,
		&user.CookingSkill,
		&user.CuisinePreferences,
		&user.OnboardingCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByAuth0ID retrieves a user by their Auth0 ID
func (r *Repository) GetByAuth0ID(ctx context.Context, auth0ID string) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `
		SELECT id, auth0_id, email, name, allergens, dietary_preferences,
		       nutritional_goals, cooking_skill, cuisine_preferences,
		       onboarding_completed, created_at, updated_at
		FROM users
		WHERE auth0_id = $1
	`, auth0ID).Scan(
		&user.ID,
		&user.Auth0ID,
		&user.Email,
		&user.Name,
		&user.Allergens,
		&user.DietaryPreferences,
		&user.NutritionalGoals,
		&user.CookingSkill,
		&user.CuisinePreferences,
		&user.OnboardingCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by their email
func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `
		SELECT id, auth0_id, email, name, allergens, dietary_preferences,
		       nutritional_goals, cooking_skill, cuisine_preferences,
		       onboarding_completed, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Auth0ID,
		&user.Email,
		&user.Name,
		&user.Allergens,
		&user.DietaryPreferences,
		&user.NutritionalGoals,
		&user.CookingSkill,
		&user.CuisinePreferences,
		&user.OnboardingCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// UpdateProfile updates a user's profile information
func (r *Repository) UpdateProfile(ctx context.Context, id string, input UpdateProfileInput) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `
		UPDATE users
		SET 
			name = COALESCE($2, name),
			allergens = COALESCE($3, allergens),
			dietary_preferences = COALESCE($4, dietary_preferences),
			nutritional_goals = COALESCE($5, nutritional_goals),
			cooking_skill = COALESCE($6, cooking_skill),
			cuisine_preferences = COALESCE($7, cuisine_preferences)
		WHERE id = $1
		RETURNING id, auth0_id, email, name, allergens, dietary_preferences,
		          nutritional_goals, cooking_skill, cuisine_preferences,
		          onboarding_completed, created_at, updated_at
	`, id, input.Name, input.Allergens, input.DietaryPreferences,
		input.NutritionalGoals, input.CookingSkill, input.CuisinePreferences).Scan(
		&user.ID,
		&user.Auth0ID,
		&user.Email,
		&user.Name,
		&user.Allergens,
		&user.DietaryPreferences,
		&user.NutritionalGoals,
		&user.CookingSkill,
		&user.CuisinePreferences,
		&user.OnboardingCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// CompleteOnboarding marks a user's onboarding as complete
func (r *Repository) CompleteOnboarding(ctx context.Context, id string, input UpdateProfileInput) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `
		UPDATE users
		SET 
			name = COALESCE($2, name),
			allergens = COALESCE($3, allergens),
			dietary_preferences = COALESCE($4, dietary_preferences),
			nutritional_goals = COALESCE($5, nutritional_goals),
			cooking_skill = COALESCE($6, cooking_skill),
			cuisine_preferences = COALESCE($7, cuisine_preferences),
			onboarding_completed = TRUE
		WHERE id = $1
		RETURNING id, auth0_id, email, name, allergens, dietary_preferences,
		          nutritional_goals, cooking_skill, cuisine_preferences,
		          onboarding_completed, created_at, updated_at
	`, id, input.Name, input.Allergens, input.DietaryPreferences,
		input.NutritionalGoals, input.CookingSkill, input.CuisinePreferences).Scan(
		&user.ID,
		&user.Auth0ID,
		&user.Email,
		&user.Name,
		&user.Allergens,
		&user.DietaryPreferences,
		&user.NutritionalGoals,
		&user.CookingSkill,
		&user.CuisinePreferences,
		&user.OnboardingCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Delete removes a user from the database
func (r *Repository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

// FindOrCreate finds a user by Auth0 ID or creates a new one
func (r *Repository) FindOrCreate(ctx context.Context, input CreateUserInput) (*User, error) {
	// Try to find existing user
	user, err := r.GetByAuth0ID(ctx, input.Auth0ID)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// Create new user
	return r.Create(ctx, input)
}
