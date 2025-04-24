package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"timebride/internal/models"
)

// PriceRepository визначає інтерфейс для роботи з цінами
type PriceRepository interface {
	Repository[models.PriceTemplate]
	GetActive(ctx context.Context) ([]*models.PriceTemplate, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PriceTemplate, error)
}

type priceRepository struct {
	baseRepository[models.PriceTemplate]
}

// NewPriceRepository створює новий репозиторій цін
func NewPriceRepository(db *gorm.DB) PriceRepository {
	return &priceRepository{
		baseRepository: baseRepository[models.PriceTemplate]{db: db},
	}
}

// GetByUserID отримує всі шаблони цін користувача
func (r *priceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.PriceTemplate, error) {
	var templates []*models.PriceTemplate
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// GetActive отримує всі активні шаблони цін
func (r *priceRepository) GetActive(ctx context.Context) ([]*models.PriceTemplate, error) {
	var templates []*models.PriceTemplate
	if err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}
