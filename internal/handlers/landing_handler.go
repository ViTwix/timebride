package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// LandingHandler обробляє запити лендінг-сторінки
type LandingHandler struct{}

// NewLandingHandler створює новий обробник лендінг-сторінки
func NewLandingHandler() *LandingHandler {
	return &LandingHandler{}
}

// HandleLanding відображає лендінг-сторінку
func (h *LandingHandler) HandleLanding(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "TimeBride - Інтелектуальна CRM для фотографів та відеографів",
	})
}
