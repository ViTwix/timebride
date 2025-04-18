CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255),
    phone VARCHAR(50),
    role VARCHAR(50) NOT NULL,
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    avatar_url VARCHAR(255),
    subscription_status VARCHAR(50) NOT NULL DEFAULT 'free',
    subscription_expires_at TIMESTAMP,
    
    -- White Label налаштування
    domain VARCHAR(255) UNIQUE,
    brand_colors JSONB DEFAULT '{"primary": "#000000", "secondary": "#ffffff"}',
    brand_logo_url VARCHAR(255),
    
    -- Налаштування календаря
    calendar_settings JSONB DEFAULT '{
        "default_view": "month",
        "working_hours": {"start": "09:00", "end": "18:00"},
        "working_days": [1,2,3,4,5],
        "slot_duration": "01:00"
    }',
    
    -- Загальні налаштування
    settings JSONB DEFAULT '{}',
    custom_fields JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    
    -- Основна інформація
    title VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    event_type VARCHAR(100) NOT NULL, -- wedding, portrait, commercial, etc.
    
    -- Часові межі
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    is_all_day BOOLEAN DEFAULT false,
    preparation_time INTERVAL,
    
    -- Інформація про клієнта
    client_name VARCHAR(255) NOT NULL,
    client_phone VARCHAR(50),
    client_email VARCHAR(255),
    client_instagram VARCHAR(255),
    
    -- Локація
    location VARCHAR(255),
    location_details TEXT,
    
    -- Фінансова інформація
    package_type VARCHAR(100),
    price DECIMAL(10,2),
    deposit DECIMAL(10,2),
    payment_status VARCHAR(50) DEFAULT 'pending',
    
    -- Команда
    team_members JSONB DEFAULT '[]',
    
    -- Матеріали
    delivery_link VARCHAR(255),
    delivery_password VARCHAR(100),
    delivery_expires_at TIMESTAMP,
    
    -- Google Calendar
    calendar_event_id VARCHAR(255),
    
    -- Додаткова інформація
    notes TEXT,
    custom_fields JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    
    name VARCHAR(255) NOT NULL,
    description TEXT,
    event_type VARCHAR(100),
    
    -- Базові налаштування
    duration INTERVAL NOT NULL DEFAULT '01:00:00',
    price DECIMAL(10,2),
    deposit DECIMAL(10,2),
    
    -- Шаблон полів
    fields_template JSONB NOT NULL DEFAULT '{}',
    
    -- Налаштування команди
    team_template JSONB DEFAULT '[]',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Індекси
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_domain ON users(domain);
CREATE INDEX idx_users_subscription ON users(subscription_status, subscription_expires_at);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_date ON bookings(start_time, end_time);
CREATE INDEX idx_bookings_client ON bookings(client_name, client_email);
CREATE INDEX idx_bookings_event_type ON bookings(event_type);

CREATE INDEX idx_templates_user ON templates(user_id);
