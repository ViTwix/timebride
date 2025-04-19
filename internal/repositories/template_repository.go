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

// TemplateRepository handles database operations for templates
type TemplateRepository interface {
	Repository[models.Template]

	// GetByUserID retrieves all templates for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Template, error)

	// GetByEventType retrieves templates by their event type
	GetByEventType(ctx context.Context, userID uuid.UUID, eventType string) ([]*models.Template, error)

	// Count returns the number of templates matching the filter
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

type templateRepository struct {
	db *gorm.DB
}

// NewTemplateRepository creates a new instance of TemplateRepository
func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepository{db: db}
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

func (r *templateRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx)

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Model(&models.Template{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
