package storage

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/google/uuid"

	"timebride/internal/models"
)

// IStorageService визначає інтерфейс сервісу сховища
type IStorageService interface {
	// UploadFile завантажує файл та повертає його URL
	UploadFile(ctx context.Context, file *multipart.FileHeader, path string) (string, error)

	// DeleteFile видаляє файл за шляхом
	DeleteFile(ctx context.Context, path string) error

	// GetFile отримує файл за ID
	GetFile(ctx context.Context, id uuid.UUID) (*models.File, error)

	// ListFiles отримує список файлів
	ListFiles(ctx context.Context, filter map[string]interface{}) ([]*models.File, error)

	// DownloadFile завантажує файл та повертає його вміст
	DownloadFile(ctx context.Context, path string) (io.ReadCloser, error)

	// GetFileURL повертає URL файлу
	GetFileURL(ctx context.Context, path string) string
}
