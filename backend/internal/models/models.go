package models

import "time"

// Account represents a user account in the database.
type Account struct {
	AccountID int      `json:"account_id"`
	Labels    []string `json:"labels"`
	Allergens []string `json:"allergens"`
}

// CreateAccountRequest is the payload for POST /users.
type CreateAccountRequest struct {
	Labels    []string `json:"labels"`
	Allergens []string `json:"allergens"`
}

// UpdateAccountRequest is the payload for PUT /users/{id}.
type UpdateAccountRequest struct {
	Labels    []string `json:"labels"`
	Allergens []string `json:"allergens"`
}

// Food represents a food item from the food table.
type Food struct {
	ID                   int      `json:"id"`
	ProductName          string   `json:"product_name"`
	EnvironmentalScore   float64  `json:"environmental_score"`
	NutriscoreScore      float64  `json:"nutriscore_score"`
	LabelsEn             []string `json:"labels_en"`
	AllergensEn          []string `json:"allergens_en"`
	TracesEn             []string `json:"traces_en"`
	ImageURL             *string  `json:"image_url"`
	ImageSmallURL        *string  `json:"image_small_url"`
	ShelfLife            *int     `json:"shelf_life"`
	Category             *string  `json:"category"`
}

// CreateFoodRequest is the payload for POST /food/{id}.
type CreateFoodRequest struct {
	ProductName        string   `json:"product_name"`
	EnvironmentalScore float64  `json:"environmental_score"`
	NutriscoreScore    float64  `json:"nutriscore_score"`
	LabelsEn           []string `json:"labels_en"`
	AllergensEn        []string `json:"allergens_en"`
	TracesEn           []string `json:"traces_en"`
	ImageURL           *string  `json:"image_url"`
	ImageSmallURL      *string  `json:"image_small_url"`
	ShelfLife          *int     `json:"shelf_life"`
	Category           *string  `json:"category"`
}

// PantryItem represents a row in the pantry table.
type PantryItem struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	FoodID   int       `json:"food_id"`
	AddedAt  time.Time `json:"added_at"`
	Quantity int       `json:"quantity"`
	Units    string    `json:"units"`
	Category *string   `json:"category"`
	IsFrozen bool      `json:"is_frozen"`
}

// CreatePantryRequest is the payload for POST /pantry.
type CreatePantryRequest struct {
	UserID   int     `json:"user_id"`
	FoodID   int     `json:"food_id"`
	Quantity int     `json:"quantity"`
	Units    string  `json:"units"`
	Category *string `json:"category"`
	IsFrozen bool    `json:"is_frozen"`
}

// UpdatePantryRequest is the payload for PATCH /pantry/{id}.
type UpdatePantryRequest struct {
	Quantity *int    `json:"quantity"`
	Units    *string `json:"units"`
	Category *string `json:"category"`
	IsFrozen *bool   `json:"is_frozen"`
}

// PantryResponse is the JSON returned from GET /pantry, joining pantry + food data.
// This matches the example JSON in the spec.
type PantryResponse struct {
	ProductName        string   `json:"product_name"`
	EnvironmentalScore float64  `json:"environmental_score"`
	NutriscoreScore    float64  `json:"nutriscore_score"`
	LabelsEn           []string `json:"labels_en"`
	AllergensEn        []string `json:"allergens_en"`
	TracesEn           []string `json:"traces_en"`
	ImageURL           *string  `json:"image_url"`
	ShelfLife          *int     `json:"shelf_life"`
	Category           *string  `json:"category"`
	Quantity           int      `json:"quantity"`
	Units              string   `json:"units"`
	IsSpoiled          bool     `json:"is_spoiled"`
	IsFrozen           bool     `json:"is_frozen"`
	AddedAt            string   `json:"added_at"`
}
