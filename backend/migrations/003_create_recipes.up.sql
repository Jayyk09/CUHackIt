-- Recipe difficulty enum
CREATE TYPE recipe_difficulty AS ENUM (
    'easy',
    'medium',
    'hard'
);

-- Recipe source enum (which agent generated it)
CREATE TYPE recipe_source AS ENUM (
    'pantry_only',      -- PantryOnlyAgent - uses only pantry items
    'flexible',         -- FlexibleRecipeAgent - pantry + up to 3 missing ingredients
    'user_created'      -- User manually created
);

-- Recipes table
CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Recipe info
    title VARCHAR(255) NOT NULL,
    description TEXT,
    cuisine VARCHAR(100),                    -- e.g., 'Italian', 'Asian', 'Mexican'
    
    -- Timing
    prep_time_minutes INTEGER,
    cook_time_minutes INTEGER,
    total_time_minutes INTEGER GENERATED ALWAYS AS (COALESCE(prep_time_minutes, 0) + COALESCE(cook_time_minutes, 0)) STORED,
    
    -- Servings
    servings INTEGER DEFAULT 2,
    
    -- Difficulty
    difficulty recipe_difficulty DEFAULT 'medium',
    
    -- Ingredients (stored as JSONB for flexibility)
    -- Format: [{"name": "chicken breast", "amount": "2", "unit": "lbs", "from_pantry": true}, ...]
    ingredients JSONB NOT NULL DEFAULT '[]',
    
    -- Missing ingredients (for flexible recipes)
    -- Format: [{"name": "heavy cream", "amount": "1", "unit": "cup"}, ...]
    missing_ingredients JSONB DEFAULT '[]',
    
    -- Instructions (stored as JSONB array)
    -- Format: ["Preheat oven to 375°F", "Season the chicken...", ...]
    instructions JSONB NOT NULL DEFAULT '[]',
    
    -- Nutrition per serving (estimated by AI)
    calories_per_serving DECIMAL(10, 2),
    protein_g DECIMAL(10, 2),
    carbs_g DECIMAL(10, 2),
    fat_g DECIMAL(10, 2),
    
    -- Source tracking
    source recipe_source NOT NULL DEFAULT 'pantry_only',
    
    -- AI generation metadata
    ai_model VARCHAR(100),                   -- e.g., 'gemini-1.5-flash'
    generation_prompt_hash VARCHAR(64),      -- Hash of the prompt used (for debugging)
    
    -- User interaction
    is_favorite BOOLEAN DEFAULT FALSE,
    times_cooked INTEGER DEFAULT 0,
    last_cooked_at TIMESTAMP WITH TIME ZONE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    notes TEXT,                              -- User's personal notes
    
    -- Tags for filtering
    tags TEXT[] DEFAULT '{}',                -- e.g., ARRAY['quick', 'healthy', 'comfort-food']
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for user's recipes
CREATE INDEX IF NOT EXISTS idx_recipes_user_id ON recipes(user_id);

-- Index for favorites
CREATE INDEX IF NOT EXISTS idx_recipes_favorites ON recipes(user_id, is_favorite) WHERE is_favorite = TRUE;

-- Index for cuisine filtering
CREATE INDEX IF NOT EXISTS idx_recipes_cuisine ON recipes(cuisine);

-- Index for source filtering
CREATE INDEX IF NOT EXISTS idx_recipes_source ON recipes(source);

-- GIN index for tags array search
CREATE INDEX IF NOT EXISTS idx_recipes_tags ON recipes USING GIN(tags);

-- GIN index for ingredients JSONB search
CREATE INDEX IF NOT EXISTS idx_recipes_ingredients ON recipes USING GIN(ingredients);

-- Trigger for recipes table
CREATE TRIGGER update_recipes_updated_at
    BEFORE UPDATE ON recipes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
