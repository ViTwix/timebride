package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"timebride/internal/models"
)

// ClientRepository визначає інтерфейс для роботи з клієнтами
type ClientRepository interface {
	Repository[models.Client]
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Client, error)
	GetCategories(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetSources(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type clientRepository struct {
	baseRepository[models.Client]
}

// NewClientRepository створює новий репозиторій клієнтів
func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{
		baseRepository: baseRepository[models.Client]{db: db},
	}
}

// GetByUserID отримує список клієнтів користувача
func (r *clientRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Client, error) {
	var clients []*models.Client
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

// GetCategories отримує всі категорії клієнтів
func (r *clientRepository) GetCategories(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var categories []string
	err := r.db.WithContext(ctx).Model(&models.Client{}).
		Where("user_id = ?", userID).
		Distinct().
		Pluck("category", &categories).
		Error
	return categories, err
}

// GetSources отримує всі джерела клієнтів
func (r *clientRepository) GetSources(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var sources []string
	err := r.db.WithContext(ctx).Model(&models.Client{}).
		Where("user_id = ?", userID).
		Distinct().
		Pluck("source", &sources).
		Error
	return sources, err
}
