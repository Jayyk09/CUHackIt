package food

import (
	"context"
	"fmt"

	"github.com/Jayyk09/CUHackIt/internal/database"
)

const foodTable = "foods"

type Product struct {
	ID                     int64    `json:"id"`
	ProductName            string   `json:"product_name"`
	NormEnvironmentalScore *float64 `json:"norm_environmental_score"`
	NutriscoreScore        *string  `json:"nutriscore_score"`
	LabelsEn               []string `json:"labels_en"`
	AllergensEn            []string `json:"allergens_en"`
	TracesEn               []string `json:"traces_en"`
	ImageURL               *string  `json:"image_url"`
	ImageSmallURL          *string  `json:"image_small_url"`
	ShelfLife              *int     `json:"shelf_life"`
	Category               *string  `json:"category"`
}

func fetchProducts(ctx context.Context, db *database.DB, search string, limit, offset int) ([]Product, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	queryBase := fmt.Sprintf(`SELECT
		id,
		product_name,
		norm_environmental_score,
		nutriscore_score,
		labels_en,
		allergens_en,
		traces_en,
		image_url,
		image_small_url,
		shelf_life,
		category
		FROM %s`, foodTable)

	var rows interface {
		Close()
		Next() bool
		Scan(...any) error
		Err() error
	}
	var err error

	if search == "" {
		rows, err = db.Pool.Query(ctx, queryBase+" ORDER BY product_name LIMIT $1 OFFSET $2", limit, offset)
	} else {
		rows, err = db.Pool.Query(ctx, queryBase+" WHERE product_name ILIKE $1 ORDER BY product_name LIMIT $2 OFFSET $3", "%"+search+"%", limit, offset)
	}
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	products := make([]Product, 0, limit)
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.NormEnvironmentalScore,
			&product.NutriscoreScore,
			&product.LabelsEn,
			&product.AllergensEn,
			&product.TracesEn,
			&product.ImageURL,
			&product.ImageSmallURL,
			&product.ShelfLife,
			&product.Category,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}

	return products, nil
}

func updateProductMetadata(ctx context.Context, db *database.DB, id int64, category *string, shelfLife *int) error {
	if db == nil || db.Pool == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	if category == nil && shelfLife == nil {
		return nil
	}

	query := fmt.Sprintf("UPDATE %s SET", foodTable)
	args := make([]any, 0, 3)
	idx := 1

	if category != nil {
		query += fmt.Sprintf(" category = $%d", idx)
		args = append(args, *category)
		idx++
	}
	if shelfLife != nil {
		if len(args) > 0 {
			query += ","
		}
		query += fmt.Sprintf(" shelf_life = $%d", idx)
		args = append(args, *shelfLife)
		idx++
	}

	query += fmt.Sprintf(" WHERE id = $%d", idx)
	args = append(args, id)

	if _, err := db.Pool.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("update product metadata: %w", err)
	}

	return nil
}
