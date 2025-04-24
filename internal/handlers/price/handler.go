package price

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/services/price"
)

// Handler обробляє запити для роботи з цінами
type Handler struct {
	priceService price.IPriceService
}

// NewHandler створює новий обробник цін
func NewHandler(priceService price.IPriceService) *Handler {
	return &Handler{
		priceService: priceService,
	}
}

// List повертає список прайс-листів
func (h *Handler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	templates, err := h.priceService.ListTemplates(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch price templates",
		})
	}

	return c.JSON(fiber.Map{
		"templates": templates,
	})
}

// Create створює новий прайс-лист
func (h *Handler) Create(c *fiber.Ctx) error {
	var template models.PriceTemplate
	if err := c.BodyParser(&template); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	if err := h.priceService.CreateTemplate(c.Context(), &template); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create price template",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(template)
}

// Get повертає прайс-лист за ID
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID",
		})
	}

	template, err := h.priceService.GetTemplate(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Template not found",
		})
	}

	return c.JSON(template)
}

// Update оновлює прайс-лист
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID",
		})
	}

	var template models.PriceTemplate
	if err := c.BodyParser(&template); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	template.ID = id
	if err := h.priceService.UpdateTemplate(c.Context(), &template); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update price template",
		})
	}

	return c.JSON(template)
}

// Delete видаляє прайс-лист
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID",
		})
	}

	if err := h.priceService.DeleteTemplate(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete price template",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
