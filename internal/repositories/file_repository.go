package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"timebride/internal/models"
)

type FileRepository interface {
	Create(ctx context.Context, file *models.File) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.File, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.File, error)
	Update(ctx context.Context, file *models.File) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*models.File, error)
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *models.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.File, error) {
	var file models.File
	if err := r.db.WithContext(ctx).First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.File, error) {
	var files []*models.File
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileRepository) Update(ctx context.Context, file *models.File) error {
	return r.db.WithContext(ctx).Save(file).Error
}

func (r *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.File{}, "id = ?", id).Error
}

func (r *fileRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.File{})

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *fileRepository) List(ctx context.Context, filter map[string]interface{}) ([]*models.File, error) {
	var files []*models.File
	query := r.db.WithContext(ctx)

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
