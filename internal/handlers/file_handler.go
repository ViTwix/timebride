package handlers

import (
	"fmt"
	"path/filepath"

	"timebride/internal/services/file"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// FileHandler обробляє HTTP запити для роботи з файлами
type FileHandler struct {
	fileService *file.Service
}

// NewFileHandler створює новий екземпляр FileHandler
func NewFileHandler(fileService *file.Service) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// HandleFileList displays the list of files
func (h *FileHandler) HandleFileList(c *fiber.Ctx) error {
	// Get user ID from context
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get user's files
	files, err := h.fileService.ListFiles(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get files",
		})
	}

	// Render template
	return c.Render("file/list", fiber.Map{
		"Title": "My Files",
		"Files": files,
	})
}

// HandleFileUpload handles file upload requests
func (h *FileHandler) HandleFileUpload(c *fiber.Ctx) error {
	// Get user ID from context
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Upload file
	uploadedFile, err := h.fileService.UploadFile(c.Context(), file, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	return c.JSON(uploadedFile)
}

// HandleFileGet handles file info retrieval
func (h *FileHandler) HandleFileGet(c *fiber.Ctx) error {
	// Get user ID from context
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get file ID from URL
	fileID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	// Get file info
	file, err := h.fileService.GetFileByID(c.Context(), fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Check if file belongs to user
	if file.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	return c.JSON(file)
}

// HandleFileDownload handles file download requests
func (h *FileHandler) HandleFileDownload(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file ID is required",
		})
	}

	id, err := uuid.Parse(fileID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid file ID",
		})
	}

	file, err := h.fileService.GetFileByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "file not found",
		})
	}

	reader, filename, err := h.fileService.DownloadFile(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to download file",
		})
	}
	defer reader.Close()

	c.Set("Content-Type", file.FileType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filename)))

	return c.SendStream(reader)
}

// HandleFileDelete handles file deletion requests
func (h *FileHandler) HandleFileDelete(c *fiber.Ctx) error {
	// Get user ID from context
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get file ID from URL
	fileID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	// Get file info first to check ownership
	file, err := h.fileService.GetFileByID(c.Context(), fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Check if file belongs to user
	if file.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Delete file
	if err := h.fileService.DeleteFile(c.Context(), fileID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	return c.JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}
