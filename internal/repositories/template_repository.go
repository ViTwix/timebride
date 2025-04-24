package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"timebride/internal/models"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
)

// TemplateRepository визначає інтерфейс для роботи з шаблонами
type TemplateRepository interface {
	Repository[models.Template]

	// GetByUserID retrieves all templates for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Template, error)

	// GetByEventType retrieves templates by their event type
	GetByEventType(ctx context.Context, userID uuid.UUID, eventType string) ([]*models.Template, error)

	// Count returns the number of templates matching the filter
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)

	// GetByType retrieves all templates of a specific type
	GetByType(ctx context.Context, templateType string) ([]*models.Template, error)
}

type templateRepository struct {
	baseRepository[models.Template]
}

// NewTemplateRepository створює новий репозиторій шаблонів
func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepository{
		baseRepository: baseRepository[models.Template]{db: db},
	}
}

func (r *templateRepository) Create(ctx context.Context, template *models.Template) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *templateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	var template models.Template
	if err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	return &template, nil
}

// GetByUserID отримує шаблони користувача
func (r *templateRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Template, error) {
	var templates []*models.Template
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *templateRepository) GetByEventType(ctx context.Context, userID uuid.UUID, eventType string) ([]*models.Template, error) {
	var templates []*models.Template
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND event_type = ?", userID, eventType).
		Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *templateRepository) Update(ctx context.Context, template *models.Template) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *templateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Template{}, "id = ?", id).Error
}

func (r *templateRepository) List(ctx context.Context, filter map[string]interface{}) ([]*models.Template, error) {
	var templates []*models.Template
	query := r.db.WithContext(ctx)

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// Count повертає кількість шаблонів за фільтром
func (r *templateRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Template{})

	if filter != nil {
		for key, value := range filter {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetByType отримує всі шаблони певного типу
func (r *templateRepository) GetByType(ctx context.Context, templateType string) ([]*models.Template, error) {
	var templates []*models.Template
	if err := r.db.WithContext(ctx).Where("type = ?", templateType).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}
