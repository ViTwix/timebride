package template

import (
	"context"

	"timebride/internal/models"

	"github.com/google/uuid"
)

// ITemplateService defines template management service interface
type ITemplateService interface {
	Create(ctx context.Context, template *models.Template) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Template, error)
	Update(ctx context.Context, template *models.Template) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalTemplates(ctx context.Context) (int64, error)
	GetActiveTemplates(ctx context.Context) (int64, error)
}
