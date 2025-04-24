package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"gorm.io/gorm"

	"timebride/internal/config"
	"timebride/internal/db"
	"timebride/internal/handlers"
	"timebride/internal/middleware"
	"timebride/internal/repositories"
	"timebride/internal/services"
	"timebride/internal/services/auth"
	"timebride/internal/services/booking"
	"timebride/internal/services/client"
	"timebride/internal/services/price"
	"timebride/internal/services/storage"
	"timebride/internal/services/team"
	"timebride/internal/services/template"
	"timebride/internal/services/user"
)

// AppModules містить основні модулі програми
type AppModules struct {
	Config     *config.Config
	DB         *gorm.DB
	Templates  *html.Engine
	Static     string
	Handlers   *handlers.Handlers
	Services   *services.Services
	Repos      *repositories.Repositories
	Middleware *middleware.Middleware
}

func main() {
	// Ініціалізуємо модулі
	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Налаштовуємо і запускаємо сервер
	server := setupServer(app)

	// Запускаємо сервер
	go func() {
		if err := server.Listen(app.Config.Server.Address); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Очікуємо сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func initializeApp() (*AppModules, error) {
	// Завантажуємо конфігурацію
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Підключаємося до бази даних
	database, err := db.Connect(cfg.Database)
	if err != nil {
		return nil, err
	}

	// Ініціалізуємо репозиторії
	repos := repositories.NewRepositories(database)

	// Ініціалізуємо сервіси
	storageService := storage.NewStorageService(cfg, repos.File)
	authService := auth.NewAuthService(cfg, repos.User)
	userService := user.NewUserService(repos.User)
	clientService := client.NewService(repos.Client, repos.File, storageService)
	teamService := team.NewTeamService(repos.Team)
	priceService := price.NewPriceService(repos.Price)
	templateService := template.NewTemplateService(repos.Template)
	bookingService := booking.NewService(repos.Booking, repos.Client)

	// Створюємо екземпляр Services
	services := services.NewServices(
		authService,
		userService,
		bookingService,
		clientService,
		teamService,
		priceService,
		storageService,
		templateService,
	)

	// Ініціалізуємо middleware
	middlewares := middleware.NewMiddleware(services)

	// Ініціалізуємо обробники
	handlers := handlers.NewHandlers(services)

	return &AppModules{
		Config:     cfg,
		DB:         database,
		Templates:  initTemplates(),
		Static:     initStatic(),
		Handlers:   handlers,
		Services:   services,
		Repos:      repos,
		Middleware: middlewares,
	}, nil
}

func setupServer(app *AppModules) *fiber.App {
	// Налаштовуємо Fiber
	server := fiber.New(fiber.Config{
		Views:                 app.Templates,
		ViewsLayout:           "layouts/main",
		AppName:               "TimeBride",
		BodyLimit:             50 * 1024 * 1024, // 50MB
		DisableStartupMessage: true,             // Вимикаємо стартове повідомлення
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// Визначаємо тип помилки та відповідний статус код
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Логуємо помилку
			log.Printf("Error: %v, Path: %s, Method: %s", err, c.Path(), c.Method())

			// Якщо це API запит, повертаємо JSON
			if strings.HasPrefix(c.Path(), "/api") {
				return c.Status(code).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			// Для веб-запитів рендеримо сторінку з помилкою
			return c.Status(code).Render("error", fiber.Map{
				"Error": err.Error(),
				"Code":  code,
			})
		},
	})

	// Middleware
	server.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	server.Use(logger.New(logger.Config{
		Format:     "${time} ${status} ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))
	server.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(app.Config.Server.CorsOrigins, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Статичні файли
	server.Static("/static", app.Static, fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		CacheDuration: 24 * time.Hour,
	})

	// Публічні маршрути
	server.Get("/", app.Handlers.Home)
	server.Get("/login", app.Handlers.Auth.ShowLoginPage)
	server.Post("/login", app.Handlers.Auth.HandleLogin)
	server.Get("/register", app.Handlers.Auth.ShowRegisterPage)
	server.Post("/register", app.Handlers.Auth.HandleRegister)

	// API маршрути
	api := server.Group("/api", app.Middleware.Auth)
	{
		// Користувачі
		users := api.Group("/users")
		users.Get("/", app.Handlers.Users.List)
		users.Get("/:id", app.Handlers.Users.Get)
		users.Put("/:id", app.Handlers.Users.Update)
		users.Delete("/:id", app.Handlers.Users.Delete)

		// Бронювання
		bookings := api.Group("/bookings")
		bookings.Get("/", app.Handlers.Bookings.List)
		bookings.Post("/", app.Handlers.Bookings.Create)
		bookings.Get("/:id", app.Handlers.Bookings.Get)
		bookings.Put("/:id", app.Handlers.Bookings.Update)
		bookings.Delete("/:id", app.Handlers.Bookings.Delete)

		// Клієнти
		clients := api.Group("/clients")
		clients.Get("/", app.Handlers.Clients.List)
		clients.Post("/", app.Handlers.Clients.Create)
		clients.Get("/:id", app.Handlers.Clients.Get)
		clients.Put("/:id", app.Handlers.Clients.Update)
		clients.Delete("/:id", app.Handlers.Clients.Delete)

		// Команда
		team := api.Group("/team")
		team.Get("/", app.Handlers.Team.List)
		team.Post("/", app.Handlers.Team.Create)
		team.Get("/:id", app.Handlers.Team.Get)
		team.Put("/:id", app.Handlers.Team.Update)
		team.Delete("/:id", app.Handlers.Team.Delete)

		// Прайс-листи
		prices := api.Group("/prices")
		prices.Get("/", app.Handlers.Prices.List)
		prices.Post("/", app.Handlers.Prices.Create)
		prices.Get("/:id", app.Handlers.Prices.Get)
		prices.Put("/:id", app.Handlers.Prices.Update)
		prices.Delete("/:id", app.Handlers.Prices.Delete)

		// Сховище
		storage := api.Group("/storage")
		storage.Get("/", app.Handlers.Storage.List)
		storage.Post("/upload", app.Handlers.Storage.Upload)
		storage.Get("/:id", app.Handlers.Storage.Download)
		storage.Delete("/:id", app.Handlers.Storage.Delete)
	}

	// Захищені веб-маршрути
	protected := server.Group("/app", app.Middleware.Auth)
	{
		protected.Get("/", app.Handlers.Dashboard)
		protected.Get("/calendar", app.Handlers.Calendar)
		protected.Get("/bookings", app.Handlers.Bookings.List)
		protected.Get("/team", app.Handlers.Team.List)
		protected.Get("/clients", app.Handlers.Clients.List)
		protected.Get("/prices", app.Handlers.Prices.List)
		protected.Get("/storage", app.Handlers.Storage.List)
		protected.Get("/settings", app.Handlers.Settings)
	}

	return server
}

func initTemplates() *html.Engine {
	engine := html.New("./web/templates", ".html")
	engine.Reload(true) // Enable template reloading for development
	return engine
}

func initStatic() string {
	return "./web/static"
}
