package client

import (
	"context"
	"mime/multipart"

	"timebride/internal/models"

	"github.com/google/uuid"
)

// IClientService визначає інтерфейс для роботи з клієнтами
type IClientService interface {
	// Get отримує клієнта за ID
	Get(ctx context.Context, id uuid.UUID) (*models.Client, error)

	// Create створює нового клієнта
	Create(ctx context.Context, client *models.Client) error

	// Update оновлює існуючого клієнта
	Update(ctx context.Context, client *models.Client) error

	// Delete видаляє клієнта
	Delete(ctx context.Context, id uuid.UUID) error

	// List отримує список клієнтів за фільтром
	List(ctx context.Context, userID uuid.UUID, page, pageSize int) (*models.ClientListResult, error)

	// GetByUserID отримує всіх клієнтів користувача
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Client, error)

	// GetCategories отримує всі категорії клієнтів
	GetCategories(ctx context.Context, userID uuid.UUID) ([]string, error)

	// GetSources отримує всі джерела клієнтів
	GetSources(ctx context.Context, userID uuid.UUID) ([]string, error)

	ListClients(ctx context.Context, userID uuid.UUID, options models.ClientListOptions) (*models.ClientListResult, error)
	GetClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID) (*models.Client, error)
	CreateClient(ctx context.Context, client *models.Client) (*models.Client, error)
	UpdateClient(ctx context.Context, client *models.Client) (*models.Client, error)
	DeleteClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID) error
	UploadAvatar(ctx context.Context, clientID uuid.UUID, file *multipart.FileHeader) (string, error)
	DeleteAvatar(ctx context.Context, clientID uuid.UUID) error
}
