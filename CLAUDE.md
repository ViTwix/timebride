# Вся комунікація і коментарі українською мовою.
# Ultrathink during each request processing

# Увесь інтерфейс сайту — як адмінка, так і клієнтська частина — повинен бути побудований виключно на базі Tabler UI (https://tabler.io) https://github.com/tabler/tabler
Врахуй такі ключові вимоги:
Мінімалістичний, легкий для сприйняття інтерфейс
Повна адаптивність
Таблиці, форми, календарі, модалки та інші елементи — лише ті, що є в Tabler UI
Якщо потрібно щось кастомізувати — роби це в стилістиці Tabler
Всі компоненти мають мати єдину візуальну мову згідно Tabler
На виході я очікую HTML+CSS/JS компоненти або SSR-ready шаблони, повністю сумісні з Tabler UI. Якщо використовуєш JS — тільки Vanilla або Alpine.js.
Якщо щось потрібно придумати — інтерпретуй це через логіку Tabler UI.

## 📋 Опис проекту
TimeBride - це SaaS CRM-платформа для фотографів, відеографів, та в майбутньому весільних агенцій, що фокусується на:
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
│ │ └── config.go # Завантаження та обробка конфігурації
│ ├── models/ # Моделі даних
│ │ ├── user.go # Модель користувача з методами валідації
│ │ ├── booking.go # Модель бронювання з методами для подій
│ │ ├── template.go # Модель шаблонів для документів та форм
│ │ ├── custom_fields.go # Обробка кастомних полів для форм
│ │ ├── file.go # Модель для керування файлами та метаданими
│ │ └── validation.go # Загальні функції валідації
│ ├── handlers/ # HTTP обробники
│ │ ├── auth_handler.go # Авторизація та реєстрація
│ │ ├── user_handler.go # Керування користувачами
│ │ ├── booking_handler.go # Обробка бронювань та подій
│ │ ├── template_handler.go # Керування шаблонами
│ │ └── file_handler.go # Завантаження та керування файлами
│ ├── services/ # Бізнес-логіка
│ │ ├── auth/
│ │ │ └── auth.go # Автентифікація та JWT токени
│ │ ├── booking/
│ │ │ └── booking.go # Логіка бронювань
│ │ ├── dashboard/
│ │ │ └── dashboard.go # Дані для дашборду
│ │ ├── file/
│ │ │ ├── service.go # Завантаження та зберігання файлів
│ │ │ └── file.go # Додаткові функції для файлів
│ │ ├── template/
│ │ │ └── service.go # Бізнес-логіка для шаблонів
│ │ └── user/
│ │ └── service.go # Керування користувачами
│ ├── repositories/ # Робота з БД
│ │ ├── repository.go # Базовий інтерфейс репозиторію
│ │ ├── user_repository.go # Запити для користувачів
│ │ ├── booking_repository.go # Запити для бронювань
│ │ ├── template_repository.go # Запити для шаблонів
│ │ └── file_repository.go # Запити для файлів
│ ├── middleware/ # Middleware компоненти
│ │ └── auth.go # Перевірка авторизації та JWT
│ ├── cache/ # Механізми кешування
│ │ └── cache.go # Інтерфейси та імплементації для Redis/MemCache/InMemory
│ ├── loadbalancer/ # Компоненти для балансування навантаження
│ │ └── loadbalancer.go # Різні стратегії балансування (Round Robin, Least Connection)
│ └── router/ # Маршрутизація
│ └── router.go # Налаштування маршрутів та middleware
├── migrations/ # SQL міграції
│ ├── 000001_init.up.sql # Початкові таблиці
│ ├── 000001_init.down.sql # Відкат міграцій
│ └── 000001_clean.sql # Видалення всіх даних
├── pkg/ # Публічні пакети
│ ├── logger/ # Логування
│ └── database/ # Підключення до БД
├── web/ # Веб-інтерфейс
│ ├── src/
│ │ ├── css/ # CSS стилі для клієнтської частини
│ │ ├── js/ # JavaScript для клієнтської частини
│ │ └── templates/ # HTML шаблони
│ │ └── partials/ # Повторювані частини шаблонів (head, sidebar, footer)
│ │ └── layout.html # Основний шаблон сторінки
│ │ └── dashboard.html # Шаблон головної сторінки
│ │ └── auth/ # Шаблони для авторизації
│ │ └── clients/ # Шаблони для клієнтів
│ │ └── calendar.html # Шаблон календаря
│ │ └── gallery.html # Шаблон галереї
│ │ └── profile.html # Шаблон профілю користувача
│ └── public/ # Публічні статичні файли
│ │ ├── css/ # Компільовані CSS стилі
│ │ │ ├── styles.css # Основні стилі
│ │ │ ├── variables.css # CSS змінні (кольори, шрифти, тіні, радіуси)
│ │ │ ├── fonts.css # Налаштування шрифту Inter
│ │ │ ├── tabler.min.css # Мінімізований Tabler CSS
│ │ │ ├── tabler-icons.min.css # Іконки Tabler
│ │ ├── js/ # JavaScript файли
│ │ ├── img/ # Зображення (favicon, логотип)
│ │ ├── fonts/ # Шрифти (Inter у різних форматах)
├── config/ # Конфігураційні файли
│ ├── scaling.yaml # Налаштування для масштабування (кешування, балансування)
├── config.yaml # Основна конфігурація
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
  - [x] Шрифти та іконки
