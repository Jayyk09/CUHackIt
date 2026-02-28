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
  /** True when these are hardcoded fallback recipes (backend was unreachable). */
  is_mock?: boolean
}

export async function generateRecipes(
  userId: string,
  mode: 'pantry_only' | 'flexible' | 'both' = 'flexible',
  count = 2,
  userPrompt = ''
): Promise<GenerateResult> {
  try {
    const res = await fetch(`${API_BASE}/users/${userId}/recipes/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ mode, recipe_count: count, user_prompt: userPrompt || undefined }),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      const msg = err.error || `Backend error (${res.status})`
      console.warn('[recipes-api] Backend returned', res.status, msg)

      // If the user typed a specific prompt, surface the error instead of
      // returning hardcoded mocks that ignore their request.
      if (userPrompt.trim()) {
        throw new Error(msg)
      }
      // Generic request — show sample recipes so the UI isn't empty.
      return getMockRecipes(mode)
    }
    return res.json()
  } catch (err) {
    // Re-throw intentional errors from above
    if (err instanceof Error && !err.message.includes('fetch')) {
      throw err
    }
    console.warn('[recipes-api] Network error', err)
    if (userPrompt.trim()) {
      throw new Error('Could not reach the recipe server — please try again')
    }
    return getMockRecipes(mode)
  }
}

/* ──────────────────────────────────────────────
 *  MOCK RECIPES — used when the backend is down
 *  or the pantry is empty.  Remove once the real
 *  pipeline is stable.
 * ────────────────────────────────────────────── */

const MOCK_RECIPES: GeneratedRecipe[] = [
  {
    title: 'Lemon-Herb Grilled Chicken with Spinach',
    description:
      'Quick pan-seared chicken breast finished with a bright lemon-garlic pan sauce, served over wilted baby spinach.',
    cuisine: 'Mediterranean',
    prep_time_minutes: 10,
    cook_time_minutes: 15,
    total_time_minutes: 25,
    servings: 2,
    difficulty: 'easy',
    ingredients: [
      { name: 'Chicken Breast', amount: '1', unit: 'lb', from_pantry: true },
      { name: 'Baby Spinach', amount: '5', unit: 'oz', from_pantry: true },
      { name: 'Lemon', amount: '1', unit: 'item', from_pantry: true },
      { name: 'Garlic', amount: '3', unit: 'cloves', from_pantry: true },
      { name: 'Olive Oil', amount: '2', unit: 'tbsp', from_pantry: true },
    ],
    instructions: [
      'Season chicken with salt, pepper, and a squeeze of lemon.',
      'Heat olive oil in a skillet over medium-high heat.',
      'Cook chicken 6-7 min per side until golden and cooked through.',
      'Remove chicken; add minced garlic and spinach to the same pan.',
      'Wilt spinach ~2 min, then squeeze remaining lemon over top.',
      'Slice chicken and serve over spinach.',
    ],
    missing_ingredients: [],
    calories_per_serving: 320,
    protein_g: 38,
    carbs_g: 5,
    fat_g: 16,
    tags: ['high-protein', 'low-carb', 'quick'],
    source: 'pantry_only',
  },
  {
    title: 'Salmon & Broccoli Rice Bowl',
    description:
      'Flaky pan-seared salmon over brown rice with steamed broccoli and a simple soy-lemon glaze.',
    cuisine: 'Asian-Fusion',
    prep_time_minutes: 10,
    cook_time_minutes: 20,
    total_time_minutes: 30,
    servings: 2,
    difficulty: 'easy',
    ingredients: [
      { name: 'Salmon Fillet', amount: '1', unit: 'lb', from_pantry: true },
      { name: 'Brown Rice', amount: '1', unit: 'cup', from_pantry: true },
      { name: 'Broccoli', amount: '1', unit: 'cup', from_pantry: true },
      { name: 'Soy Sauce', amount: '2', unit: 'tbsp', from_pantry: true },
      { name: 'Lemon', amount: '0.5', unit: 'item', from_pantry: true },
      { name: 'Garlic', amount: '2', unit: 'cloves', from_pantry: true },
    ],
    instructions: [
      'Cook brown rice according to package directions.',
      'Steam broccoli until bright green and tender-crisp, ~4 min.',
      'Season salmon with salt—sear skin-side down 4 min, flip and cook 3 min.',
      'Mix soy sauce, lemon juice, and minced garlic for a quick glaze.',
      'Plate rice, top with broccoli and salmon, drizzle glaze.',
    ],
    missing_ingredients: [],
    calories_per_serving: 480,
    protein_g: 34,
    carbs_g: 48,
    fat_g: 15,
    tags: ['omega-3', 'balanced', 'meal-prep'],
    source: 'pantry_only',
  },
  {
    title: 'Penne alla Pomodoro with Bell Peppers',
    description:
      'Classic Italian penne tossed in a quick fresh-tomato sauce with roasted bell peppers and a touch of cheddar.',
    cuisine: 'Italian',
    prep_time_minutes: 10,
    cook_time_minutes: 20,
    total_time_minutes: 30,
    servings: 3,
    difficulty: 'easy',
    ingredients: [
      { name: 'Pasta (Penne)', amount: '8', unit: 'oz', from_pantry: true },
      { name: 'Tomatoes', amount: '4', unit: 'item', from_pantry: true },
      { name: 'Bell Peppers', amount: '2', unit: 'item', from_pantry: true },
      { name: 'Garlic', amount: '3', unit: 'cloves', from_pantry: true },
      { name: 'Olive Oil', amount: '2', unit: 'tbsp', from_pantry: true },
      { name: 'Cheddar Cheese', amount: '2', unit: 'oz', from_pantry: true },
    ],
    instructions: [
      'Boil penne to al dente. Reserve ½ cup pasta water.',
      'Roast diced bell peppers in olive oil at 400 °F for 12 min.',
      'Sauté garlic 30 sec, add chopped tomatoes, simmer 8 min.',
      'Toss pasta and roasted peppers into the sauce; add pasta water to loosen.',
      'Top with shredded cheddar and fresh-cracked pepper.',
    ],
    missing_ingredients: [],
    calories_per_serving: 410,
    protein_g: 14,
    carbs_g: 58,
    fat_g: 14,
    tags: ['vegetarian', 'comfort-food'],
    source: 'pantry_only',
  },
  {
    title: 'Greek Yogurt Chicken Salad Wrap',
    description:
      'Creamy, high-protein chicken salad made with Greek yogurt instead of mayo, served in a lettuce wrap.',
    cuisine: 'American',
    prep_time_minutes: 15,
    cook_time_minutes: 0,
    total_time_minutes: 15,
    servings: 2,
    difficulty: 'easy',
    ingredients: [
      { name: 'Chicken Breast', amount: '1', unit: 'lb', from_pantry: true },
      { name: 'Greek Yogurt', amount: '0.5', unit: 'cup', from_pantry: true },
      { name: 'Lemon', amount: '0.5', unit: 'item', from_pantry: true },
      { name: 'Baby Spinach', amount: '2', unit: 'oz', from_pantry: true },
      { name: 'Garlic', amount: '1', unit: 'clove', from_pantry: true },
      { name: 'Tortillas or Lettuce Wraps', amount: '4', unit: 'item', from_pantry: false },
    ],
    instructions: [
      'Use leftover or poached chicken—shred with two forks.',
      'Mix yogurt, lemon juice, minced garlic, salt and pepper.',
      'Fold chicken into the yogurt mixture.',
      'Spoon onto large lettuce leaves or tortillas, top with spinach.',
    ],
    missing_ingredients: [
      { name: 'Tortillas or Lettuce Wraps', amount: '4', unit: 'item', from_pantry: false },
    ],
    calories_per_serving: 290,
    protein_g: 40,
    carbs_g: 8,
    fat_g: 10,
    tags: ['high-protein', 'no-cook', 'meal-prep'],
    source: 'flexible',
  },
  {
    title: 'Egg Fried Rice with Vegetables',
    description:
      'A fast weeknight fried rice loaded with eggs, bell peppers, and broccoli in a soy-garlic sauce.',
    cuisine: 'Chinese',
    prep_time_minutes: 5,
    cook_time_minutes: 12,
    total_time_minutes: 17,
    servings: 2,
    difficulty: 'easy',
    ingredients: [
      { name: 'Brown Rice', amount: '2', unit: 'cups cooked', from_pantry: true },
      { name: 'Eggs', amount: '3', unit: 'item', from_pantry: true },
      { name: 'Bell Peppers', amount: '1', unit: 'item', from_pantry: true },
      { name: 'Broccoli', amount: '0.5', unit: 'cup', from_pantry: true },
      { name: 'Soy Sauce', amount: '2', unit: 'tbsp', from_pantry: true },
      { name: 'Garlic', amount: '2', unit: 'cloves', from_pantry: true },
      { name: 'Olive Oil', amount: '1', unit: 'tbsp', from_pantry: true },
    ],
    instructions: [
      'Heat oil in a wok or large skillet over high heat.',
      'Scramble eggs, break into bits, set aside.',
      'Sauté diced peppers and broccoli 3 min.',
      'Add cold rice; stir-fry 3 min until slightly crispy.',
      'Return eggs, add soy sauce and minced garlic; toss 1 min.',
    ],
    missing_ingredients: [],
    calories_per_serving: 380,
    protein_g: 16,
    carbs_g: 50,
    fat_g: 12,
    tags: ['quick', 'budget-friendly', 'one-pan'],
    source: 'pantry_only',
  },
]

function getMockRecipes(
  mode: 'pantry_only' | 'flexible' | 'both'
): GenerateResult {
  const pantryOnly = MOCK_RECIPES.filter((r) => r.source === 'pantry_only')
  const flexible = MOCK_RECIPES.filter((r) => r.source === 'flexible')

  return {
    pantry_only_recipes: mode !== 'flexible' ? pantryOnly : undefined,
    flexible_recipes: mode !== 'pantry_only' ? flexible : undefined,
    all_recipes:
      mode === 'both'
        ? MOCK_RECIPES
        : mode === 'pantry_only'
          ? pantryOnly
          : [...pantryOnly, ...flexible],
    generated_at: new Date().toISOString(),
    total_count: MOCK_RECIPES.length,
    filtered_count: MOCK_RECIPES.length,
    is_mock: true,
  }
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
