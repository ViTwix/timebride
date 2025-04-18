CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Основна інформація
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255),
    phone VARCHAR(50),
    
    -- Роль та доступи
    role VARCHAR(50) NOT NULL,  -- master_admin, admin, user
    parent_admin_id UUID REFERENCES users(id),  -- для user ролі, посилання на їх admin
    
    -- Налаштування локалізації
    language VARCHAR(10) NOT NULL DEFAULT 'en',  -- uk, en
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    
    -- Налаштування інтерфейсу
    theme VARCHAR(20) NOT NULL DEFAULT 'light',  -- light, dark, custom
    theme_settings JSONB DEFAULT '{}',  -- для custom теми
    
    -- Підписка та інтеграції
    subscription_plan VARCHAR(50) NOT NULL DEFAULT 'free',
    subscription_expires_at TIMESTAMP,
    google_calendar_token TEXT,
    telegram_chat_id VARCHAR(100),
    
    -- OAuth провайдери
    oauth_providers JSONB DEFAULT '{}',  -- {"google": {"id": "...", "email": "..."}, "apple": {...}}
    
    -- White Label налаштування
    domain VARCHAR(255) UNIQUE,
    brand_colors JSONB DEFAULT '{
        "primary": "#000000",
        "secondary": "#ffffff",
        "accent": "#000000",
        "background": "#ffffff",
        "text": "#000000"
    }',
    brand_logo_url VARCHAR(255),
    
    -- Налаштування календаря
    calendar_settings JSONB DEFAULT '{
        "default_view": "month",
        "first_day_of_week": 1
    }',
    
    -- Налаштування команди
    team_settings JSONB DEFAULT '{
        "roles": [],
        "permissions": {}
    }',
    
    -- Загальні налаштування
    settings JSONB DEFAULT '{
        "notifications": {
            "email": true,
            "telegram": false
        },
        "deadlines": {
            "default_delivery_time": "30d"
        },
        "templates": {}
    }',
    
    -- Кастомні поля
    custom_fields JSONB DEFAULT '{}',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    
    -- Основна інформація
    title VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL DEFAULT 'wedding',
    status VARCHAR(50) NOT NULL DEFAULT 'booked',
    
    -- Часові межі
    date DATE NOT NULL,
    start_time TIME,
    end_time TIME,
    deadline TIMESTAMP,
    
    -- Локація
    location TEXT,
    location_details TEXT,
    
    -- Клієнт
    client_name VARCHAR(255),
    client_phone VARCHAR(50),
    client_instagram VARCHAR(255),
    client_notes TEXT,
    
    -- Фінансова інформація
    total_price DECIMAL(10,2) DEFAULT 0,
    deposit DECIMAL(10,2) DEFAULT 0,
    shooting_day_payment DECIMAL(10,2) DEFAULT 0,
    final_payment DECIMAL(10,2) DEFAULT 0,
    
    -- Додаткові витрати
    travel_expenses DECIMAL(10,2) DEFAULT 0,
    rental_expenses DECIMAL(10,2) DEFAULT 0,
    other_expenses DECIMAL(10,2) DEFAULT 0,
    
    -- Розрахункові поля
    remaining_payment DECIMAL(10,2) DEFAULT 0,
    net_profit DECIMAL(10,2) DEFAULT 0,
    
    -- Матеріали
    raw_footage_drive VARCHAR(255),
    delivery_link VARCHAR(255),
    contract_link VARCHAR(255),
    
    -- Інтеграції
    calendar_event_id VARCHAR(255),
    
    -- Додаткові налаштування
    settings JSONB DEFAULT '{
        "deadline_template": "6m",
        "notifications": {
            "email": true,
            "telegram": false
        }
    }',
    
    -- Історія платежів
    payments_history JSONB DEFAULT '[]',
    
    -- Додаємо нові поля для прайс-листів
    package_id UUID REFERENCES service_packages(id),
    package_snapshot JSONB DEFAULT NULL,
    selected_services JSONB DEFAULT '[]',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE booking_team_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID REFERENCES bookings(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    
    role VARCHAR(100) NOT NULL,
    name VARCHAR(255),
    payment DECIMAL(10,2) DEFAULT 0,
    payment_status VARCHAR(50) DEFAULT 'pending',
    
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

-- Додаємо нові таблиці для прайс-листів перед bookings
CREATE TABLE service_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    
    settings JSONB DEFAULT '{
        "icon": null,
        "color": null,
        "show_in_menu": true
    }',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, slug)
);