- [x] Модульна структура шаблонів
  - [x] Partials для повторюваних елементів
  - [x] Шаблон для контенту
- [x] Стилізація компонентів
  - [x] Кольорова схема
  - [x] Шрифт Inter (Regular, Medium, Bold, Italic)
  - [x] Консистентний дизайн елементів
- [x] Календар бронювань
- [x] Галерея зображень
- [x] Адаптивний дизайн

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
- [x] Масштабованість
  - [x] Кешування
    - [x] Механізми кешування (Redis, Memcached, In-Memory)
    - [x] Стратегії інвалідації кешу
  - [x] Балансування навантаження
    - [x] Різні стратегії балансування (Round Robin, Least Connection, IP Hash)
    - [x] Перевірка здоров'я серверів
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

## 📱 UI Компоненти та Оформлення

### Шрифти
- **Основний шрифт**: Inter
- **Варіації**: Regular (400), Medium (500), Bold (700), Italic
- **Підтримка**: Повна підтримка кирилиці
- **Файли**: Шрифти завантажуються як через Google Fonts, так і локально

### Кольорова схема
- **Фон**: #FFFFFF
- **Фон секцій**: #F8F5F2
- **Світлий акцент**: #E3D5CA (використовується для фону хедера і футера)
- **Нейтральний акцент**: #D6CCC2
- **Контрастний акцент**: #D5BDAF
- **Основний текст**: #2A2A2A
- **Другорядний текст**: #4A4744

### CSS структура
- **variables.css** - Базові змінні оформлення (кольори, шрифти, тіні, радіуси)
- **fonts.css** - Налаштування шрифтів
- **styles.css** - Основні стилі, перевизначення компонентів Tabler

### Шаблони UI
- **Модульна структура** з використанням partials:
  - **head.html** - Метатеги, CSS та JavaScript підключення
  - **sidebar.html** - Бокове меню та навігація
  - **footer.html** - Підвал сайту
  - **layout.html** - Основний шаблон, який об'єднує всі partials

## 🔄 Поточний статус
- Етап: 4 - Інтерфейс та UX
- Фокус: Вдосконалення UI та UX компонентів
- Наступний крок: Завершення інтеграції всіх шаблонів та їх стилізація

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
- Впроваджено модульну структуру шаблонів
- Встановлено шрифт Inter з різними варіаціями
- Налаштовано єдину кольорову схему
- Реалізовано кешування для підвищення продуктивності
- Додано компоненти для балансування навантаження
- Використана архітектура, що забезпечує горизонтальне масштабування
- Підготовлено план технічних покращень для масштабування проекту
- Враховується необхідність високої якості коду та тестування
- Планується впровадження сучасних практик розробки та моніторингу
