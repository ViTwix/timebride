package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"timebride/internal/config"
	"timebride/internal/models"
	"timebride/internal/repositories"
)

type storageService struct {
	config      *config.Config
	fileRepo    repositories.FileRepository
	storagePath string
}

// NewStorageService creates a new storage service instance
func NewStorageService(cfg *config.Config, fileRepo repositories.FileRepository) IStorageService {
	return &storageService{
		config:      cfg,
		fileRepo:    fileRepo,
		storagePath: cfg.Storage.Path,
	}
}

// UploadFile завантажує файл та створює запис в БД
func (s *storageService) UploadFile(ctx context.Context, file *multipart.FileHeader, path string) (string, error) {
	// Створюємо директорію, якщо не існує
	fullPath := filepath.Join(s.storagePath, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Відкриваємо файл для запису
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Відкриваємо вхідний файл
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Копіюємо вміст
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return s.GetFileURL(ctx, path), nil
}

// DownloadFile завантажує файл зі сховища
func (s *storageService) DownloadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.storagePath, path)
	return os.Open(fullPath)
}

// DeleteFile видаляє файл зі сховища та БД
func (s *storageService) DeleteFile(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.storagePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// ListFiles повертає список файлів користувача
func (s *storageService) ListFiles(ctx context.Context, filter map[string]interface{}) ([]*models.File, error) {
	return s.fileRepo.List(ctx, filter)
}

// GetFileByID повертає файл за ID
func (s *storageService) GetFile(ctx context.Context, id uuid.UUID) (*models.File, error) {
	return s.fileRepo.GetByID(ctx, id)
}

// GetTotalFiles повертає загальну кількість файлів
func (s *storageService) GetTotalFiles(ctx context.Context) (int64, error) {
	return s.fileRepo.Count(ctx, nil)
}

// GetTotalSize повертає загальний розмір файлів
func (s *storageService) GetTotalSize(ctx context.Context) (int64, error) {
	files, err := s.fileRepo.List(ctx, nil)
	if err != nil {
		return 0, err
	}

	var totalSize int64
	for _, file := range files {
		totalSize += file.Size
	}

	return totalSize, nil
}

// GetFileURL повертає URL файлу
func (s *storageService) GetFileURL(ctx context.Context, path string) string {
	// TODO: Implement proper URL generation based on storage provider
	return fmt.Sprintf("/storage/%s", path)
}

// Допоміжні методи

func (s *storageService) saveFile(src io.Reader, dst string) error {
	// TODO: Реалізувати збереження файлу
	return nil
}

func (s *storageService) openFile(path string) (io.ReadCloser, error) {
	// TODO: Реалізувати відкриття файлу
	return nil, nil
}

func (s *storageService) deleteFile(path string) error {
	// TODO: Реалізувати видалення файлу
	return nil
}
