-- Users table with profile fields for personalization
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auth0_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    
    -- Profile fields (collected during onboarding)
    allergens TEXT[] DEFAULT '{}',           -- e.g., ARRAY['peanuts', 'shellfish', 'dairy']
    dietary_preferences TEXT[] DEFAULT '{}', -- e.g., ARRAY['vegetarian', 'low-carb', 'keto']
    nutritional_goals TEXT[] DEFAULT '{}',   -- e.g., ARRAY['high-protein', 'low-sodium', 'low-sugar']
    cooking_skill VARCHAR(50) DEFAULT 'beginner', -- 'beginner', 'intermediate', 'advanced'
    cuisine_preferences TEXT[] DEFAULT '{}', -- e.g., ARRAY['italian', 'asian', 'mexican', 'mediterranean']
    
    -- Onboarding status
    onboarding_completed BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for Auth0 lookups
CREATE INDEX IF NOT EXISTS idx_users_auth0_id ON users(auth0_id);

-- Index for email lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Function to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
