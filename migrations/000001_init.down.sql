DROP TRIGGER IF EXISTS calculate_booking_finances ON bookings;
DROP TRIGGER IF EXISTS update_users_timestamp ON users;
DROP TRIGGER IF EXISTS update_service_categories_timestamp ON service_categories;
DROP TRIGGER IF EXISTS update_service_packages_timestamp ON service_packages;
DROP TRIGGER IF EXISTS update_additional_services_timestamp ON additional_services;
DROP TRIGGER IF EXISTS update_bookings_timestamp ON bookings;
DROP TRIGGER IF EXISTS update_booking_team_members_timestamp ON booking_team_members;

DROP FUNCTION IF EXISTS update_booking_calculations();
DROP FUNCTION IF EXISTS update_timestamp();

DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS booking_team_members;
DROP TABLE IF EXISTS additional_services;
DROP TABLE IF EXISTS service_packages;
DROP TABLE IF EXISTS service_categories;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
