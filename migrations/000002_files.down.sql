-- Видалення тригера
DROP TRIGGER IF EXISTS update_files_updated_at ON files;

-- Видалення індексів
DROP INDEX IF EXISTS idx_files_user_id;
DROP INDEX IF EXISTS idx_files_booking_id;
DROP INDEX IF EXISTS idx_files_deleted_at;

-- Видалення таблиці
DROP TABLE IF EXISTS files; 