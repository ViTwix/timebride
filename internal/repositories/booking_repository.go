package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"timebride/internal/models"
)

var (
	ErrBookingNotFound = errors.New("booking not found")
)

// BookingRepository handles database operations for bookings
type BookingRepository interface {
	Repository[models.Booking]

	// GetByUserID retrieves all bookings for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Booking, error)

	// GetByDateRange retrieves bookings within a date range
	GetByDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Booking, error)

	// GetByStatus retrieves bookings by their status
	GetByStatus(ctx context.Context, userID uuid.UUID, status string) ([]*models.Booking, error)

	// GetByEventType retrieves bookings by their event type
	GetByEventType(ctx context.Context, userID uuid.UUID, eventType string) ([]*models.Booking, error)

	// Count returns the number of bookings matching the filter
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

type bookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository creates a new instance of BookingRepository
func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.WithContext(ctx).First(&booking, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Booking, error) {
	var bookings []*models.Booking
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Booking, error) {
	var bookings []*models.Booking
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND start_time >= ? AND end_time <= ?", userID, start, end).
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) GetByStatus(ctx context.Context, userID uuid.UUID, status string) ([]*models.Booking, error) {
	var bookings []*models.Booking
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, status).
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) GetByEventType(ctx context.Context, userID uuid.UUID, eventType string) ([]*models.Booking, error) {
	var bookings []*models.Booking
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND event_type = ?", userID, eventType).
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	return r.db.WithContext(ctx).Save(booking).Error
}

func (r *bookingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Booking{}, "id = ?", id).Error
}

func (r *bookingRepository) List(ctx context.Context, filter map[string]interface{}) ([]*models.Booking, error) {
	var bookings []*models.Booking
	query := r.db.WithContext(ctx)

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Booking{})

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
