-- Drop users table and related objects
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_auth0_id;
DROP TABLE IF EXISTS users;
