package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"timebride/internal/config"
	"timebride/internal/handlers"
	"timebride/internal/middleware"
	"timebride/internal/repositories"
	"timebride/internal/router"
	"timebride/internal/services/auth"
	"timebride/internal/services/booking"
	"timebride/internal/services/file"
	"timebride/internal/services/template"
	"timebride/internal/services/user"
	"timebride/pkg/database"
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
	authSvc := auth.NewService(userRepo, cfg.JWT)

	// Створюємо конфігурацію авторизації
	authConfig := middleware.AuthConfig{
		SecretKey:     cfg.JWT.SecretKey,
		TokenDuration: cfg.JWT.TokenDuration,
	}

	// Створюємо хендлери
	authHandler := handlers.NewAuthHandler(userSvc, authConfig, authSvc)
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
