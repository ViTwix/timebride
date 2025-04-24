package price

import (
	"context"

	"timebride/internal/models"

	"github.com/google/uuid"
)

// IPriceService defines pricing management service interface
type IPriceService interface {
	CreateTemplate(ctx context.Context, template *models.PriceTemplate) error
	UpdateTemplate(ctx context.Context, template *models.PriceTemplate) error
	DeleteTemplate(ctx context.Context, templateID uuid.UUID) error
	GetTemplate(ctx context.Context, templateID uuid.UUID) (*models.PriceTemplate, error)
	ListTemplates(ctx context.Context, userID uuid.UUID) ([]*models.PriceTemplate, error)
}
