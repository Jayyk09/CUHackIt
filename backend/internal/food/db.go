package food

import (
	"strings"
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

// FoodProduct represents a food product from search
type FoodProduct struct {
	ID                 string       `json:"id"`
	Name               string       `json:"name"`
	Brand              string       `json:"brand,omitempty"`
	Barcode            string       `json:"barcode,omitempty"`
	Category           FoodCategory `json:"category"`
	ImageURL           string       `json:"image_url,omitempty"`
	CaloriesPerServing float64      `json:"calories_per_serving,omitempty"`
	ProteinG           float64      `json:"protein_g,omitempty"`
	CarbsG             float64      `json:"carbs_g,omitempty"`
	FatG               float64      `json:"fat_g,omitempty"`
	FiberG             float64      `json:"fiber_g,omitempty"`
	SugarG             float64      `json:"sugar_g,omitempty"`
	SodiumMg           float64      `json:"sodium_mg,omitempty"`
	ServingSize        string       `json:"serving_size,omitempty"`
	ShelfLifeDays      int          `json:"shelf_life_days,omitempty"`
}

// MockFoodDatabase is an in-memory mock database of food products
var MockFoodDatabase = []FoodProduct{
	// Produce
	{ID: "prod-001", Name: "Organic Bananas", Brand: "Dole", Category: CategoryProduce, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 105, ProteinG: 1.3, CarbsG: 27, FatG: 0.4, FiberG: 3.1, SugarG: 14, ServingSize: "1 medium", ShelfLifeDays: 7},
	{ID: "prod-002", Name: "Baby Spinach", Brand: "Earthbound Farm", Category: CategoryProduce, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 20, ProteinG: 2, CarbsG: 3, FatG: 0, FiberG: 2, ServingSize: "3 cups", ShelfLifeDays: 5},
	{ID: "prod-003", Name: "Avocados", Brand: "Hass", Category: CategoryProduce, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 240, ProteinG: 3, CarbsG: 12, FatG: 22, FiberG: 10, ServingSize: "1 avocado", ShelfLifeDays: 5},
	{ID: "prod-004", Name: "Lemons", Category: CategoryProduce, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 17, ProteinG: 0.6, CarbsG: 5, FatG: 0.2, ServingSize: "1 lemon", ShelfLifeDays: 21},
	{ID: "prod-005", Name: "Cherry Tomatoes", Category: CategoryProduce, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 27, ProteinG: 1, CarbsG: 6, FatG: 0.3, FiberG: 2, ServingSize: "1 cup", ShelfLifeDays: 7},
	{ID: "prod-006", Name: "Garlic", Category: CategoryProduce, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 4, ProteinG: 0.2, CarbsG: 1, FatG: 0, ServingSize: "1 clove", ShelfLifeDays: 60},
	{ID: "prod-007", Name: "Red Onion", Category: CategoryProduce, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 44, ProteinG: 1, CarbsG: 10, FatG: 0.1, FiberG: 2, ServingSize: "1 medium", ShelfLifeDays: 30},
	{ID: "prod-008", Name: "Bell Peppers", Category: CategoryProduce, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 24, ProteinG: 1, CarbsG: 6, FatG: 0.2, FiberG: 2, ServingSize: "1 medium", ShelfLifeDays: 14},

	// Dairy
	{ID: "dairy-001", Name: "Whole Milk", Brand: "Organic Valley", Category: CategoryDairy, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 150, ProteinG: 8, CarbsG: 12, FatG: 8, SugarG: 12, ServingSize: "1 cup", ShelfLifeDays: 14},
	{ID: "dairy-002", Name: "Greek Yogurt", Brand: "Fage", Category: CategoryDairy, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 100, ProteinG: 17, CarbsG: 6, FatG: 0.7, SugarG: 4, ServingSize: "170g", ShelfLifeDays: 21},
	{ID: "dairy-003", Name: "Sharp Cheddar", Brand: "Tillamook", Category: CategoryDairy, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 110, ProteinG: 7, CarbsG: 0, FatG: 9, SodiumMg: 180, ServingSize: "1 oz", ShelfLifeDays: 45},
	{ID: "dairy-004", Name: "Unsalted Butter", Brand: "Kerrygold", Category: CategoryDairy, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 100, ProteinG: 0, CarbsG: 0, FatG: 11, ServingSize: "1 tbsp", ShelfLifeDays: 60},
	{ID: "dairy-005", Name: "Large Eggs", Brand: "Pete & Gerry's", Category: CategoryDairy, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 70, ProteinG: 6, CarbsG: 0, FatG: 5, ServingSize: "1 egg", ShelfLifeDays: 35},
	{ID: "dairy-006", Name: "Heavy Cream", Brand: "Organic Valley", Category: CategoryDairy, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 50, ProteinG: 0, CarbsG: 0, FatG: 5, ServingSize: "1 tbsp", ShelfLifeDays: 21},

	// Meat
	{ID: "meat-001", Name: "Chicken Breast", Brand: "Bell & Evans", Category: CategoryMeat, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 165, ProteinG: 31, CarbsG: 0, FatG: 3.6, ServingSize: "4 oz", ShelfLifeDays: 3},
	{ID: "meat-002", Name: "Ground Beef 85/15", Brand: "Grass-fed", Category: CategoryMeat, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 240, ProteinG: 21, CarbsG: 0, FatG: 17, ServingSize: "4 oz", ShelfLifeDays: 3},
	{ID: "meat-003", Name: "Bacon", Brand: "Applegate", Category: CategoryMeat, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 60, ProteinG: 4, CarbsG: 0, FatG: 5, SodiumMg: 290, ServingSize: "2 slices", ShelfLifeDays: 14},
	{ID: "meat-004", Name: "Italian Sausage", Brand: "Johnsonville", Category: CategoryMeat, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 280, ProteinG: 16, CarbsG: 2, FatG: 23, ServingSize: "1 link", ShelfLifeDays: 7},

	// Seafood
	{ID: "sea-001", Name: "Atlantic Salmon", Brand: "Wild Caught", Category: CategorySeafood, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 208, ProteinG: 20, CarbsG: 0, FatG: 13, ServingSize: "4 oz", ShelfLifeDays: 2},
	{ID: "sea-002", Name: "Shrimp", Brand: "Wild Gulf", Category: CategorySeafood, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 84, ProteinG: 18, CarbsG: 0, FatG: 1, ServingSize: "4 oz", ShelfLifeDays: 2},

	// Pantry
	{ID: "pantry-001", Name: "Extra Virgin Olive Oil", Brand: "California Olive Ranch", Category: CategoryPantry, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 120, ProteinG: 0, CarbsG: 0, FatG: 14, ServingSize: "1 tbsp", ShelfLifeDays: 365},
	{ID: "pantry-002", Name: "Jasmine Rice", Brand: "Thai Kitchen", Category: CategoryPantry, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 160, ProteinG: 3, CarbsG: 36, FatG: 0, ServingSize: "1/4 cup dry", ShelfLifeDays: 730},
	{ID: "pantry-003", Name: "Spaghetti", Brand: "Barilla", Category: CategoryPantry, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 200, ProteinG: 7, CarbsG: 42, FatG: 1, FiberG: 2, ServingSize: "2 oz dry", ShelfLifeDays: 730},
	{ID: "pantry-004", Name: "Crushed Tomatoes", Brand: "San Marzano", Category: CategoryPantry, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 20, ProteinG: 1, CarbsG: 4, FatG: 0, FiberG: 1, ServingSize: "1/4 cup", ShelfLifeDays: 730},
	{ID: "pantry-005", Name: "Black Beans", Brand: "Goya", Category: CategoryPantry, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 110, ProteinG: 7, CarbsG: 20, FatG: 0.5, FiberG: 8, ServingSize: "1/2 cup", ShelfLifeDays: 730},
	{ID: "pantry-006", Name: "Chicken Broth", Brand: "Pacific Foods", Category: CategoryPantry, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 10, ProteinG: 1, CarbsG: 1, FatG: 0, SodiumMg: 570, ServingSize: "1 cup", ShelfLifeDays: 365},
	{ID: "pantry-007", Name: "Honey", Brand: "Local Raw", Category: CategoryPantry, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 60, ProteinG: 0, CarbsG: 17, FatG: 0, SugarG: 17, ServingSize: "1 tbsp", ShelfLifeDays: 1095},
	{ID: "pantry-008", Name: "Soy Sauce", Brand: "Kikkoman", Category: CategoryPantry, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 10, ProteinG: 1, CarbsG: 1, FatG: 0, SodiumMg: 920, ServingSize: "1 tbsp", ShelfLifeDays: 730},

	// Frozen
	{ID: "frozen-001", Name: "Mixed Vegetables", Brand: "Green Giant", Category: CategoryFrozen, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 60, ProteinG: 2, CarbsG: 10, FatG: 0, FiberG: 2, ServingSize: "3/4 cup", ShelfLifeDays: 365},
	{ID: "frozen-002", Name: "Frozen Blueberries", Brand: "Wyman's", Category: CategoryFrozen, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 70, ProteinG: 0, CarbsG: 17, FatG: 0, FiberG: 3, SugarG: 11, ServingSize: "1 cup", ShelfLifeDays: 365},

	// Bakery
	{ID: "bakery-001", Name: "Sourdough Bread", Brand: "Acme", Category: CategoryBakery, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 120, ProteinG: 4, CarbsG: 24, FatG: 0.5, FiberG: 1, ServingSize: "1 slice", ShelfLifeDays: 5},
	{ID: "bakery-002", Name: "Whole Wheat Tortillas", Brand: "Mission", Category: CategoryBakery, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 130, ProteinG: 4, CarbsG: 22, FatG: 3, FiberG: 3, ServingSize: "1 tortilla", ShelfLifeDays: 30},

	// Snacks
	{ID: "snack-001", Name: "Almonds", Brand: "Blue Diamond", Category: CategorySnacks, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 170, ProteinG: 6, CarbsG: 6, FatG: 15, FiberG: 3, ServingSize: "1/4 cup", ShelfLifeDays: 180},
	{ID: "snack-002", Name: "Dark Chocolate", Brand: "Ghirardelli", Category: CategorySnacks, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 210, ProteinG: 2, CarbsG: 24, FatG: 13, FiberG: 3, SugarG: 18, ServingSize: "3 squares", ShelfLifeDays: 365},

	// Beverage
	{ID: "bev-001", Name: "Orange Juice", Brand: "Tropicana", Category: CategoryBeverage, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 110, ProteinG: 2, CarbsG: 26, FatG: 0, SugarG: 22, ServingSize: "8 oz", ShelfLifeDays: 10},
	{ID: "bev-002", Name: "Oat Milk", Brand: "Oatly", Category: CategoryBeverage, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 120, ProteinG: 3, CarbsG: 16, FatG: 5, FiberG: 2, SugarG: 7, ServingSize: "1 cup", ShelfLifeDays: 10},

	// Deli
	{ID: "deli-001", Name: "Turkey Breast", Brand: "Boar's Head", Category: CategoryDeli, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 60, ProteinG: 12, CarbsG: 1, FatG: 0.5, SodiumMg: 440, ServingSize: "2 oz", ShelfLifeDays: 7},
	{ID: "deli-002", Name: "Swiss Cheese", Brand: "Boar's Head", Category: CategoryDeli, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 110, ProteinG: 8, CarbsG: 0, FatG: 8, SodiumMg: 60, ServingSize: "1 oz", ShelfLifeDays: 14},

	// Specialty
	{ID: "spec-001", Name: "Tofu", Brand: "House Foods", Category: CategorySpecialty, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 80, ProteinG: 8, CarbsG: 2, FatG: 4.5, ServingSize: "3 oz", ShelfLifeDays: 14},
	{ID: "spec-002", Name: "Miso Paste", Brand: "Miso Master", Category: CategorySpecialty, ImageURL: "https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg", CaloriesPerServing: 25, ProteinG: 2, CarbsG: 3, FatG: 1, SodiumMg: 630, ServingSize: "1 tbsp", ShelfLifeDays: 365},
}

// SearchProducts searches for food products by name
func SearchProducts(query string, limit int) []FoodProduct {
	if query == "" || len(query) < 3 {
		return []FoodProduct{}
	}

	query = strings.ToLower(query)
	var results []FoodProduct

	for _, product := range MockFoodDatabase {
		// Search in name and brand
		if strings.Contains(strings.ToLower(product.Name), query) ||
			strings.Contains(strings.ToLower(product.Brand), query) {
			results = append(results, product)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}

	return results
}

// GetProductByID retrieves a product by its ID
func GetProductByID(id string) *FoodProduct {
	for _, product := range MockFoodDatabase {
		if product.ID == id {
			return &product
		}
	}
	return nil
}

// GetProductsByCategory retrieves products by category
func GetProductsByCategory(category FoodCategory, limit int) []FoodProduct {
	var results []FoodProduct

	for _, product := range MockFoodDatabase {
		if product.Category == category {
			results = append(results, product)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}

	return results
}

// GetAllCategories returns all available categories with item counts
func GetAllCategories() map[FoodCategory]int {
	counts := make(map[FoodCategory]int)
	for _, product := range MockFoodDatabase {
		counts[product.Category]++
	}
	return counts
}
