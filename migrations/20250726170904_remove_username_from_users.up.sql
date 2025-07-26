-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- 1. Удаляем индекс по username (он больше не нужен)
DROP INDEX IF EXISTS idx_users_username;

-- 2. Удаляем колонку username
ALTER TABLE users
DROP COLUMN IF EXISTS username;