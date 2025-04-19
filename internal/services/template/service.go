package template

import (
	"context"

	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/repositories"
)

type Service struct {
	repo repositories.TemplateRepository
}

func NewService(repo repositories.TemplateRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, template *models.Template) error {
	template.ID = uuid.New()
	return s.repo.Create(ctx, template)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Template, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) Update(ctx context.Context, template *models.Template) error {
	return s.repo.Update(ctx, template)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// GetTotalTemplates returns the total number of templates
func (s *Service) GetTotalTemplates(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx, map[string]interface{}{})
}

// GetActiveTemplates returns the number of active templates
func (s *Service) GetActiveTemplates(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx, map[string]interface{}{
		"status": "active",
	})
}
