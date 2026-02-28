package database

import "context"

// Migrate creates tables if they don't already exist.
// Suitable for development / hackathon use — swap for a proper
// migration tool (goose, golang-migrate, etc.) in production.
func (db *DB) Migrate(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id         TEXT PRIMARY KEY,
			email      TEXT NOT NULL DEFAULT '',
			name       TEXT NOT NULL DEFAULT '',
			labels     TEXT[] NOT NULL DEFAULT '{}',
			allergens  TEXT[] NOT NULL DEFAULT '{}',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS email      TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS name       TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS labels     TEXT[] NOT NULL DEFAULT '{}'`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS allergens  TEXT[] NOT NULL DEFAULT '{}'`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT now()`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`,
	}

	for _, stmt := range statements {
		if _, err := db.Pool.Exec(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}
