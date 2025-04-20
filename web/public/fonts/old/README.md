# Інструкція по встановленню шрифтів Inter

## Завантаження шрифтів Inter

Для правильної роботи сайту з кирилицею, потрібно замінити поточні файли-заглушки на справжні шрифти Inter. Ось як це зробити:

1. Перейдіть на сайт Google Fonts: https://fonts.google.com/specimen/Inter

2. Натисніть "Download family" для завантаження всіх варіантів шрифту

3. Розпакуйте ZIP-архів

4. Конвертуйте файли TTF у WOFF та WOFF2:
   - Можна використати онлайн-конвертер: https://cloudconvert.com/ttf-to-woff або https://transfonter.org/
   - Або використайте локальний інструмент, такий як FontForge

5. Розмістіть отримані файли у цій директорії:
   - Inter-Regular.woff і Inter-Regular.woff2
   - Inter-Medium.woff і Inter-Medium.woff2
   - Inter-Bold.woff і Inter-Bold.woff2
   - Inter-Italic.woff і Inter-Italic.woff2

## Альтернативний варіант

Якщо не потрібна локальна копія, то можна покладатися виключно на Google Fonts. В цьому випадку:

1. Переконайтесь, що в файлі `web/src/templates/partials/head.html` є рядок:
```html
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700&display=swap&subset=cyrillic,cyrillic-ext" rel="stylesheet">
```

2. Параметр `subset=cyrillic,cyrillic-ext` забезпечує підтримку кирилиці.

Зверніть увагу: використання локальних шрифтів рекомендоване для кращої продуктивності і незалежності від зовнішніх сервісів. 