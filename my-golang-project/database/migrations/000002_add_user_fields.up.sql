-- Добавляем новые колонки в таблицу users
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS email varchar(255) UNIQUE,
ADD COLUMN IF NOT EXISTS age int,
ADD COLUMN IF NOT EXISTS created_at timestamp DEFAULT now();

-- Обновляем существующие записи (если нужно)
UPDATE users SET email = 'john.doe@example.com' WHERE id = 1 AND email IS NULL;
UPDATE users SET email = 'jane.smith@example.com' WHERE id = 2 AND email IS NULL;
UPDATE users SET age = 30 WHERE id = 1 AND age IS NULL;
UPDATE users SET age = 25 WHERE id = 2 AND age IS NULL;