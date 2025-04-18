package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"timebride/internal/config"
	"timebride/internal/repositories"
	"timebride/internal/router"
	"timebride/internal/services"
	"timebride/pkg/database"
)

func main() {
	// Завантажуємо конфігурацію
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Підключаємося до бази даних
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Створюємо репозиторії
	userRepo := repositories.NewUserRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)
	templateRepo := repositories.NewTemplateRepository(db)

	// Створюємо сервіси
	userService := services.NewUserService(userRepo)
	bookingService := services.NewBookingService(bookingRepo)
	templateService := services.NewTemplateService(templateRepo)

	// Створюємо роутер
	handler := router.Router(cfg, userService, bookingService, templateService)

	// Налаштовуємо сервер
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаємо сервер
	log.Printf("Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
