package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	"github.com/gofiber/template/html/v2"
)

// AppModules містить основні модулі програми
type AppModules struct {
	Config      *config.Config
	DB          *gorm.DB
	Templates   *html.Engine
	Static      string
	Controllers string
	Public      string
	Handlers    *handlers.Handlers
	Services    *services.Services
	Repos       *repositories.Repositories
	Middleware  *middleware.Middleware
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
	authService := auth.NewAuthService(cfg, repos.User)
	userService := user.NewUserService(repos.User)
	storageService := storage.NewStorageService(cfg, repos.File)
	clientService := client.NewService(repos.Client, repos.File, storageService)
	bookingService := booking.NewService(repos.Booking, repos.Client)
	teamService := team.NewTeamService(repos.Team)
	priceService := price.NewPriceService(repos.Price)
	templateService := template.NewTemplateService(repos.Template)

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

	// Ініціалізуємо шаблонізатор
	templates := initTemplates()

	// Ініціалізуємо статичні файли
	static := initStatic()

	// Ініціалізуємо контролери
	controllers := initControllers()

	// Ініціалізуємо публічні файли
	public := initPublic()

	// Ініціалізуємо middleware
	middleware := middleware.NewMiddleware(services)

	// Ініціалізуємо хендлери
	handlers := handlers.NewHandlers(services)

	return &AppModules{
		Config:      cfg,
		DB:          database,
		Templates:   templates,
		Static:      static,
		Controllers: controllers,
		Public:      public,
		Handlers:    handlers,
		Services:    services,
		Repos:       repos,
		Middleware:  middleware,
	}, nil
}

func setupServer(app *AppModules) *fiber.App {
	log.Println("Setting up server...")

	// Створюємо новий екземпляр Fiber
	server := fiber.New(fiber.Config{
		Views: app.Templates,
	})

	// Налаштовуємо middleware
	server.Use(recover.New())
	server.Use(logger.New())
	server.Use(cors.New())

	// Додаємо app до контексту запиту
	server.Use(func(c *fiber.Ctx) error {
		c.Locals("app", app)
		return c.Next()
	})

	// Налаштовуємо статичні файли
	log.Printf("Static files path: %s", app.Static)
	if _, err := os.Stat(app.Static); os.IsNotExist(err) {
		log.Printf("WARNING: Static directory does not exist: %s", app.Static)
	} else {
		log.Printf("Static directory exists: %s", app.Static)
		// Перевірка наявності файлів в директорії
		files, err := filepath.Glob(filepath.Join(app.Static, "**", "*"))
		if err != nil {
			log.Printf("Error listing static files: %v", err)
		} else {
			log.Printf("Static files: %v", files)
		}
	}

	// Налаштовуємо статичні файли
	server.Static("/static", app.Static)

	// Налаштовуємо публічні файли
	log.Printf("Public files path: %s", app.Public)
	if _, err := os.Stat(app.Public); os.IsNotExist(err) {
		log.Printf("WARNING: Public directory does not exist: %s", app.Public)
	} else {
		log.Printf("Public directory exists: %s", app.Public)
	}
	server.Static("/", app.Public)

	// Налаштовуємо маршрути
	setupRoutes(server, app)

	return server
}

func setupRoutes(app *fiber.App, modules *AppModules) {
	// Публічні маршрути
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("auth/login", fiber.Map{})
	})

	app.Post("/login", modules.Handlers.Auth.HandleLogin)

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("auth/register", fiber.Map{})
	})

	app.Post("/register", modules.Handlers.Auth.HandleRegister)

	app.Get("/logout", modules.Handlers.Auth.HandleLogout)

	// Захищені маршрути
	api := app.Group("/dashboard", modules.Middleware.Auth)
	{
		// Дашборд
		api.Get("/", func(c *fiber.Ctx) error {
			return c.Render("dashboard/index", fiber.Map{})
		})

		// Клієнти
		api.Get("/clients", func(c *fiber.Ctx) error {
			return c.Render("dashboard/clients/index", fiber.Map{})
		})

		api.Post("/clients", modules.Handlers.Clients.Create)
		api.Get("/clients/:id", modules.Handlers.Clients.Get)
		api.Put("/clients/:id", modules.Handlers.Clients.Update)
		api.Delete("/clients/:id", modules.Handlers.Clients.Delete)

		// Бронювання
		api.Get("/bookings", func(c *fiber.Ctx) error {
			return c.Render("dashboard/bookings/index", fiber.Map{})
		})

		api.Post("/bookings", modules.Handlers.Bookings.Create)
		api.Get("/bookings/:id", modules.Handlers.Bookings.Get)
		api.Put("/bookings/:id", modules.Handlers.Bookings.Update)
		api.Delete("/bookings/:id", modules.Handlers.Bookings.Delete)

		// Профіль
		api.Get("/profile", func(c *fiber.Ctx) error {
			return c.Render("dashboard/profile", fiber.Map{})
		})

		api.Put("/profile", modules.Handlers.Users.Update)
	}
}

func initTemplates() *html.Engine {
	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		templateDir = "./web/templates"
	}
	log.Printf("Initializing templates from directory: %s", templateDir)

	engine := html.New(templateDir, ".html")

	// Add debug logging for template loading
	engine.AddFunc("debug", func(v interface{}) string {
		log.Printf("Template debug: %v", v)
		return fmt.Sprintf("%v", v)
	})

	if err := engine.Load(); err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	log.Printf("Templates loaded successfully")
	return engine
}

func initStatic() string {
	// Використовуємо змінну середовища для шляху до статичних файлів
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./web/static" // Значення за замовчуванням
	}
	log.Printf("Static directory: %s", staticDir)
	return staticDir
}

func initControllers() string {
	// Використовуємо змінну середовища для шляху до контролерів
	controllersDir := os.Getenv("CONTROLLERS_DIR")
	if controllersDir == "" {
		controllersDir = "./web/controllers" // Значення за замовчуванням
	}
	log.Printf("Controllers directory: %s", controllersDir)
	return controllersDir
}

func initPublic() string {
	// Використовуємо змінну середовища для шляху до публічних файлів
	publicDir := os.Getenv("PUBLIC_DIR")
	if publicDir == "" {
		publicDir = "./web/public" // Значення за замовчуванням
	}
	log.Printf("Public directory: %s", publicDir)
	return publicDir
}

func renderTemplate(c *fiber.Ctx, template string, data fiber.Map) error {
	app := c.Locals("app").(*AppModules)
	return app.Templates.Render(c, template, data)
}
