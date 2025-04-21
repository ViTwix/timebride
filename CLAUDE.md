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
│   ├── cache/                  # Механізми кешування
│   ├── router/                 # Маршрутизація
│   └── utils/                  # Загальні утиліти
├── migrations/                 # SQL міграції
├── pkg/                        # Публічні пакети
│   ├── logger/                 # Логування
│   └── database/               # Підключення до БД
├── web/
│   ├── templates/              # HTML шаблони
│   │   └── layout.html         # Основний шаблон сайту
│   │   └── dashboard/          # Шаблони дашборду
│   ├── public/                 # Публічні статичні файли
│   │   ├── css/
│   │   │   ├── tabler.min.css  # Основний стиль Tabler
│   │   │   ├── tabler-icons.min.css  # Іконки Tabler
│   │   │   └── main.css        # Єдиний CSS файл проекту
│   │   ├── js/                 # JavaScript файли
│   │   │   ├── tabler.min.js   # Основний JS файл Tabler
│   │   │   ├── main.js         # Головний JS файл проекту
│   │   │   └── force-styles.js # Скрипт для примусового застосування стилів
│   │   ├── img/                # Зображення
│   │   └── fonts/              # Шрифти (включно з Tabler Icons)
├── config.yaml                 # Основна конфігурація
├── docker-compose.yml          # Docker конфігурація
├── Dockerfile                  # Dockerfile для збірки
└── Makefile                    # Команди для розробки
```

## 📝 Оптимізації та спрощення

### 🔄 Оптимізація CSS
1. **Об'єднано всі стилі в один файл main.css** - замість багатьох окремих файлів використовується один головний, що:
   - Зменшує кількість HTTP запитів
   - Покращує продуктивність
   - Спрощує підтримку кодової бази
   - Ліквідує дублювання стилів
2. **Спрощено каскадність селекторів** - усунуто надмірно специфічні селектори
3. **Структуровано CSS** - код розділений на логічні секції з коментарями:
   - Змінні
   - Шрифти
   - Кнопки
   - Основні стилі
   - Responsive
   - Utility Classes
4. **Реалізовано перевизначення таблерівських класів** - забезпечено єдину кольорову палітру:
   - Перевизначено клас `.text-secondary` для використання кольорів з нашої палітри
   - Додано спеціальні стилі для заголовків (`h1.display-3`, `h2.display-5`)
   - Створено додаткові класи для нашої кольорової схеми
5. **Додано fallback для правильного застосування стилів** - скрипт для запобігання проблем з відображенням:
   - Примусове застосування кольорів для тексту
   - Гарантоване застосування шрифту Inter до всіх елементів
   - Обхід кешування стилів

### 🧩 Оптимізація серверної частини
1. **Спрощено структуру сервісів** - об'єднано схожі сервіси
2. **Покращено ін'єкцію залежностей** - зменшено зв'язаність компонентів
3. **Видалено повторювані функції** - утилітарні функції винесено в окремий пакет
4. **Оптимізовано роутинг** - налаштовано статичні маршрути:
   - Додано маршрут `/fonts` для доступу до шрифтів Tabler Icons
   - Дозволено публічний доступ до шрифтів без аутентифікації

### 🗂 Оптимізація шаблонів
1. **Зменшено вкладеність** - спрощено ієрархію шаблонів
2. **Уніфіковано підхід до макетів** - всі сторінки використовують єдиний базовий макет
3. **Зменшено кількість часткових шаблонів** - об'єднано дрібні компоненти
4. **Оновлено підключення стилів** - використання єдиного CSS файлу замість багатьох окремих
5. **Додано preload для шрифтів** - прискорено завантаження шрифтів:
   - Додано preload для Inter і Tabler Icons
   - Реалізовано механізм виявлення завантаження шрифтів

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

### CSS структура (Оптимізована)
- **main.css** - Єдиний файл стилів, що містить:
  - CSS змінні
  - Шрифти та їх налаштування
  - Перевизначення Tabler класів
  - Компоненти UI
  - Адаптивні стилі
  - Утилітарні класи

## 🔄 Поточний статус
- Етап: 4 - Інтерфейс та UX
- Фокус: Оптимізація та спрощення компонентів, вирішення проблеми зі стилями
- Наступний крок: Тестування та покращення UI компонентів

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
- Налаштовано правильне відображення шрифтів Tabler Icons
- Реалізовано механізм перевизначення кольорів для забезпечення єдиної стилістики
- Додано скрипт для примусового застосування стилів
- Додано додаткові класи для адаптивності
- Оптимізовано структуру шляхів до статичних файлів
- Впроваджено перевірку завантаження шрифтів
- Вирішено проблему з кольорами тексту заголовків
