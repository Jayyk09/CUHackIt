const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

export interface PantryItem {
  id: string
  user_id: string
  name: string
  brand?: string
  barcode?: string
  open_food_facts_id?: string
  category: string
  quantity: number
  unit: string
  purchase_date?: string
  expiration_date?: string
  shelf_life_days?: number
  calories_per_serving?: number
  protein_g?: number
  carbs_g?: number
  fat_g?: number
  fiber_g?: number
  sugar_g?: number
  sodium_mg?: number
  serving_size?: string
  image_url?: string
  is_expired: boolean
  is_expiring_soon: boolean
  created_at: string
  updated_at: string
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
 * List pantry items that are expiring soon.
 */
export async function listExpiringSoon(userId: string): Promise<PantryItem[]> {
  const res = await fetch(`${API_BASE}/users/${userId}/pantry/expiring`)
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Failed to fetch expiring items')
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
