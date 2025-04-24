-- Drop triggers
DROP TRIGGER IF EXISTS update_bookings_updated_at ON bookings;
DROP TRIGGER IF EXISTS update_price_templates_updated_at ON price_templates;
DROP TRIGGER IF EXISTS update_team_members_updated_at ON team_members;
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in correct order (respect foreign keys)
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS price_templates CASCADE;
DROP TABLE IF EXISTS team_members CASCADE;
DROP TABLE IF EXISTS clients CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp"; 