package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository визначає базовий інтерфейс для всіх репозиторіїв
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*T, error)
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

// Repositories містить всі репозиторії програми
type Repositories struct {
	User     UserRepository
	Booking  BookingRepository
	Client   ClientRepository
	Team     TeamRepository
	Price    PriceRepository
	Template TemplateRepository
	File     FileRepository
}

// NewRepositories створює нову структуру репозиторіїв
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:     NewUserRepository(db),
		Booking:  NewBookingRepository(db),
		Client:   NewClientRepository(db),
		Team:     NewTeamRepository(db),
		Price:    NewPriceRepository(db),
		Template: NewTemplateRepository(db),
		File:     NewFileRepository(db),
	}
}

// baseRepository реалізує базові CRUD операції
type baseRepository[T any] struct {
	db *gorm.DB
}

// Create створює новий запис
func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update оновлює існуючий запис
func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete видаляє запис
func (r *baseRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(new(T), "id = ?", id).Error
}

// GetByID отримує запис за ID
func (r *baseRepository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// List отримує список записів за фільтром
func (r *baseRepository[T]) List(ctx context.Context, filter map[string]interface{}) ([]*T, error) {
	var entities []*T
	query := r.db.WithContext(ctx)
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}
	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Count підраховує кількість записів за фільтром
func (r *baseRepository[T]) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T))
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
