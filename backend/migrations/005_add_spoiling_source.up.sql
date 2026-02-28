-- Add 'spoiling' to the recipe_source enum
ALTER TYPE recipe_source ADD VALUE IF NOT EXISTS 'spoiling';