CREATE TABLE service_packages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    category_id UUID REFERENCES service_categories(id),
    
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    
    -- Базові налаштування
    settings JSONB DEFAULT '{
        "duration_hours": null,
        "team_size": null,
        "cameras_count": null,
        "location_type": "any",
        "highlight": false,
        "featured": false
    }',
    
    -- Цінова політика
    pricing JSONB DEFAULT '{
        "base_price": 0,
        "currency": "USD",
        "deposit": {
            "type": "percentage",
            "value": 30
        },
        "payment_schedule": [
            {"stage": "booking", "type": "deposit"},
            {"stage": "shooting_day", "type": "remaining"}
        ],
        "discounts": [],
        "seasonal_adjustments": []
    }',
    
    -- Що входить у пакет
    included_items JSONB DEFAULT '{
        "services": [],
        "deliverables": {
            "main_video_duration": "30-90",
            "highlights_duration": "3-5"
        },
        "equipment": []
    }',
    
    -- Команда за замовчуванням
    team_template JSONB DEFAULT '[]',
    
    -- Додаткові послуги
    additional_services JSONB DEFAULT '{
        "available": [],
        "required": []
    }',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, slug)
);

CREATE TABLE additional_services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    
    settings JSONB DEFAULT '{
        "duration": null,
        "type": "addon"
    }',
    
    pricing JSONB DEFAULT '{
        "base_price": 0,
        "pricing_type": "fixed"
    }',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, slug)
);

-- Індекси
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_parent_admin ON users(parent_admin_id);
CREATE INDEX idx_users_domain ON users(domain);
CREATE INDEX idx_users_subscription ON users(subscription_plan, subscription_expires_at);

CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_date ON bookings(date);
CREATE INDEX idx_bookings_event_type ON bookings(event_type);
CREATE INDEX idx_bookings_client ON bookings(client_name, client_phone);

CREATE INDEX idx_booking_team_booking ON booking_team_members(booking_id);
CREATE INDEX idx_booking_team_user ON booking_team_members(user_id);
CREATE INDEX idx_booking_team_role ON booking_team_members(role);

CREATE INDEX idx_templates_user ON templates(user_id);

-- Оновлюємо тригер для розрахунків з урахуванням команди
CREATE OR REPLACE FUNCTION update_booking_calculations()
RETURNS TRIGGER AS $$
DECLARE
    team_expenses DECIMAL(10,2);
    total_expenses DECIMAL(10,2);
BEGIN
    -- Отримуємо суму витрат на команду
    SELECT COALESCE(SUM(payment), 0)
    INTO team_expenses
    FROM booking_team_members
    WHERE booking_id = NEW.id;
    
    -- Сума всіх витрат
    total_expenses := team_expenses + 
                     COALESCE(NEW.travel_expenses, 0) + 
                     COALESCE(NEW.rental_expenses, 0) + 
                     COALESCE(NEW.other_expenses, 0);
    
    -- Розрахунок remaining_payment
    NEW.remaining_payment = NEW.total_price - 
                          COALESCE(NEW.deposit, 0) - 
                          COALESCE(NEW.shooting_day_payment, 0) - 
                          COALESCE(NEW.final_payment, 0);
    
    -- Розрахунок net_profit
    NEW.net_profit = NEW.total_price - total_expenses;
    
    -- Перевірка на від'ємні значення
    IF NEW.remaining_payment < 0 THEN
        NEW.remaining_payment = 0;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER calculate_booking_finances
    BEFORE INSERT OR UPDATE ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION update_booking_calculations();

-- Додаємо нові індекси
CREATE INDEX idx_service_categories_user ON service_categories(user_id);
CREATE INDEX idx_service_categories_active ON service_categories(is_active);

CREATE INDEX idx_service_packages_user ON service_packages(user_id);
CREATE INDEX idx_service_packages_category ON service_packages(category_id);
CREATE INDEX idx_service_packages_active ON service_packages(is_active);

CREATE INDEX idx_additional_services_user ON additional_services(user_id);
CREATE INDEX idx_additional_services_active ON additional_services(is_active);
