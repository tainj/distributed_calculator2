DROP TRIGGER IF EXISTS update_examples_updated_at ON examples;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS examples;
DROP TABLE IF EXISTS users;

DROP FUNCTION IF EXISTS update_updated_at_column();