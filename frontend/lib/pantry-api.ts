const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

/** Pantry item joined with food details from the backend */
export interface PantryItem {
  // pantry_items columns
  id: number
  user_id: string
  food_id: number
  quantity: number
  is_frozen: boolean
  added_at: string

  // foods columns (joined)
  product_name: string
  environmental_score?: number
  nutriscore_score?: number
  labels_en: string[]
  allergens_en: string[]
  traces_en: string[]
  image_url?: string
  image_small_url?: string
  norm_environmental_score?: number
  shelf_life?: number
  category?: string
}

export interface CategorySummary {
  [category: string]: number
}

/**
 * List all pantry items for a user. Optionally filter by category.
 */
export async function listPantryItems(
  userId: string,
  category?: string
): Promise<PantryItem[]> {
  const url = category
    ? `${API_BASE}/users/${userId}/pantry?category=${encodeURIComponent(category)}`
    : `${API_BASE}/users/${userId}/pantry`

  const res = await fetch(url)
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Failed to fetch pantry items')
  }
  const data = await res.json()
  return Array.isArray(data) ? data : []
}

/**
 * Get a count of pantry items grouped by category.
 */
export async function getCategorySummary(userId: string): Promise<CategorySummary> {
  const res = await fetch(`${API_BASE}/users/${userId}/pantry/summary`)
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Failed to fetch category summary')
  }
  return res.json()
}
