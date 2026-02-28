package users

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Jayyk09/CUHackIt/internal/database"
)

// User represents a row in the users table.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Labels    []string  `json:"labels"`
	Allergens []string  `json:"allergens"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUser inserts a new user. Labels and allergens default to empty.
func CreateUser(ctx context.Context, db *database.DB, id, email, name string) (*User, error) {
	u := &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Labels:    []string{},
		Allergens: []string{},
	}

	err := db.Pool.QueryRow(ctx,
		`INSERT INTO users (id, email, name)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (id) DO NOTHING
		 RETURNING created_at, updated_at`,
		u.ID, u.Email, u.Name,
	).Scan(&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		// ON CONFLICT DO NOTHING returns no rows for duplicates — treat as existing user.
		if errors.Is(err, pgx.ErrNoRows) {
			return GetUser(ctx, db, id)
		}
		return nil, err
	}
	return u, nil
}

// GetUser retrieves a user by ID (Auth0 sub).
func GetUser(ctx context.Context, db *database.DB, id string) (*User, error) {
	u := &User{}
	err := db.Pool.QueryRow(ctx,
		`SELECT id, email, name, labels, allergens, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Labels, &u.Allergens, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, err
	}
	return u, nil
}

// UserExists checks whether a user row exists.
func UserExists(ctx context.Context, db *database.DB, id string) (bool, error) {
	var exists bool
	err := db.Pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id,
	).Scan(&exists)
	return exists, err
}

// UpdatePreferences sets the labels and allergens arrays.
func UpdatePreferences(ctx context.Context, db *database.DB, id string, labels, allergens []string) (*User, error) {
	if labels == nil {
		labels = []string{}
	}
	if allergens == nil {
		allergens = []string{}
	}

	u := &User{}
	err := db.Pool.QueryRow(ctx,
		`INSERT INTO users (id, labels, allergens)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (id) DO UPDATE
		   SET labels = EXCLUDED.labels,
		       allergens = EXCLUDED.allergens,
		       updated_at = now()
		 RETURNING id, email, name, labels, allergens, created_at, updated_at`,
		id, labels, allergens,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Labels, &u.Allergens, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}
