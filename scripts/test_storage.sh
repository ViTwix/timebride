#!/bin/bash

# Перевіряємо наявність змінних середовища
if [ -z "$BACKBLAZE_ACCOUNT_ID" ] || [ -z "$BACKBLAZE_APPLICATION_KEY" ]; then
    echo "Помилка: Не встановлені змінні середовища BACKBLAZE_ACCOUNT_ID та/або BACKBLAZE_APPLICATION_KEY"
    exit 1
fi

# Запускаємо тест
echo "Запуск тесту з'єднання з Backblaze B2..."
go test -v ./internal/services -run TestBackblazeConnection

# Перевіряємо результат
if [ $? -eq 0 ]; then
    echo "Тест успішно пройдено! З'єднання з Backblaze B2 працює."
else
    echo "Тест не пройдено. Перевірте налаштування та логи помилок."
fi 