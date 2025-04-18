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
│ ├── handlers/ # HTTP обробники
│ │ └── project_handler.go
│ ├── services/ # Бізнес-логіка
│ ├── repositories/ # Робота з БД
│ └── middleware/ # Middleware компоненти
├── migrations/ # SQL міграції
├── pkg/ # Публічні пакети
│ ├── logger/
│ └── database/
├── web/ # Веб-ресурси
│ ├── templates/
│ ├── static/
│ └── assets/
├── config.yaml # Конфігурація
├── docker-compose.yml # Docker конфігурація
└── Makefile # Команди для розробки

## 📝 План розробки

### ✅ Етап 1: Базова інфраструктура (Виконано)
- [x] Ініціалізація Go проекту
- [x] Налаштування Docker і docker-compose
- [x] Базова структура проекту
- [x] Підключення до PostgreSQL

### 🏃‍♂️ Етап 2: Основні компоненти (В процесі)
- [ ] Моделі даних
  - [ ] User (користувач)
  - [ ] Booking (бронювання)
- [ ] Міграції бази даних
- [ ] Репозиторії
- [ ] Аутентифікація та авторизація

### 📅 Етап 3: Основний функціонал
- [ ] API ендпоінти для бронювань
- [ ] Управління користувачами
- [ ] Завантаження файлів
- [ ] Система сповіщень

### 🎨 Етап 4: Інтерфейс та UX
- [ ] Базовий веб-інтерфейс
- [ ] Календар бронювань
- [ ] Галерея зображень
- [ ] Адаптивний дизайн

### 🚀 Етап 5: Розширений функціонал
- [ ] Інтеграція з Google Calendar
- [ ] Платіжна система
- [ ] Кастомні поля для бронювань
- [ ] Шаблони документів

## 📊 Структура бази даних

### Users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(50) NOT NULL,
    settings JSONB DEFAULT '{}',
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
    status VARCHAR(50) NOT NULL,
    client_name VARCHAR(255),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    is_all_day BOOLEAN DEFAULT false,
    custom_fields JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🛠 API Endpoints (Заплановані)

### Аутентифікація
POST /api/v1/auth/register
POST /api/v1/auth/login
GET /api/v1/auth/me

### Бронювання
GET    /api/v1/bookings
POST   /api/v1/bookings
GET    /api/v1/bookings/:id
PUT    /api/v1/bookings/:id
DELETE /api/v1/bookings/:id

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
- Етап: 2 - Основні компоненти
- Фокус: Налаштування моделей даних та міграцій
- Наступний крок: Реалізація моделей User та Booking

## 📌 Нотатки та важливі моменти
- Використовуємо UUID для ідентифікаторів
- JSONB для гнучких полів
- Готуємося до мультимовності
- Враховуємо можливість white-label рішення
