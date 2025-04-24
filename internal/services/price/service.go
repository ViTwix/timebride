package price

import (
	"context"

	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/repositories"
)

type priceService struct {
	priceRepo repositories.Repository[models.PriceTemplate]
}

// NewPriceService creates a new price service instance
func NewPriceService(priceRepo repositories.Repository[models.PriceTemplate]) IPriceService {
	return &priceService{
		priceRepo: priceRepo,
	}
}

func (s *priceService) CreateTemplate(ctx context.Context, template *models.PriceTemplate) error {
	return s.priceRepo.Create(ctx, template)
}

func (s *priceService) UpdateTemplate(ctx context.Context, template *models.PriceTemplate) error {
	return s.priceRepo.Update(ctx, template)
}

func (s *priceService) DeleteTemplate(ctx context.Context, templateID uuid.UUID) error {
	return s.priceRepo.Delete(ctx, templateID)
}

func (s *priceService) GetTemplate(ctx context.Context, templateID uuid.UUID) (*models.PriceTemplate, error) {
	return s.priceRepo.GetByID(ctx, templateID)
}

func (s *priceService) ListTemplates(ctx context.Context, userID uuid.UUID) ([]*models.PriceTemplate, error) {
	return s.priceRepo.List(ctx, map[string]interface{}{"user_id": userID})
}
