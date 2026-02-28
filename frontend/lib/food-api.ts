export interface FoodProduct {
  id: number
  product_name: string
  norm_environmental_score: number | null
  norm_nutriscore: string | null
  labels_en: string[] | null
  allergens_en: string[] | null
  traces_en: string[] | null
  image_url: string | null
  image_small_url: string | null
  shelf_life: number | null
  category: string | null
}

export interface FoodItem {
  id: number
  product_name: string
  category: string
  environmental_score: number
  nutriscore_score: number
  labels_en: string[]
  allergens_en: string[]
  traces_en: string[]
  shelf_life: number
  is_spoiled: boolean
  image_url: string
}

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'
const FALLBACK_IMAGE_URL = '/placeholder.jpg'

function toNumber(value: string | number | null): number {
  if (value === null) return 0
  const parsed = typeof value === 'number' ? value : Number.parseFloat(value)
  return Number.isNaN(parsed) ? 0 : parsed
}

export function mapProductToFoodItem(product: FoodProduct): FoodItem {
  return {
    id: product.id,
    product_name: product.product_name,
    category: product.category?.trim() || 'specialty',
    environmental_score: Math.round(toNumber(product.norm_environmental_score)),
    nutriscore_score: Math.round(toNumber(product.norm_nutriscore)),
    labels_en: product.labels_en ?? [],
    allergens_en: product.allergens_en ?? [],
    traces_en: product.traces_en ?? [],
    shelf_life: product.shelf_life ?? 0,
    is_spoiled: false,
    image_url: product.image_url ?? product.image_small_url ?? FALLBACK_IMAGE_URL,
  }
}

export async function searchFoodProducts({
  search,
  limit,
  offset,
  signal,
}: {
  search: string
  limit: number
  offset: number
  signal?: AbortSignal
}): Promise<FoodItem[]> {
  const params = new URLSearchParams()
  params.set('limit', String(limit))
  params.set('offset', String(offset))
  params.set('search', search)
  params.set('q', search)

  const res = await fetch(`${API_URL}/food?${params.toString()}`, {
    method: 'GET',
    credentials: 'include',
    signal,
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to fetch food results')
  }

  const data = (await res.json()) as FoodProduct[]
  if (!Array.isArray(data)) {
    return []
  }

  return data.map(mapProductToFoodItem)
}
