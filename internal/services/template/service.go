package template

import (
	"context"

	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/repositories"
)

type templateService struct {
	templateRepo repositories.Repository[models.Template]
}

// NewTemplateService creates a new template service instance
func NewTemplateService(templateRepo repositories.Repository[models.Template]) ITemplateService {
	return &templateService{
		templateRepo: templateRepo,
	}
}

// Create створює новий шаблон
func (s *templateService) Create(ctx context.Context, template *models.Template) error {
	return s.templateRepo.Create(ctx, template)
}

// GetByID отримує шаблон за ID
func (s *templateService) GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	return s.templateRepo.GetByID(ctx, id)
}

// GetByUserID отримує всі шаблони користувача
func (s *templateService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Template, error) {
	return s.templateRepo.List(ctx, map[string]interface{}{"user_id": userID})
}

// Update оновлює шаблон
func (s *templateService) Update(ctx context.Context, template *models.Template) error {
	return s.templateRepo.Update(ctx, template)
}

// Delete видаляє шаблон
func (s *templateService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.templateRepo.Delete(ctx, id)
}

// GetTotalTemplates повертає загальну кількість шаблонів
func (s *templateService) GetTotalTemplates(ctx context.Context) (int64, error) {
	templates, err := s.templateRepo.List(ctx, nil)
	if err != nil {
		return 0, err
	}
	return int64(len(templates)), nil
}

// GetActiveTemplates повертає кількість активних шаблонів
func (s *templateService) GetActiveTemplates(ctx context.Context) (int64, error) {
	templates, err := s.templateRepo.List(ctx, map[string]interface{}{"is_active": true})
	if err != nil {
		return 0, err
	}
	return int64(len(templates)), nil
}
