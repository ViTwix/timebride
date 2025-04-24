package booking

import (
	"context"
	"time"

	"github.com/google/uuid"

	"timebride/internal/models"
)

// IBookingService визначає інтерфейс для роботи з бронюваннями
type IBookingService interface {
	// Get отримує бронювання за ID
	Get(ctx context.Context, id uuid.UUID) (*models.Booking, error)

	// Create створює нове бронювання
	Create(ctx context.Context, input *models.BookingCreate) (*models.Booking, error)

	// Update оновлює існуюче бронювання
	Update(ctx context.Context, id uuid.UUID, input *models.BookingUpdate) (*models.Booking, error)

	// Delete видаляє бронювання
	Delete(ctx context.Context, id uuid.UUID) error

	// List отримує список бронювань за фільтром
	List(ctx context.Context, filter map[string]interface{}) ([]*models.Booking, error)

	// GetByDateRange отримує бронювання за діапазоном дат
	GetByDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Booking, error)

	// GetByUserID отримує всі бронювання користувача
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Booking, error)

	// GetByClientID отримує всі бронювання клієнта
	GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*models.Booking, error)

	// CountUpcoming підраховує кількість майбутніх бронювань
	CountUpcoming(ctx context.Context, userID uuid.UUID) (int64, error)

	// CountInDateRange підраховує кількість бронювань в діапазоні дат
	CountInDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) (int64, error)

	// GetRecent отримує останні бронювання
	GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Booking, error)

	// GetByClient отримує всі бронювання клієнта
	GetByClient(ctx context.Context, clientID uuid.UUID) ([]models.Booking, error)
}
