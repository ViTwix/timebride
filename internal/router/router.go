package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"

	"timebride/internal/handlers"
	"timebride/internal/middleware"
	"timebride/internal/utils"
)

type Router struct {
	app          *fiber.App
	sessionStore *session.Store
	handlers     *handlers.Handlers
}

func New(h *handlers.Handlers) *Router {
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		ExposeHeaders:    "Content-Length, Content-Type",
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	// Статичні файли
	app.Static("/css", "./web/public/css")
	app.Static("/js", "./web/public/js")
	app.Static("/img", "./web/public/img")
	app.Static("/fonts", "./web/public/fonts")

	return &Router{
		app:          app,
		sessionStore: store,
		handlers:     h,
	}
}

func (r *Router) SetupRoutes() {
	// Публічні маршрути
	r.app.Get("/", r.handlers.Home)
	r.app.Get("/login", r.handlers.Auth.ShowLoginPage)
	r.app.Post("/login", r.handlers.Auth.HandleLogin)
	r.app.Get("/register", r.handlers.Auth.ShowRegisterPage)
	r.app.Post("/register", r.handlers.Auth.HandleRegister)

	// OAuth маршрути
	r.app.Get("/oauth/:provider", r.handlers.Auth.OAuthRedirect)
	r.app.Get("/oauth/:provider/callback", r.handlers.Auth.OAuthCallback)

	// Захищені маршрути
	app := r.app.Group("/app")

	// Передаємо секретний ключ для перевірки JWT токена
	app.Use(middleware.Auth)

	// Дашборд (головна сторінка для авторизованих користувачів)
	app.Get("/", r.handlers.Dashboard)
	app.Get("/dashboard", r.handlers.Dashboard)

	// Бронювання
	app.Get("/bookings", r.handlers.Bookings.List)
	app.Post("/bookings", r.handlers.Bookings.Create)
	app.Get("/bookings/:id", r.handlers.Bookings.Get)
	app.Put("/bookings/:id", r.handlers.Bookings.Update)
	app.Delete("/bookings/:id", r.handlers.Bookings.Delete)

	// Клієнти
	app.Get("/clients", r.handlers.Clients.List)
	app.Post("/clients", r.handlers.Clients.Create)
	app.Get("/clients/:id", r.handlers.Clients.Get)
	app.Put("/clients/:id", r.handlers.Clients.Update)
	app.Delete("/clients/:id", r.handlers.Clients.Delete)

	// Команда
	app.Get("/team", r.handlers.Team.List)
	app.Post("/team", r.handlers.Team.Create)
	app.Get("/team/:id", r.handlers.Team.Get)
	app.Put("/team/:id", r.handlers.Team.Update)
	app.Delete("/team/:id", r.handlers.Team.Delete)

	// Ціни
	app.Get("/prices", r.handlers.Prices.List)
	app.Post("/prices", r.handlers.Prices.Create)
	app.Get("/prices/:id", r.handlers.Prices.Get)
	app.Put("/prices/:id", r.handlers.Prices.Update)
	app.Delete("/prices/:id", r.handlers.Prices.Delete)

	// Файли
	app.Get("/storage", r.handlers.Storage.List)
	app.Post("/storage/upload", r.handlers.Storage.Upload)
	app.Get("/storage/:id", r.handlers.Storage.Download)
	app.Delete("/storage/:id", r.handlers.Storage.Delete)

	// Профіль користувача
	app.Get("/profile", r.handlers.Users.Get)
	app.Put("/profile", r.handlers.Users.Update)
	app.Get("/settings", r.handlers.Settings)
}

func (r *Router) Start(addr string) error {
	return r.app.Listen(addr)
}
