Вся комунікація і коментарі українською мовою.

# TimeBride CRM

## 📋 Опис проекту
TimeBride - це SaaS CRM-платформа для фотографів, відеографів, та в майбутньому весільних агенцій, що фокусується на:
- Максимальній простоті та мінімалізмі інтерфейсу
- Легкому управлінні зйомками (бронюваннями) та процесом віддачі матеріалу
- Високій швидкодії та оптимізованій роботі
- Синхронізації всіх бронювань з Google Calendar, iCal, Google Sheets та в майбутньому можливо інших сервісів

## 🏗 Поточна структура проекту
.
├── cmd/
│ └── app/
│ └── main.go # Точка входу в програму
├── internal/
│ ├── config/ # Конфігурація програми
│ │ └── config.go
│ ├── models/ # Моделі даних
│ │ ├── user.go
│ │ ├── booking.go
│ │ ├── template.go
│ │ ├── custom_fields.go
│ │ ├── file.go
│ │ └── validation.go
│ ├── handlers/ # HTTP обробники
│ │ ├── auth_handler.go
│ │ ├── user_handler.go
│ │ ├── booking_handler.go
│ │ ├── template_handler.go
│ │ └── file_handler.go
│ ├── services/ # Бізнес-логіка
│ │ ├── auth/
│ │ │ └── auth.go
│ │ ├── booking/
│ │ │ └── booking.go
│ │ ├── dashboard/
│ │ │ └── dashboard.go
│ │ ├── file/
│ │ │ ├── service.go
│ │ │ └── file.go
│ │ ├── template/
│ │ │ └── service.go
│ │ └── user/
│ │ └── service.go
│ ├── repositories/ # Робота з БД
│ │ ├── repository.go
│ │ ├── user_repository.go
│ │ ├── booking_repository.go
│ │ ├── template_repository.go
│ │ └── file_repository.go
│ ├── middleware/ # Middleware компоненти
│ │ └── auth.go
│ └── router/ # Маршрутизація
│ └── router.go
├── migrations/ # SQL міграції
│ ├── 000001_init.up.sql
│ ├── 000001_init.down.sql
│ └── 000001_clean.sql
├── pkg/ # Публічні пакети
│ ├── logger/
│ └── database/
├── web/ # Веб-інтерфейс
│ ├── src/
│ │ ├── css/
│ │ │ ├── main.css
│ │ │ ├── utilities.css
│ │ │ └── layout.css
│ │ ├── js/
│ │ └── templates/
│ └── static/
│ └── css/
├── config.yaml # Конфігурація
├── docker-compose.yml # Docker конфігурація
├── Dockerfile # Dockerfile для збірки
├── .env # Змінні середовища
└── Makefile # Команди для розробки

## 📝 План розробки

### ✅ Етап 1: Базова інфраструктура (Виконано)
- [x] Ініціалізація Go проекту
- [x] Налаштування Docker і docker-compose
- [x] Базова структура проекту
- [x] Підключення до PostgreSQL

### ✅ Етап 2: Основні компоненти (Виконано)
- [x] Моделі даних
  - [x] User (користувач)
  - [x] Booking (бронювання)
  - [x] Template (шаблон)
  - [x] CustomFields (кастомні поля)
  - [x] File (файл)
  - [x] Validation (валідація)
- [x] Міграції бази даних
- [x] Репозиторії
  - [x] UserRepository
  - [x] BookingRepository
  - [x] TemplateRepository
  - [x] FileRepository
- [x] Сервіси
  - [x] UserService
  - [x] BookingService
  - [x] TemplateService
  - [x] FileService
  - [x] AuthService
  - [x] DashboardService
- [x] HTTP обробники
  - [x] AuthHandler
  - [x] UserHandler
  - [x] BookingHandler
  - [x] TemplateHandler
  - [x] FileHandler

### ✅ Етап 3: Основний функціонал (Виконано)
- [x] Аутентифікація та авторизація
  - [x] JWT токени
  - [x] Middleware для перевірки авторизації
- [x] API ендпоінти для бронювань
  - [x] CRUD операції
  - [x] Фільтрація та пошук
- [x] Управління користувачами
  - [x] CRUD операції
  - [x] Ролі та права доступу
- [x] Завантаження файлів
  - [x] Зберігання в Backblaze B2
  - [x] Обробка файлів
  - [x] Метадані та теги

### 🏃‍♂️ Етап 4: Інтерфейс та UX (В процесі)
- [x] Базова структура веб-інтерфейсу
  - [x] CSS утиліти та компоненти
  - [x] Базовий layout
- [ ] Календар бронювань
- [ ] Галерея зображень
- [ ] Адаптивний дизайн

### 🚀 Етап 5: Розширений функціонал
- [ ] Інтеграція з Google Calendar
- [ ] Платіжна система
- [x] Кастомні поля для бронювань
- [x] Шаблони документів

### 🔧 Етап 6: Технічні покращення
- [x] Валідація даних
  - [x] Інтеграція go-playground/validator
  - [x] Кастомні правила валідації
  - [x] Валідація вхідних даних
- [x] Обробка помилок
  - [x] Централізований механізм
  - [x] Структуроване логування
  - [x] Стандартизовані відповіді
