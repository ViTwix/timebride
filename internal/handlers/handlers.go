package handlers

import (
	"github.com/gofiber/fiber/v2"

	"timebride/internal/handlers/auth"
	"timebride/internal/handlers/booking"
	"timebride/internal/handlers/client"
	"timebride/internal/handlers/interfaces"
	"timebride/internal/handlers/price"
	"timebride/internal/handlers/storage"
	"timebride/internal/handlers/team"
	"timebride/internal/handlers/user"
	"timebride/internal/services"
)

// Handlers містить всі HTTP обробники
type Handlers struct {
	Auth     interfaces.IAuthHandler
	Users    interfaces.IUserHandler
	Bookings interfaces.IBookingHandler
	Clients  interfaces.IClientHandler
	Team     interfaces.ITeamHandler
	Prices   interfaces.IPriceHandler
	Storage  interfaces.IStorageHandler
}

// NewHandlers створює нову структуру обробників
func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Auth:     auth.NewHandler(services.Auth),
		Users:    user.NewHandler(services.User),
		Bookings: booking.NewHandler(services.Booking, services.User),
		Clients:  client.NewHandler(services.Client, services.Booking),
		Team:     team.NewHandler(services.Team),
		Prices:   price.NewHandler(services.Price),
		Storage:  storage.NewHandler(services.Storage),
	}
}

// Home обробляє головну сторінку
func (h *Handlers) Home(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{
		"Title": "TimeBride - CRM для фото та відео підрядників",
	})
}

// Dashboard обробляє сторінку дашборду
func (h *Handlers) Dashboard(c *fiber.Ctx) error {
	return c.Render("dashboard", fiber.Map{
		"Title": "Дашборд",
	})
}

// Calendar обробляє сторінку календаря
func (h *Handlers) Calendar(c *fiber.Ctx) error {
	return c.Render("calendar", fiber.Map{
		"Title": "Календар",
	})
}

// Settings обробляє сторінку налаштувань
func (h *Handlers) Settings(c *fiber.Ctx) error {
	return c.Render("settings", fiber.Map{
		"Title": "Налаштування",
	})
}
