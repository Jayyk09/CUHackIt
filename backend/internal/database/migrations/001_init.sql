-- 001_init.sql
-- Creates the initial tables for Sift

CREATE TABLE IF NOT EXISTS account (
    account_id SERIAL PRIMARY KEY,
    labels TEXT[] DEFAULT '{}',
    allergens TEXT[] DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS food (
    id SERIAL PRIMARY KEY,
    product_name TEXT NOT NULL,
    environmental_score REAL NOT NULL,
    nutriscore_score REAL NOT NULL,
    labels_en TEXT[] DEFAULT '{}',
    allergens_en TEXT[] DEFAULT '{}',
    traces_en TEXT[] DEFAULT '{}',
    image_url TEXT,
    image_small_url TEXT,
    shelf_life INT,
    category TEXT
);

CREATE TABLE IF NOT EXISTS pantry (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES account(account_id) ON DELETE CASCADE,
    food_id INT NOT NULL REFERENCES food(id) ON DELETE CASCADE,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    quantity INT NOT NULL DEFAULT 1,
    units TEXT NOT NULL DEFAULT 'unit',
    category TEXT,
    is_frozen BOOLEAN DEFAULT FALSE
);
