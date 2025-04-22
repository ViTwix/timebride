# Вся комунікація і коментарі українською мовою.
# Ultrathink during each request processing

# Увесь інтерфейс сайту — як адмінка, так і клієнтська частина — повинен бути побудований виключно на базі Tabler UI (https://tabler.io) https://github.com/tabler/tabler
Врахуй такі ключові вимоги:
Мінімалістичний, легкий для сприйняття інтерфейс
Повна адаптивність
Таблиці, форми, календарі, модалки та інші елементи — лише ті, що є в Tabler UI
Якщо потрібно щось кастомізувати — роби це в стилістиці Tabler
Всі компоненти мають мати єдину візуальну мову згідно Tabler
На виході я очікую HTML+CSS/JS компоненти або SSR-ready шаблони, повністю сумісні з Tabler UI. Якщо використовуєш JS — Alpine.js.
Якщо щось потрібно придумати — інтерпретуй це через логіку Tabler UI.

## 📋 Опис проекту
TimeBride - це SaaS CRM-платформа для фотографів, відеографів, та в майбутньому весільних агенцій, що фокусується на:
- Легкому управлінні зйомками (бронюваннями) та процесом віддачі матеріалу
- Високій швидкодії та оптимізованій роботі
- Синхронізації всіх бронювань з Google Calendar, iCal, Google Sheets та в майбутньому можливо інших сервісів
- Командній роботі з різними рівнями доступу

## 🏗 Оптимізована структура проекту
```
.
├── cmd/
│   └── app/
│       └── main.go             # Точка входу в програму
├── internal/
│   ├── config/                 # Конфігурація програми
│   ├── models/                 # Моделі даних (user, booking, template, file)
│   ├── handlers/               # HTTP обробники
│   ├── services/               # Бізнес-логіка
│   ├── repositories/           # Робота з БД
│   ├── middleware/             # Middleware компоненти
│   ├── database/               # Підключення до БД
│   └── utils/                  # Загальні утиліти
├── migrations/                 # SQL міграції
├── pkg/                        # Публічні пакети
│   └── logger/                 # Логування
├── web/
│   ├── templates/              # HTML шаблони
│   │   └── layout.html         # Основний шаблон сайту
│   │   └── booking/            # Шаблони для бронювань
│   ├── public/                 # Публічні статичні файли
│   │   ├── css/
│   │   │   ├── tabler.min.css  # Основний стиль Tabler
│   │   │   ├── tabler-icons.min.css  # Іконки Tabler
│   │   │   └── main.css        # Єдиний CSS файл проекту
│   │   ├── js/                 # JavaScript файли
│   │   │   ├── tabler.min.js   # Основний JS файл Tabler
│   │   │   └── main.js         # Головний JS файл проекту
│   │   ├── img/                # Зображення
│   │   └── fonts/              # Шрифти (включно з Tabler Icons)
├── config.yaml                 # Основна конфігурація
├── docker-compose.yml          # Docker конфігурація
├── Dockerfile                  # Dockerfile для збірки
└── Makefile                    # Команди для розробки
```

## 📊 Оновлена структура бази даних

### Users
```sql
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
    role TEXT NOT NULL DEFAULT 'user', -- 'owner', 'admin', 'retoucher', 'editor', 'operator', 'assistant', 'user'
    permissions JSONB DEFAULT '{"view_financials": false, "edit_projects": false, "view_projects": true, "manage_team": false}', -- Права доступу
    
    country TEXT,
    city TEXT,
    language TEXT NOT NULL DEFAULT 'uk',
    settings JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Bookings
```sql
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    event_type TEXT NOT NULL DEFAULT 'wedding', -- 'wedding', 'portrait', 'family', 'event', 'commercial'
    event_date DATE NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    location TEXT,
    client_name TEXT NOT NULL,
    client_phone TEXT,
    client_email TEXT NOT NULL,
    instagram_handle TEXT,
    contract_url TEXT,
    total_price NUMERIC NOT NULL DEFAULT 0,
    deposit NUMERIC NOT NULL DEFAULT 0,
    payment_status TEXT NOT NULL DEFAULT 'pending', -- 'pending', 'paid', 'partial', 'after_delivery'
    status TEXT NOT NULL DEFAULT 'active', -- 'active', 'completed', 'archived', 'cancelled', 'pending'
    notes TEXT,
    custom_fields JSONB DEFAULT '{}',
    team_members JSONB DEFAULT '[]', -- Масив user_id учасників команди та їх ролей у проекті
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## 🛠 API Endpoints (Реалізовані)

### Аутентифікація
**Локальна:**
- POST `/api/auth/register` - Реєстрація нового користувача
- POST `/api/auth/login` - Вхід за допомогою email та пароля
- POST `/api/auth/refresh` - Оновлення токену

**OAuth:**
- GET `/api/auth/oauth/:provider` - Початок OAuth авторизації (Google, Facebook, Apple)
- GET `/api/auth/oauth/:provider/callback` - Callback для OAuth авторизації

### Бронювання
- GET `/api/bookings` - Отримати список бронювань
- GET `/api/bookings/:id` - Отримати деталі бронювання
- POST `/api/bookings` - Створити нове бронювання
- PUT `/api/bookings/:id` - Оновити бронювання
- DELETE `/api/bookings/:id` - Видалити бронювання
- PUT `/api/bookings/:id/status` - Оновити статус бронювання
- PUT `/api/bookings/:id/payment-status` - Оновити статус оплати
- POST `/api/bookings/:id/team-members` - Додати члена команди до бронювання
- DELETE `/api/bookings/:id/team-members/:member-id` - Видалити члена команди з бронювання

