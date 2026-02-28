// Mock data for pantry items and search results
// All items use the same Open Food Facts image URL for now

const OFF_IMAGE_URL = 'https://images.openfoodfacts.org/images/products/000/000/004/7838/front_fr.12.400.jpg'

export interface FoodItem {
  id: string
  product_name: string
  category: string
  quantity: number
  environmental_score: number
  nutriscore_score: number
  labels_en: string[]
  allergens_en: string[]
  traces_en: string[]
  shelf_life: number
  is_spoiled: boolean
  image_url: string
}

export const mockPantryItems: FoodItem[] = [
  {
    id: 'pantry-001',
    product_name: 'Organic Atlantic Salmon',
    category: 'seafood',
    quantity: 500,
    environmental_score: 72,
    nutriscore_score: 85,
    labels_en: ['Organic', 'Sustainable', 'Wild Caught'],
    allergens_en: ['Fish'],
    traces_en: ['Shellfish'],
    shelf_life: 14,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'pantry-002',
    product_name: 'Heirloom Tomatoes',
    category: 'produce',
    quantity: 1,
    environmental_score: 91,
    nutriscore_score: 95,
    labels_en: ['Organic', 'Local', 'Seasonal'],
    allergens_en: [],
    traces_en: [],
    shelf_life: 7,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'pantry-003',
    product_name: 'Free-Range Eggs',
    category: 'dairy',
    quantity: 12,
    environmental_score: 68,
    nutriscore_score: 78,
    labels_en: ['Free-Range', 'Farm Fresh', 'Omega-3'],
    allergens_en: ['Eggs'],
    traces_en: ['Milk'],
    shelf_life: 21,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'pantry-004',
    product_name: 'Extra Virgin Olive Oil',
    category: 'pantry',
    quantity: 500,
    environmental_score: 84,
    nutriscore_score: 92,
    labels_en: ['Cold Pressed', 'First Cold Press', 'DOP'],
    allergens_en: [],
    traces_en: [],
    shelf_life: 365,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'pantry-005',
    product_name: 'Grass-Fed Ground Beef',
    category: 'meat',
    quantity: 450,
    environmental_score: 52,
    nutriscore_score: 70,
    labels_en: ['Grass-Fed', 'Hormone-Free', 'Local Farm'],
    allergens_en: [],
    traces_en: [],
    shelf_life: 5,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'pantry-006',
    product_name: 'Aged Cheddar Cheese',
    category: 'dairy',
    quantity: 250,
    environmental_score: 58,
    nutriscore_score: 65,
    labels_en: ['Aged', 'Artisanal', 'Grass-Fed'],
    allergens_en: ['Milk'],
    traces_en: ['Lactose'],
    shelf_life: 180,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'pantry-007',
    product_name: 'Sourdough Bread',
    category: 'bakery',
    quantity: 1,
    environmental_score: 78,
    nutriscore_score: 72,
    labels_en: ['Artisan', 'Slow Fermented', 'No Preservatives'],
    allergens_en: ['Gluten', 'Wheat'],
    traces_en: ['Sesame'],
    shelf_life: 5,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'pantry-008',
    product_name: 'Frozen Wild Blueberries',
    category: 'frozen',
    quantity: 340,
    environmental_score: 85,
    nutriscore_score: 98,
    labels_en: ['Wild Harvested', 'No Added Sugar', 'Flash Frozen'],
    allergens_en: [],
    traces_en: [],
    shelf_life: 365,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
]

// Extended product database for search results
const searchableProducts: FoodItem[] = [
  ...mockPantryItems,
  {
    id: 'search-001',
    product_name: 'Organic Whole Milk',
    category: 'dairy',
    quantity: 1,
    environmental_score: 62,
    nutriscore_score: 75,
    labels_en: ['Organic', 'Grass-Fed', 'Non-Homogenized'],
    allergens_en: ['Milk'],
    traces_en: [],
    shelf_life: 14,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'search-002',
    product_name: 'Greek Yogurt',
    category: 'dairy',
    quantity: 500,
    environmental_score: 65,
    nutriscore_score: 82,
    labels_en: ['High Protein', 'Probiotic', 'No Added Sugar'],
    allergens_en: ['Milk'],
    traces_en: [],
    shelf_life: 21,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'search-003',
    product_name: 'Avocados',
    category: 'produce',
    quantity: 3,
    environmental_score: 55,
    nutriscore_score: 90,
    labels_en: ['Organic', 'Ripe'],
    allergens_en: [],
    traces_en: [],
    shelf_life: 5,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'search-004',
    product_name: 'Smoked Salmon',
    category: 'seafood',
    quantity: 200,
    environmental_score: 68,
    nutriscore_score: 80,
    labels_en: ['Wild Caught', 'Cold Smoked', 'Scottish'],
    allergens_en: ['Fish'],
    traces_en: [],
    shelf_life: 21,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'search-005',
    product_name: 'Organic Chicken Breast',
    category: 'meat',
    quantity: 500,
    environmental_score: 60,
    nutriscore_score: 88,
    labels_en: ['Organic', 'Free-Range', 'Air-Chilled'],
    allergens_en: [],
    traces_en: [],
    shelf_life: 7,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'search-006',
    product_name: 'Almond Butter',
    category: 'pantry',
    quantity: 340,
    environmental_score: 70,
    nutriscore_score: 78,
    labels_en: ['Raw', 'No Salt', 'Single Ingredient'],
    allergens_en: ['Tree Nuts'],
    traces_en: ['Peanuts'],
    shelf_life: 180,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'search-007',
    product_name: 'Dark Chocolate 85%',
    category: 'snacks',
    quantity: 100,
    environmental_score: 72,
    nutriscore_score: 68,
    labels_en: ['Fair Trade', 'Single Origin', 'Vegan'],
    allergens_en: ['Soy'],
    traces_en: ['Milk', 'Nuts'],
    shelf_life: 365,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'search-008',
    product_name: 'Sparkling Water',
    category: 'beverage',
    quantity: 1,
    environmental_score: 88,
    nutriscore_score: 100,
    labels_en: ['Natural', 'No Additives', 'Glass Bottle'],
    allergens_en: [],
    traces_en: [],
    shelf_life: 730,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'search-009',
    product_name: 'Prosciutto di Parma',
    category: 'deli',
    quantity: 150,
    environmental_score: 55,
    nutriscore_score: 62,
    labels_en: ['DOP', 'Aged 18 Months', 'Imported'],
    allergens_en: [],
    traces_en: [],
    shelf_life: 60,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
  {
    id: 'search-010',
    product_name: 'Truffle Oil',
    category: 'specialty',
    quantity: 100,
    environmental_score: 65,
    nutriscore_score: 70,
    labels_en: ['Black Truffle', 'Italian', 'Extra Virgin Base'],
    allergens_en: [],
    traces_en: [],
    shelf_life: 365,
    is_spoiled: false,
    image_url: OFF_IMAGE_URL,
  },
]

/**
 * Mock search function - simulates backend /search endpoint
 * Filters products by name, category, or labels
 */
export function mockSearchProducts(query: string): FoodItem[] {
  const normalizedQuery = query.toLowerCase().trim()
  
  if (normalizedQuery.length < 3) {
    return []
  }
  
  return searchableProducts.filter((item) => {
    const matchesName = item.product_name.toLowerCase().includes(normalizedQuery)
    const matchesCategory = item.category.toLowerCase().includes(normalizedQuery)
    const matchesLabels = item.labels_en.some((label) =>
      label.toLowerCase().includes(normalizedQuery)
    )
    
    return matchesName || matchesCategory || matchesLabels
  })
}
