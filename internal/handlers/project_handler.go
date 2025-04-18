package handlers

import (
    "github.com/gofiber/fiber/v2"
    "timebride/internal/repositories"
)

type ProjectHandler struct {
    Repo *repositories.ProjectRepository
}

func NewProjectHandler(repo *repositories.ProjectRepository) *ProjectHandler {
    return &ProjectHandler{Repo: repo}
}

func (h *ProjectHandler) GetAll(c *fiber.Ctx) error {
    projects, err := h.Repo.GetAll()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.JSON(projects)
}