- [ ] Тестування
  - [x] Юніт-тести для сервісів
  - [ ] Інтеграційні тести
  - [ ] Навантажувальні тести
- [ ] Документація
  - [ ] Swagger/OpenAPI
  - [ ] Автоматична генерація
  - [ ] Приклади використання
- [x] Безпека
  - [x] Шифрування даних
  - [x] Безпечне зберігання
  - [x] Аудит безпеки
- [ ] Масштабованість
  - [ ] Кешування
  - [ ] Оптимізація запитів
  - [ ] Балансування навантаження
- [x] Логування
  - [x] Структуроване логування
  - [x] Рівні логування
  - [x] Ротація логів
- [ ] CI/CD
  - [ ] Автоматизація тестування
  - [ ] Автоматичне розгортання
  - [ ] Моніторинг процесу
- [ ] Моніторинг
  - [ ] Метрики продуктивності
  - [ ] Алерти
  - [ ] Дашборди
- [ ] Локалізація
  - [ ] Підтримка i18n
  - [ ] Управління перекладами
  - [ ] Тестування локалізації

## 📊 Структура бази даних

### Users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255),
    phone VARCHAR(50),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    theme VARCHAR(20) DEFAULT 'light',
    subscription_plan VARCHAR(50) DEFAULT 'free',
    google_calendar_token TEXT,
    telegram_chat_id VARCHAR(100),
    settings JSONB DEFAULT '{}',
    team_settings JSONB DEFAULT '{}',
    custom_fields JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Bookings
```sql
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    event_type VARCHAR(50) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    is_all_day BOOLEAN DEFAULT false,
    client_name VARCHAR(255) NOT NULL,
    client_phone VARCHAR(50),
    client_email VARCHAR(255),
    location VARCHAR(255),
    price DECIMAL(10,2),
    currency VARCHAR(3) DEFAULT 'USD',
    payment_status VARCHAR(50) DEFAULT 'pending',
    team_members JSONB DEFAULT '[]',
    custom_fields JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Templates
```sql
CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    event_type VARCHAR(50) NOT NULL,
    fields_template JSONB DEFAULT '{}',
    team_template JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Files
```sql
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) NOT NULL,
    booking_id UUID REFERENCES bookings(id),
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    bucket VARCHAR(100) NOT NULL,
    key VARCHAR(255) NOT NULL,
    cdn_url VARCHAR(255) NOT NULL,
    description TEXT,
    tags TEXT,
    metadata TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## 🛠 API Endpoints (Реалізовані)

### Аутентифікація
POST /api/v1/auth/register
POST /api/v1/auth/login
GET /api/v1/auth/me

### Користувачі
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
PUT    /api/v1/users/:id/password
PUT    /api/v1/users/:id/settings

### Бронювання
GET    /api/v1/bookings
POST   /api/v1/bookings
GET    /api/v1/bookings/:id
PUT    /api/v1/bookings/:id
DELETE /api/v1/bookings/:id
PUT    /api/v1/bookings/:id/status
PUT    /api/v1/bookings/:id/payment-status
POST   /api/v1/bookings/:id/team-members
DELETE /api/v1/bookings/:id/team-members/:member-id

### Шаблони
GET    /api/v1/templates
POST   /api/v1/templates
GET    /api/v1/templates/:id
PUT    /api/v1/templates/:id
DELETE /api/v1/templates/:id

### Файли
POST   /api/v1/files
GET    /api/v1/files
GET    /api/v1/files/:id
DELETE /api/v1/files/:id
GET    /api/v1/files/:id/download

## 📚 Команди розробки

```bash
# Запуск сервера
make run

# Запуск Docker контейнерів
make docker-up

# Зупинка Docker контейнерів
make docker-down

# Запуск тестів
make test

# Застосування міграцій
make migrate-up

# Відкат міграцій
make migrate-down
```

## 🔄 Поточний статус
- Етап: 4 - Інтерфейс та UX
- Фокус: Розробка веб-інтерфейсу та компонентів
- Наступний крок: Реалізація календаря бронювань та галереї зображень

## 📌 Нотатки та важливі моменти
- Використовуємо UUID для ідентифікаторів
- JSONB для гнучких полів (Settings, CustomFields, TeamMembers)
- Готуємося до мультимовності
- Враховуємо можливість white-label рішення
- Реалізовано базову структуру для роботи з кастомними полями
- Додано підтримку різних типів подій та статусів бронювань
- Реалізовано валідацію шаблонів для полів та команди
- Реалізовано JWT аутентифікацію та авторизацію
- Додано middleware для перевірки авторизації
- Реалізовано CRUD операції для всіх основних сутностей
- Додано підтримку ролей та прав доступу
- Реалізовано завантаження файлів на Backblaze B2 з підтримкою CDN
- Файли зберігаються з структурованою ієрархією: {user-id}/{YYYY-MM}/{safe-name}-{unique-id}{ext}
- Додано метадані та теги для файлів
- Реалізовано безпечне іменування файлів з санітизацією небезпечних символів
- Реорганізовано сервіси в окремі пакети для кращої структури коду
- Додано базові CSS утиліти та компоненти для веб-інтерфейсу
- Підготовлено план технічних покращень для масштабування проекту
- Враховується необхідність високої якості коду та тестування
- Планується впровадження сучасних практик розробки та моніторингу
