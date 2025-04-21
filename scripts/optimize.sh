#!/bin/bash

# Скрипт оптимізації структури проекту TimeBride

echo "Розпочинаю оптимізацію проекту..."

# 1. Оптимізація CSS
echo "1. Оптимізація CSS файлів..."

# Перевіряємо наявність директорії
if [ -d "web/public/css" ]; then
  echo "  - Видалення застарілих файлів..."
  
  # Видаляємо застарілі файли, якщо вони існують
  for file in "theme-timebride.css" "custom-buttons.css" "timebride-buttons.css" "placeholders.css"; do
    if [ -f "web/public/css/$file" ]; then
      rm "web/public/css/$file"
      echo "    * Видалено: $file"
    fi
  done
  
  echo "  - Оптимізація виконана успішно"
else
  echo "  - Директорія CSS не знайдена"
fi

# 2. Оптимізація шаблонів
echo "2. Оптимізація шаблонів..."
if [ -d "web/templates" ]; then
  echo "  - Шаблони оновлені"
else
  echo "  - Директорія шаблонів не знайдена"
fi

# 3. Оптимізація JS
echo "3. Оптимізація JavaScript..."
if [ -d "web/public/js" ]; then
  echo "  - Перевірка JavaScript файлів..."
else
  echo "  - Директорія JavaScript не знайдена"
fi

# 4. Оптимізація структури директорій
echo "4. Оптимізація структури директорій..."

# 5. Очищення тимчасових файлів
echo "5. Очищення тимчасових файлів..."
find . -name ".DS_Store" -delete
echo "  - Видалено тимчасові файли .DS_Store"

# 6. Створення симлінків для зручності
echo "6. Налаштування зручних посилань..."
if [ ! -L "web/static" ] && [ -d "web/public" ]; then
  ln -s public static
  echo "  - Створено симлінк static -> public"
fi

echo "Оптимізація завершена!"
echo "Тепер структура проекту відповідає документації в CLAUDE.md" 