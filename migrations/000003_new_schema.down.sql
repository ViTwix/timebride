-- Видаляємо тригери
DROP TRIGGER IF EXISTS set_timestamp_bookings ON bookings;
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP TRIGGER IF EXISTS check_parent_role_trigger ON users;

-- Видаляємо функції тригерів
DROP FUNCTION IF EXISTS trigger_set_timestamp();
DROP FUNCTION IF EXISTS check_parent_role();

-- Видаляємо таблиці
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS users CASCADE; 