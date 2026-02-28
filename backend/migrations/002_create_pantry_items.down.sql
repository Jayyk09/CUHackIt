-- Drop pantry_items table and related objects
DROP TRIGGER IF EXISTS update_pantry_items_updated_at ON pantry_items;
DROP INDEX IF EXISTS idx_pantry_items_user_category;
DROP INDEX IF EXISTS idx_pantry_items_expiration;
DROP INDEX IF EXISTS idx_pantry_items_category;
DROP INDEX IF EXISTS idx_pantry_items_user_id;
DROP TABLE IF EXISTS pantry_items;
DROP TYPE IF EXISTS food_category;
