const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

export interface Ingredient {
  name: string
  amount: string
  unit?: string
  from_pantry: boolean
}

export interface GeneratedRecipe {
  title: string
  description: string
  cuisine?: string
  prep_time_minutes?: number
  cook_time_minutes?: number
  total_time_minutes?: number
  servings?: number
  difficulty?: 'easy' | 'medium' | 'hard'
  ingredients: Ingredient[]
  instructions: string[]
  missing_ingredients?: Ingredient[]
  calories_per_serving?: number
  protein_g?: number
  carbs_g?: number
  fat_g?: number
  tags?: string[]
  source: 'pantry_only' | 'flexible' | 'user_created'
}

export interface SavedRecipe {
  id: string
  user_id: string
  title: string
  description?: string
  cuisine?: string
  prep_time_minutes?: number
  cook_time_minutes?: number
  total_time_minutes?: number
  servings: number
  difficulty: 'easy' | 'medium' | 'hard'
  ingredients: Ingredient[]
  missing_ingredients?: Ingredient[]
  instructions: string[]
  calories_per_serving?: number
  protein_g?: number
  carbs_g?: number
  fat_g?: number
  source: 'pantry_only' | 'flexible' | 'user_created'
  ai_model?: string
  is_favorite: boolean
  times_cooked: number
  last_cooked_at?: string
  rating?: number
  notes?: string
  tags: string[]
  created_at: string
  updated_at: string
}

export interface GenerateResult {
  pantry_only_recipes?: GeneratedRecipe[]
  flexible_recipes?: GeneratedRecipe[]
  all_recipes: GeneratedRecipe[]
  generated_at: string
  total_count: number
  filtered_count: number
}

export async function generateRecipes(
  userId: string,
  mode: 'pantry_only' | 'flexible' | 'both' = 'flexible',
  count = 2
): Promise<GenerateResult> {
  const res = await fetch(`${API_BASE}/users/${userId}/recipes/generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ mode, recipe_count: count }),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Failed to generate recipes')
  }
  return res.json()
}

export async function saveRecipe(userId: string, recipe: GeneratedRecipe): Promise<SavedRecipe> {
  const res = await fetch(`${API_BASE}/users/${userId}/recipes`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      title: recipe.title,
      description: recipe.description ?? '',
      cuisine: recipe.cuisine ?? '',
      prep_time_minutes: recipe.prep_time_minutes,
      cook_time_minutes: recipe.cook_time_minutes,
      servings: recipe.servings ?? 2,
      difficulty: recipe.difficulty ?? 'medium',
      ingredients: recipe.ingredients,
      missing_ingredients: recipe.missing_ingredients ?? [],
      instructions: recipe.instructions,
      calories_per_serving: recipe.calories_per_serving,
      protein_g: recipe.protein_g,
      carbs_g: recipe.carbs_g,
      fat_g: recipe.fat_g,
      source: recipe.source,
      tags: recipe.tags ?? [],
    }),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Failed to save recipe')
  }
  return res.json()
}

export async function listRecipes(userId: string): Promise<SavedRecipe[]> {
  const res = await fetch(`${API_BASE}/users/${userId}/recipes`)
  if (!res.ok) throw new Error('Failed to fetch recipes')
  const data = await res.json()
  return Array.isArray(data) ? data : []
}

export async function toggleFavorite(userId: string, recipeId: string): Promise<SavedRecipe> {
  const res = await fetch(`${API_BASE}/users/${userId}/recipes/${recipeId}/favorite`, {
    method: 'POST',
  })
  if (!res.ok) throw new Error('Failed to toggle favorite')
  return res.json()
}
