-- Додаємо розширення uuid-ossp для генерації UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Видаляємо старі таблиці
DROP TABLE IF EXISTS files CASCADE;
DROP TABLE IF EXISTS booking_team_members CASCADE;
DROP TABLE IF EXISTS additional_services CASCADE;
DROP TABLE IF EXISTS service_packages CASCADE;
DROP TABLE IF EXISTS service_categories CASCADE;
DROP TABLE IF EXISTS templates CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Створюємо нову таблицю users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT,
    auth_provider TEXT NOT NULL DEFAULT 'local', -- 'local', 'google', 'apple', 'facebook'
    provider_id TEXT,
    name TEXT NOT NULL,
    avatar_url TEXT,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_2fa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    two_fa_secret TEXT,
    phone TEXT,
    
    -- Поля для роботи з командою
    parent_id UUID REFERENCES users(id) ON DELETE SET NULL, -- ID адміністратора, який створив цього користувача
    role TEXT NOT NULL DEFAULT 'member', -- 'owner', 'admin', 'photographer', 'videographer', 'editor', 'retoucher', 'assistant', 'member'
    permissions JSONB DEFAULT '{"finances": false, "clients": true, "bookings": true, "templates": false, "team": false}', -- Права доступу
    
    language TEXT NOT NULL DEFAULT 'uk',
    settings JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Додаткові обмеження
    CONSTRAINT users_auth_provider_check 
        CHECK (auth_provider IN ('local', 'google', 'apple', 'facebook')),
    CONSTRAINT users_role_check 
        CHECK (role IN ('owner', 'admin', 'photographer', 'videographer', 'editor', 'retoucher', 'assistant', 'member'))
);

-- Додаємо обмеження для локальної та OAuth автентифікації після створення таблиці
ALTER TABLE users
ADD CONSTRAINT users_local_auth_password_check 
    CHECK ((auth_provider = 'local' AND password_hash IS NOT NULL) 
        OR auth_provider != 'local');

ALTER TABLE users
ADD CONSTRAINT users_oauth_provider_id_check 
    CHECK ((auth_provider != 'local' AND provider_id IS NOT NULL) 
        OR auth_provider = 'local');

-- Додаємо обмеження для перевірки parent_id
ALTER TABLE users
ADD CONSTRAINT users_parent_id_check
    CHECK (parent_id != id); -- Користувач не може бути своїм власним керівником

-- Додаємо обмеження на ієрархію ролей (тільки owner/admin можуть бути керівниками)
CREATE OR REPLACE FUNCTION check_parent_role() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.parent_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM users 
            WHERE id = NEW.parent_id AND (role = 'owner' OR role = 'admin')
        ) THEN
            RAISE EXCEPTION 'Parent user must have role owner or admin';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_parent_role_trigger
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION check_parent_role();

-- Створюємо нову таблицю bookings
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    event_date DATE NOT NULL,
    location TEXT,
    client_name TEXT NOT NULL,
    client_phone TEXT,
    instagram_handle TEXT,
    contract_url TEXT,
    total_price NUMERIC NOT NULL DEFAULT 0,
    deposit NUMERIC NOT NULL DEFAULT 0,
    payment_status TEXT NOT NULL DEFAULT 'pending', -- 'pending', 'paid', 'partial', 'after_delivery'
    status TEXT NOT NULL DEFAULT 'active', -- 'active', 'completed', 'archived', 'cancelled'
    notes TEXT,
    
    -- Додаємо поле для учасників команди, які беруть участь у цьому бронюванні
    team_members JSONB DEFAULT '[]', -- Масив user_id учасників команди та їх ролей у проекті
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Додаткові обмеження
    CONSTRAINT bookings_payment_status_check 
        CHECK (payment_status IN ('pending', 'paid', 'partial', 'after_delivery')),
    CONSTRAINT bookings_status_check 
        CHECK (status IN ('active', 'completed', 'archived', 'cancelled'))
);

-- Індекси для таблиці users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_auth_provider_provider_id ON users(auth_provider, provider_id);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_parent_id ON users(parent_id); -- Додаємо індекс для пошуку членів команди

-- Індекси для таблиці bookings
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_event_date ON bookings(event_date);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_payment_status ON bookings(payment_status);

-- Тригер для автоматичного оновлення updated_at
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Тригери для оновлення поля updated_at
CREATE TRIGGER set_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_bookings
BEFORE UPDATE ON bookings
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp(); 