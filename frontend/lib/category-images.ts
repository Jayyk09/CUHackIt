// Category to image mapping based on /public/categories/ directory
// Each category has one image

const categoryImageMap: Record<string, string> = {
  produce: '/categories/produce.png',
  dairy: '/categories/dairy.png',
  meat: '/categories/meat.png',
  seafood: '/categories/seafood.png',
  pantry: '/categories/pantry.png',
  frozen: '/categories/frozen.png',
  bakery: '/categories/bakery.png',
  snacks: '/categories/snacks.png',
  beverage: '/categories/beverage.png',
  deli: '/categories/deli.png',
  specialty: '/categories/specialty.png',
}

// Fallback image if category not found
const fallbackImage = '/categories/specialty.png'

/**
 * Get category image for a given category.
 */
export function getCategoryImage(category: string): string {
  const normalizedCategory = category.toLowerCase().trim()
  return categoryImageMap[normalizedCategory] || fallbackImage
}

/**
 * Get all available categories
 */
export function getAvailableCategories(): string[] {
  return Object.keys(categoryImageMap)
}
