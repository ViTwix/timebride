package router

import (
	"log"
	"timebride/internal/handlers"
	"timebride/internal/middleware"
	"timebride/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

type Router struct {
	app              *fiber.App
	sessionStore     *session.Store
	authHandler      *handlers.AuthHandler
	dashboardHandler *handlers.DashboardHandler
	bookingHandler   *handlers.BookingHandler
	templateHandler  *handlers.TemplateHandler
	fileHandler      *handlers.FileHandler
	userHandler      *handlers.UserHandler
}

func New(
	authHandler *handlers.AuthHandler,
	dashboardHandler *handlers.DashboardHandler,
	bookingHandler *handlers.BookingHandler,
	templateHandler *handlers.TemplateHandler,
	fileHandler *handlers.FileHandler,
	userHandler *handlers.UserHandler,
) *Router {
	// Ініціалізуємо HTML шаблонізатор
	engine := html.New("./web/src/templates", ".html")
	engine.Reload(true) // Enable reload in development

	// Додаємо helper функції до шаблонізатора
	engine.AddFuncMap(utils.TemplateFunctions())

	// Створюємо додаток Fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Ініціалізуємо сесії
	store := session.New()

	// Додаємо middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(cors.New())

	// Статичні файли
	app.Static("/css", "./web/public/css")
	app.Static("/js", "./web/public/js")
	app.Static("/img", "./web/public/img")

	return &Router{
		app:              app,
		sessionStore:     store,
		authHandler:      authHandler,
		dashboardHandler: dashboardHandler,
		bookingHandler:   bookingHandler,
		templateHandler:  templateHandler,
		fileHandler:      fileHandler,
		userHandler:      userHandler,
	}
}

func (r *Router) SetupRoutes() {
	// Публічні маршрути
	r.app.Get("/login", r.authHandler.ShowLoginPage)
	r.app.Post("/login", r.authHandler.HandleLogin)
	r.app.Get("/register", r.authHandler.ShowRegisterPage)
	r.app.Post("/register", r.authHandler.HandleRegister)

	// Захищені маршрути
	app := r.app.Group("/")

	// Передаємо секретний ключ для перевірки JWT токена
	app.Use(func(c *fiber.Ctx) error {
		secretKey := r.authHandler.GetJWTSecret()
		if secretKey == "" {
			log.Println("Warning: JWT secret is empty!")
		} else {
			log.Println("JWT secret successfully loaded!")
		}
		c.Locals("jwt_secret", secretKey)
		return c.Next()
	})

	app.Use(middleware.Auth)

	// Дашборд
	app.Get("/", r.dashboardHandler.HandleDashboard)
	app.Get("/dashboard", r.dashboardHandler.HandleDashboard)

	// Бронювання
	app.Get("/bookings", r.bookingHandler.HandleBookingList)
	app.Post("/bookings", r.bookingHandler.HandleBookingCreate)
	app.Put("/bookings/:id", r.bookingHandler.HandleBookingUpdate)
	app.Delete("/bookings/:id", r.bookingHandler.HandleBookingDelete)

	// Шаблони
	app.Get("/templates", r.templateHandler.HandleTemplateList)
	app.Post("/templates", r.templateHandler.HandleTemplateCreate)
	app.Put("/templates/:id", r.templateHandler.HandleTemplateUpdate)
	app.Delete("/templates/:id", r.templateHandler.HandleTemplateDelete)

	// Файли
	app.Get("/files", r.fileHandler.HandleFileList)
	app.Post("/files", r.fileHandler.HandleFileUpload)
	app.Get("/files/:id", r.fileHandler.HandleFileDownload)
	app.Delete("/files/:id", r.fileHandler.HandleFileDelete)

	// Профіль користувача
	app.Get("/profile", r.userHandler.HandleProfile)
	app.Get("/settings", r.userHandler.HandleSettings)
	app.Put("/profile", r.userHandler.HandleUpdateProfile)
	app.Put("/password", r.userHandler.HandleUpdatePassword)
	app.Get("/logout", r.authHandler.HandleLogout)
}

func (r *Router) Start(addr string) error {
	return r.app.Listen(addr)
}
