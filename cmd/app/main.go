package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"timebride/internal/config"
	"timebride/internal/pkg/database"
)

func main() {
	// Завантаження конфігурації
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Підключення до бази даних
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer db.Close()

	// Створення Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Базові роути
	api := app.Group("/api/v1")

	// Роут перевірки здоров'я системи
	api.Get("/health", func(c *fiber.Ctx) error {
		// Перевіряємо також підключення до бази даних
		if err := db.HealthCheck(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "error",
				"error":  "Database connection failed",
			})
		}

		return c.JSON(fiber.Map{
			"status":      "ok",
			"environment": cfg.Server.Env,
			"database":    "connected",
		})
	})

	// Запуск сервера
	log.Printf("Starting server on port %d in %s mode", cfg.Server.Port, cfg.Server.Env)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)))
}
