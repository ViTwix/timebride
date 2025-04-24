package booking

import (
	"context"
	"time"

	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/repositories"
)

// Service реалізує інтерфейс IBookingService
type Service struct {
	bookingRepo repositories.BookingRepository
	clientRepo  repositories.ClientRepository
}

// NewService створює новий екземпляр сервісу бронювань
func NewService(
	bookingRepo repositories.BookingRepository,
	clientRepo repositories.ClientRepository,
) IBookingService {
	return &Service{
		bookingRepo: bookingRepo,
		clientRepo:  clientRepo,
	}
}

// Get отримує бронювання за ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	return s.bookingRepo.GetByID(ctx, id)
}

// Create створює нове бронювання
func (s *Service) Create(ctx context.Context, input *models.BookingCreate) (*models.Booking, error) {
	userUUID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, err
	}

	clientUUID, err := uuid.Parse(input.ClientID)
	if err != nil {
		return nil, err
	}

	booking := &models.Booking{
		ID:           uuid.New(),
		UserID:       userUUID,
		ClientID:     clientUUID,
		Title:        input.Title,
		EventType:    input.EventType,
		EventDate:    input.EventDate,
		StartTime:    input.StartTime,
		EndTime:      input.EndTime,
		Status:       models.BookingStatusDraft,
		Location:     input.Location,
		Description:  input.Description,
		PackageName:  input.PackageName,
		DeadlineDays: input.DeadlineDays,
		PriceTotal:   input.Amount,
		TeamMembers:  input.TeamMembers,
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

// Update оновлює існуюче бронювання
func (s *Service) Update(ctx context.Context, id uuid.UUID, input *models.BookingUpdate) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		booking.Title = *input.Title
	}
	if input.EventType != nil {
		booking.EventType = *input.EventType
	}
	if input.StartTime != nil {
		booking.StartTime = *input.StartTime
	}
	if input.EndTime != nil {
		booking.EndTime = *input.EndTime
	}
	if input.Status != nil {
		booking.Status = *input.Status
	}
	if input.Location != nil {
		booking.Location = *input.Location
	}
	if input.Description != nil {
		booking.Description = *input.Description
	}
	if input.Amount != nil {
		booking.PriceTotal = *input.Amount
	}
	if input.PackageName != nil {
		booking.PackageName = *input.PackageName
	}
	if input.DeadlineDays != nil {
		booking.DeadlineDays = *input.DeadlineDays
	}
	if input.TeamMembers != nil {
		booking.TeamMembers = *input.TeamMembers
	}

	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

// Delete видаляє бронювання
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.bookingRepo.Delete(ctx, id)
}

// List отримує список бронювань за фільтром
func (s *Service) List(ctx context.Context, filter map[string]interface{}) ([]*models.Booking, error) {
	return s.bookingRepo.List(ctx, filter)
}

// GetByDateRange отримує бронювання за діапазоном дат
func (s *Service) GetByDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Booking, error) {
	return s.bookingRepo.GetByDateRange(ctx, userID, start, end)
}

// GetByUserID отримує всі бронювання користувача
func (s *Service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Booking, error) {
	return s.bookingRepo.GetByUserID(ctx, userID)
}

// GetByClientID отримує всі бронювання клієнта
func (s *Service) GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*models.Booking, error) {
	return s.bookingRepo.GetByClientID(ctx, clientID)
}

// CountUpcoming підраховує кількість майбутніх бронювань
func (s *Service) CountUpcoming(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.bookingRepo.CountUpcoming(ctx, userID)
}

// CountInDateRange підраховує кількість бронювань в діапазоні дат
func (s *Service) CountInDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) (int64, error) {
	return s.bookingRepo.CountInDateRange(ctx, userID, start, end)
}

// GetRecent отримує останні бронювання
func (s *Service) GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Booking, error) {
	return s.bookingRepo.GetRecent(ctx, userID, limit)
}

// GetByClient отримує всі бронювання клієнта
func (s *Service) GetByClient(ctx context.Context, clientID uuid.UUID) ([]models.Booking, error) {
	bookings, err := s.bookingRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	result := make([]models.Booking, len(bookings))
	for i, booking := range bookings {
		result[i] = *booking
	}
	return result, nil
}
