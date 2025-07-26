-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- Так как таблица была пустая, просто добавляем колонку назад
ALTER TABLE users
ADD COLUMN username VARCHAR(50) NOT NULL UNIQUE;