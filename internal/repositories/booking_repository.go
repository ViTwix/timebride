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

	// CountUpcoming counts upcoming bookings
	CountUpcoming(ctx context.Context, userID uuid.UUID) (int64, error)

	// CountInDateRange counts bookings in a date range
	CountInDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) (int64, error)

	// GetRecent retrieves recent bookings
	GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Booking, error)

	// GetByClientID retrieves bookings by client ID
	GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*models.Booking, error)
}

type bookingRepository struct {
	baseRepository[models.Booking]
}

// NewBookingRepository creates a new instance of BookingRepository
func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{
		baseRepository: baseRepository[models.Booking]{db: db},
	}
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
		Where("user_id = ? AND start_time BETWEEN ? AND ?", userID, start, end).
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

// CountUpcoming counts upcoming bookings
func (r *bookingRepository) CountUpcoming(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("user_id = ? AND start_time > ?", userID, time.Now()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountInDateRange counts bookings in a date range
func (r *bookingRepository) CountInDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("user_id = ? AND start_time BETWEEN ? AND ?", userID, start, end).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetRecent retrieves recent bookings
func (r *bookingRepository) GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

// GetByClientID retrieves bookings by client ID
func (r *bookingRepository) GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*models.Booking, error) {
	var bookings []*models.Booking
	if err := r.db.WithContext(ctx).Where("client_id = ?", clientID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}
