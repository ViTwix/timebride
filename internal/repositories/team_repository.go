package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"timebride/internal/models"
)

// TeamRepository визначає інтерфейс для роботи з командою
type TeamRepository interface {
	Repository[models.TeamMember]
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TeamMember, error)
}

type teamRepository struct {
	baseRepository[models.TeamMember]
}

// NewTeamRepository створює новий репозиторій команди
func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{
		baseRepository: baseRepository[models.TeamMember]{db: db},
	}
}

// GetByUserID отримує список членів команди користувача
func (r *teamRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TeamMember, error) {
	var members []*models.TeamMember
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}
