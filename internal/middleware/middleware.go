package middleware

import (
	"github.com/gofiber/fiber/v2"

	"timebride/internal/services"
)

// Middleware містить всі middleware обробники
type Middleware struct {
	services *services.Services
}

// NewMiddleware створює нову структуру middleware
func NewMiddleware(services *services.Services) *Middleware {
	return &Middleware{
		services: services,
	}
}

// Auth перевіряє аутентифікацію користувача
func (m *Middleware) Auth(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		return c.Redirect("/login")
	}

	user, err := m.services.Auth.Verify(c.Context(), token)
	if err != nil {
		return c.Redirect("/login")
	}

	c.Locals("user", user)
	return c.Next()
}
