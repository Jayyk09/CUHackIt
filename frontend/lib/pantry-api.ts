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
  norm_nutriscore_score?: number
  shelf_life?: number
  category?: string
}

export interface CategorySummary {
  [category: string]: number
}

/**
 * Calculate how many days remain before a pantry item expires.
 * Uses: added_at + shelf_life (with 4x multiplier if frozen).
 * Returns null if shelf_life is missing.
 */
export function getDaysRemaining(item: PantryItem): number | null {
  if (item.shelf_life == null) return null
  const shelfDays = item.is_frozen ? item.shelf_life * 4 : item.shelf_life
  const addedAt = new Date(item.added_at)
  const expiresAt = new Date(addedAt.getTime() + shelfDays * 24 * 60 * 60 * 1000)
  const now = new Date()
  return Math.ceil((expiresAt.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
}

/**
 * Returns pantry items that are expiring within the given threshold (default 3 days).
 */
export function getExpiringItems(items: PantryItem[], thresholdDays = 3): PantryItem[] {
  return items.filter((item) => {
    const days = getDaysRemaining(item)
    return days !== null && days <= thresholdDays
  })
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
