package database

import (
	"log"
	"timebride/internal/models"

	"gorm.io/gorm"
)

// RunMigrations виконує міграції бази даних
func RunMigrations(db *gorm.DB) error {
	log.Println("Запуск міграцій...")

	// Перевірка наявності таблиці users
	if !db.Migrator().HasTable(&models.User{}) {
		log.Println("Таблиця users не існує, створюю...")
		if err := db.Migrator().CreateTable(&models.User{}); err != nil {
			log.Printf("Помилка при створенні таблиці users: %v", err)
			return err
		}
		log.Println("Таблиця users успішно створена")
	}

	// Перевірка та оновлення поля name в таблиці users
	if db.Migrator().HasColumn(&models.User{}, "name") {
		log.Println("Колонка 'name' вже існує в таблиці users")

		// Перевіряємо чи є записи з NULL в полі name
		var nullNameCount int64
		db.Raw("SELECT COUNT(*) FROM users WHERE name IS NULL").Scan(&nullNameCount)

		if nullNameCount > 0 {
			log.Printf("Знайдено %d записів з NULL в полі name, оновлюю...", nullNameCount)
			// Заповнюємо дані на основі email
			if err := db.Exec("UPDATE users SET name = split_part(email, '@', 1) WHERE name IS NULL").Error; err != nil {
				log.Printf("Помилка при оновленні значень name: %v", err)
				return err
			}
			log.Println("Записи успішно оновлені")
		}
	} else {
		log.Println("Колонка 'name' не існує, додаю...")
		// Спочатку додаємо колонку як nullable
		if err := db.Exec("ALTER TABLE users ADD COLUMN name VARCHAR(255)").Error; err != nil {
			log.Printf("Помилка при додаванні колонки 'name': %v", err)
			return err
		}
		log.Println("Колонка 'name' додана як nullable")

		// Заповнюємо дані на основі email
		if err := db.Exec("UPDATE users SET name = split_part(email, '@', 1) WHERE name IS NULL").Error; err != nil {
			log.Printf("Помилка при оновленні значень name: %v", err)
			return err
		}
		log.Println("Заповнено значення name для всіх користувачів")

		// Після заповнення даних встановлюємо NOT NULL обмеження
		if err := db.Exec("ALTER TABLE users ALTER COLUMN name SET NOT NULL").Error; err != nil {
			log.Printf("Помилка при встановленні NOT NULL для 'name': %v", err)
			return err
		}
		log.Println("Встановлено обмеження NOT NULL для колонки 'name'")
	}

	// Автоматична міграція інших моделей
	log.Println("Виконання автоматичних міграцій...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Booking{},
		&models.Template{},
		&models.File{},
	); err != nil {
		log.Printf("Помилка при виконанні міграцій: %v", err)
		return err
	}

	log.Println("Міграції успішно виконані")
	return nil
}
