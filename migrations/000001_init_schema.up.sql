-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table (підрядники)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    company_name VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    language VARCHAR(10) NOT NULL DEFAULT 'UA',
    theme VARCHAR(20) NOT NULL DEFAULT 'light',
    default_currency VARCHAR(10) NOT NULL DEFAULT 'UAH',
    subscription_plan VARCHAR(50) NOT NULL DEFAULT 'free',
    subscription_expires_at TIMESTAMP WITH TIME ZONE,
    google_calendar_token TEXT,
    telegram_chat_id VARCHAR(100),
    oauth_providers JSONB DEFAULT '{}',
    storage_limit_gb INTEGER NOT NULL DEFAULT 100,
    storage_used_gb FLOAT NOT NULL DEFAULT 0,
    gallery_template VARCHAR(100),
    gallery_logo_url VARCHAR(255),
    gallery_social_links JSONB DEFAULT '[]',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Clients table (клієнти підрядників)
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(255),
    notes TEXT,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Team Members table (команда підрядника)
CREATE TABLE team_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    role VARCHAR(50) NOT NULL,
    permissions JSONB DEFAULT '{}',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Price Templates table (прайс-листи)
CREATE TABLE price_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'UAH',
    price DECIMAL(10,2) NOT NULL,
    deposit DECIMAL(10,2) NOT NULL,
    description TEXT,
    team_payments JSONB DEFAULT '[]',
    deadline_days INTEGER NOT NULL DEFAULT 180,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Bookings table (бронювання)
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES clients(id),
    title VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL,
    location VARCHAR(255),
    package_name VARCHAR(255),
    currency VARCHAR(10) NOT NULL DEFAULT 'UAH',
    price_total DECIMAL(10,2) NOT NULL,
    price_prepayment DECIMAL(10,2) NOT NULL,
    price_extra DECIMAL(10,2) DEFAULT 0,
    team_payments JSONB DEFAULT '[]',
    price_profit DECIMAL(10,2) GENERATED ALWAYS AS (
        price_total - price_prepayment - price_extra - (
            SELECT COALESCE(SUM((payment->>'amount')::DECIMAL), 0) 
            FROM jsonb_array_elements(team_payments) AS payment
        )
    ) STORED,
    price_left_to_pay DECIMAL(10,2) GENERATED ALWAYS AS (
        price_total - price_prepayment
    ) STORED,
    deadline_days INTEGER NOT NULL DEFAULT 180,
    contract_file_url VARCHAR(255),
    delivery_page_url VARCHAR(255),
    disk_code VARCHAR(100),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    booking_others JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE INDEX idx_clients_user_id ON clients(user_id);
CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_phone ON clients(phone);
CREATE INDEX idx_clients_deleted_at ON clients(deleted_at);

CREATE INDEX idx_team_members_user_id ON team_members(user_id);
CREATE INDEX idx_team_members_email ON team_members(email);
CREATE INDEX idx_team_members_role ON team_members(role);
CREATE INDEX idx_team_members_deleted_at ON team_members(deleted_at);

CREATE INDEX idx_price_templates_user_id ON price_templates(user_id);
CREATE INDEX idx_price_templates_event_type ON price_templates(event_type);
CREATE INDEX idx_price_templates_deleted_at ON price_templates(deleted_at);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_client_id ON bookings(client_id);
CREATE INDEX idx_bookings_event_date ON bookings(event_date);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_event_type ON bookings(event_type);
CREATE INDEX idx_bookings_deleted_at ON bookings(deleted_at);

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_clients_updated_at
    BEFORE UPDATE ON clients
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_team_members_updated_at
    BEFORE UPDATE ON team_members
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_price_templates_updated_at
    BEFORE UPDATE ON price_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookings_updated_at
    BEFORE UPDATE ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 