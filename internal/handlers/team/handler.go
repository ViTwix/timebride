package team

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/services/team"
)

// Handler обробляє запити для роботи з командою
type Handler struct {
	teamService team.ITeamService
}

// NewHandler створює новий обробник команди
func NewHandler(teamService team.ITeamService) *Handler {
	return &Handler{
		teamService: teamService,
	}
}

// List повертає список членів команди
func (h *Handler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	members, err := h.teamService.ListMembers(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch team members",
		})
	}

	return c.JSON(fiber.Map{
		"members": members,
	})
}

// Create створює нового члена команди
func (h *Handler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var member models.TeamMember
	if err := c.BodyParser(&member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	if err := h.teamService.CreateMember(c.Context(), userUUID, &member); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create team member",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(member)
}

// Get повертає члена команди за ID
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid member ID",
		})
	}

	member, err := h.teamService.GetMember(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Member not found",
		})
	}

	return c.JSON(member)
}

// Update оновлює члена команди
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid member ID",
		})
	}

	var member models.TeamMember
	if err := c.BodyParser(&member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	member.ID = id
	if err := h.teamService.UpdateMember(c.Context(), &member); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update team member",
		})
	}

	return c.JSON(member)
}

// Delete видаляє члена команди
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid member ID",
		})
	}

	if err := h.teamService.DeleteMember(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete team member",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
