package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/crypto/bcrypt"

	"timebride/internal/config"
	"timebride/internal/database"
	"timebride/internal/handlers"
	"timebride/internal/models"
	"timebride/internal/repositories"
	"timebride/internal/router"
	"timebride/internal/services"
	"timebride/internal/services/booking"
	"timebride/internal/services/file"
	"timebride/internal/services/template"
	"timebride/internal/services/user"
)

func main() {
	ctx := context.Background()

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

	// Виконуємо міграції бази даних
	log.Println("Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Створення адміністратора, якщо він не існує
	log.Println("Checking for admin user...")
	var adminCount int64
	if err := db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&adminCount).Error; err != nil {
		log.Printf("Error checking for admin user: %v", err)
	} else if adminCount == 0 {
		log.Println("Admin user not found, creating default admin...")
		// Створення хешу паролю
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
		} else {
			// Конвертуємо дозволи в JSON
			permissions, err := json.Marshal(models.DefaultPermissions(models.RoleAdmin))
			if err != nil {
				log.Printf("Error marshaling permissions: %v", err)
				return
			}

			// Створення адміністратора
			adminUser := models.User{
				Email:        "admin@timebride.com",
				PasswordHash: string(passwordHash),
				Name:         "Адміністратор системи",
				Role:         models.RoleAdmin,
				Language:     "uk",
				Permissions:  permissions,
			}
			if err := db.Create(&adminUser).Error; err != nil {
				log.Printf("Error creating admin user: %v", err)
			} else {
				log.Println("Default admin user created successfully")
			}
		}
	} else {
		log.Println("Admin user already exists")
	}

	// Створюємо S3 клієнт
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.Storage.Backblaze.Endpoint,
			SigningRegion:     cfg.Storage.Backblaze.Region,
			HostnameImmutable: true,
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithEndpointResolverWithOptions(customResolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.Storage.Backblaze.AccountID,
			cfg.Storage.Backblaze.ApplicationKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Failed to create AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	// Створюємо репозиторії
	userRepo := repositories.NewUserRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)
	templateRepo := repositories.NewTemplateRepository(db)
	fileRepo := repositories.NewFileRepository(db)

	// Створюємо сервіси
	userSvc := user.NewService(userRepo)
	bookingSvc := booking.New(bookingRepo, userSvc, db)
	templateSvc := template.NewService(templateRepo)
	fileSvc := file.NewService(fileRepo, s3Client, cfg.Storage)
	authSvc := services.NewAuthService(userRepo, cfg)

	// Створюємо хендлери
	authHandler := handlers.NewAuthHandler(authSvc)
	dashboardHandler := handlers.NewDashboardHandler(userSvc, bookingSvc, templateSvc, fileSvc)
	bookingHandler := handlers.NewBookingHandler(bookingSvc, userSvc)
	templateHandler := handlers.NewTemplateHandler(templateSvc)
	fileHandler := handlers.NewFileHandler(fileSvc)
	userHandler := handlers.NewUserHandler(userSvc)

	// Створюємо роутер
	r := router.New(
		authHandler,
		dashboardHandler,
		bookingHandler,
		templateHandler,
		fileHandler,
		userHandler,
	)

	// Налаштовуємо маршрути
	r.SetupRoutes()

	// Запускаємо сервер
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