### Web інтерфейс
- GET `/bookings` - Список бронювань
- GET `/bookings/new` - Форма створення нового бронювання
- POST `/bookings` - Створення бронювання
- GET `/bookings/:id/edit` - Форма редагування бронювання
- PUT `/bookings/:id` - Оновлення бронювання
- DELETE `/bookings/:id` - Видалення бронювання

## 🔐 Аутентифікація та авторизація

### Типи аутентифікації
- **Локальна** (email/пароль)
- **OAuth інтеграції:**
  - Google
  - Facebook
  - Apple

### Токени
- **Access Token** - JWT токен для доступу до захищених ресурсів (час життя 15 хвилин)
- **Refresh Token** - JWT токен для оновлення access token (час життя 7 днів)

### Структура JWT токена
```json
{
  "sub": "user_uuid",
  "email": "user@example.com",
  "role": "user",
  "exp": 1637172747,
  "iat": 1637172147
}
```

### Рівні доступу
- **owner** - Власник акаунту, повний доступ до всіх функцій
- **admin** - Адміністратор, може керувати командою та має доступ до фінансів
- **retoucher** - Ретушер, може переглядати та редагувати проекти
- **editor** - Редактор, може переглядати та редагувати проекти
- **operator** - Оператор, може переглядати проекти
- **assistant** - Асистент, може переглядати проекти
- **user** - Звичайний користувач

### Дозволи
- **view_financials** - Перегляд фінансової інформації
- **edit_projects** - Редагування проектів
- **view_projects** - Перегляд проектів
- **manage_team** - Управління командою

## 📚 Команди розробки

```bash
# Запуск сервера
go run cmd/app/main.go

# Запуск Docker контейнерів
docker-compose up -d

# Зупинка Docker контейнерів
docker-compose down

# Запуск тестів
go test ./...

# Застосування міграцій
migrate -path migrations -database "postgres://..." up

# Відкат міграцій
migrate -path migrations -database "postgres://..." down
```

## 📱 UI Компоненти та Оформлення

### Шрифти
- **Основний шрифт**: Inter
- **Варіації**: Regular (400), Medium (500), Bold (700), Italic
- **Підтримка**: Повна підтримка кирилиці
- **Завантаження**: 
  - Через Google Fonts для онлайн-підтримки
  - Локальні файли для офлайн-підтримки 
  - Preload для прискорення завантаження

### Кольорова схема
- **Фон**: #FFFFFF (--color-background)
- **Фон секцій**: #F8F5F2 (--color-section-bg)
- **Світлий акцент**: #E3D5CA (--color-light-accent)
- **Нейтральний акцент**: #D6CCC2 (--color-neutral-accent)
- **Контрастний акцент**: #D5BDAF (--color-contrast-accent)
- **Основний текст**: #2A2A2A (--color-text-primary)
- **Другорядний текст**: #4A4744 (--color-text-secondary)

## 🔄 Поточний статус
- Етап: 9 - Стабілізація та посилення безпеки CORS
- Фокус: Виправлення проблем CORS та вдосконалення структури проекту
- Імплементовано: 
  - Повноцінну локальну автентифікацію з системою логіна та реєстрації
  - OAuth повністю функціональний (Google, Facebook, Apple)
  - Роботу з JWT токенами, включаючи оновлення через refresh токен
  - Систему ролей та дозволів з підтримкою різних рівнів доступу
  - Безпечну CORS конфігурацію, що відповідає вимогам безпеки
  - Усунуто проблему з AllowOrigins та AllowCredentials в CORS конфігурації
  - Спрощено структуру проекту шляхом використання єдиної точки входу (cmd/app/main.go)
  - Покращену систему обробки помилок в middleware для авторизації
  - Модель бронювань з підтримкою різних типів подій (wedding, portrait, family, event, commercial)
  - Форми для створення та редагування бронювань з підтримкою типу події
  - REST API для всіх операцій з бронюваннями та автентифікацією
  - Middleware для авторизації та перевірки прав користувача
  - Виправлено обробку JSONB полів (Permissions, Settings) для коректного збереження в базі даних
  - Оновлено процеси реєстрації та створення адміністратора для коректної конвертації об'єктів в JSON
  - Забезпечено правильне використання JSONB полів у всьому додатку для уникнення помилок типізації
- Наступний крок: Розширення функціональності користувацького профілю та додавання розширеної фільтрації бронювань за типом події та датою

## 📌 Нотатки та важливі моменти
- Використовуємо UUID для ідентифікаторів
- JSONB для гнучких полів (Settings, CustomFields, TeamMembers, Permissions)
- Повна підтримка множинних методів автентифікації
- JWT аутентифікація з refresh токенами
- Багаторівнева система прав доступу
- Підтримка команди з ієрархічною структурою
- Реалізовано можливість реєстрації через OAuth (Google, Facebook, Apple)
- Додано підтримку різних типів подій (event_type) у бронюваннях
- Підтримка рольової автентифікації через middleware
- Інтеграція з Tabler UI для єдиного стилю всіх компонентів
- Взаємодія між клієнтським інтерфейсом і серверними ендпоінтами через REST API

## ✅ Модель користувача (User)
Успішно оновлено з підтримкою:
- Різних ролей користувачів
- Управління командою
- OAuth автентифікації
- Детальних налаштувань дозволів

## ✅ Модель бронювання (Booking)
Успішно оновлено з підтримкою:
- Різних типів подій (event_type)
- Статусів оплати і бронювання
- Повного циклу управління бронюваннями
- Командної роботи через team_members
- Кастомних полів для розширення функціоналу
