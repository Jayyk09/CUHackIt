-- Drop recipes table and related objects
DROP TRIGGER IF EXISTS update_recipes_updated_at ON recipes;
DROP INDEX IF EXISTS idx_recipes_ingredients;
DROP INDEX IF EXISTS idx_recipes_tags;
DROP INDEX IF EXISTS idx_recipes_source;
DROP INDEX IF EXISTS idx_recipes_cuisine;
DROP INDEX IF EXISTS idx_recipes_favorites;
DROP INDEX IF EXISTS idx_recipes_user_id;
DROP TABLE IF EXISTS recipes;
DROP TYPE IF EXISTS recipe_source;
DROP TYPE IF EXISTS recipe_difficulty;
