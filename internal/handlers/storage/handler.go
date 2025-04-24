package storage

import (
	"timebride/internal/services/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handler обробляє запити для роботи з файлами
type Handler struct {
	storageService storage.IStorageService
}

// NewHandler створює новий обробник файлів
func NewHandler(storageService storage.IStorageService) *Handler {
	return &Handler{
		storageService: storageService,
	}
}

// List повертає список файлів
func (h *Handler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	filter := map[string]interface{}{
		"user_id": userUUID,
	}

	files, err := h.storageService.ListFiles(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch files",
		})
	}

	return c.JSON(fiber.Map{
		"files": files,
	})
}

// Upload завантажує файл
func (h *Handler) Upload(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file provided",
		})
	}

	uploadedFile, err := h.storageService.UploadFile(c.Context(), file, userUUID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(uploadedFile)
}

// Download завантажує файл
func (h *Handler) Download(c *fiber.Ctx) error {
	fileID := c.Params("id")
	id, err := uuid.Parse(fileID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	reader, err := h.storageService.DownloadFile(c.Context(), id.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}
	defer reader.Close()

	return c.SendStream(reader)
}

// Delete видаляє файл
func (h *Handler) Delete(c *fiber.Ctx) error {
	fileID := c.Params("id")
	id, err := uuid.Parse(fileID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	if err := h.storageService.DeleteFile(c.Context(), id.String()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
