package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/services/user"
)

// Handler обробляє запити для роботи з користувачами
type Handler struct {
	userService user.IUserService
}

// NewHandler створює новий обробник користувачів
func NewHandler(userService user.IUserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// List повертає список користувачів
func (h *Handler) List(c *fiber.Ctx) error {
	users, err := h.userService.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Конвертуємо в публічне представлення
	publicUsers := make([]*models.PublicUser, len(users))
	for i, user := range users {
		publicUsers[i] = user.ToPublic()
	}

	return c.JSON(fiber.Map{
		"users": publicUsers,
	})
}

// Get повертає користувача за ID
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.userService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user.ToPublic(),
	})
}

// Update оновлює користувача
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var input struct {
		FullName    string `json:"full_name"`
		CompanyName string `json:"company_name"`
		Phone       string `json:"phone"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	user, err := h.userService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	user.FullName = input.FullName
	user.CompanyName = input.CompanyName
	user.Phone = input.Phone

	if err := h.userService.Update(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": user.ToPublic(),
	})
}

// Delete видаляє користувача
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.userService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
