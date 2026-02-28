-- Food categories enum
CREATE TYPE food_category AS ENUM (
    'PRODUCE',
    'DAIRY',
    'MEAT',
    'SEAFOOD',
    'PANTRY',
    'FROZEN',
    'BAKERY',
    'SNACKS',
    'BEVERAGE',
    'DELI',
    'SPECIALTY'
);

-- Pantry items table
CREATE TABLE IF NOT EXISTS pantry_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Food identification
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255),
    barcode VARCHAR(100),                    -- UPC/EAN barcode
    open_food_facts_id VARCHAR(100),         -- Open Food Facts product ID
    
    -- Categorization
    category food_category NOT NULL,
    
    -- Quantity tracking
    quantity DECIMAL(10, 2) DEFAULT 1,
    unit VARCHAR(50) DEFAULT 'item',         -- 'item', 'oz', 'lb', 'g', 'kg', 'ml', 'l', 'cup', etc.
    
    -- Expiration tracking
    purchase_date DATE DEFAULT CURRENT_DATE,
    expiration_date DATE,
    shelf_life_days INTEGER,                 -- Estimated days until expiration from purchase
    
    -- Nutrition (from Open Food Facts or manual entry)
    calories_per_serving DECIMAL(10, 2),
    protein_g DECIMAL(10, 2),
    carbs_g DECIMAL(10, 2),
    fat_g DECIMAL(10, 2),
    fiber_g DECIMAL(10, 2),
    sugar_g DECIMAL(10, 2),
    sodium_mg DECIMAL(10, 2),
    serving_size VARCHAR(50),
    
    -- Image
    image_url TEXT,
    
    -- Status (computed by application based on expiration_date)
    is_expired BOOLEAN NOT NULL DEFAULT FALSE,
    is_expiring_soon BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for user's pantry lookups
CREATE INDEX IF NOT EXISTS idx_pantry_items_user_id ON pantry_items(user_id);

-- Index for category filtering
CREATE INDEX IF NOT EXISTS idx_pantry_items_category ON pantry_items(category);

-- Index for expiration queries
CREATE INDEX IF NOT EXISTS idx_pantry_items_expiration ON pantry_items(expiration_date);

-- Composite index for user + category
CREATE INDEX IF NOT EXISTS idx_pantry_items_user_category ON pantry_items(user_id, category);

-- Trigger for pantry_items table
CREATE TRIGGER update_pantry_items_updated_at
    BEFORE UPDATE ON pantry_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
