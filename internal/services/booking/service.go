package booking

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"timebride/internal/models"
	"timebride/internal/repositories"
	"timebride/internal/services/user"
)

// Service handles booking-related business logic
type Service struct {
	repo        repositories.BookingRepository
	userService *user.Service
	db          *gorm.DB
}

// New creates a new booking service
func New(repo repositories.BookingRepository, userService *user.Service, db *gorm.DB) *Service {
	return &Service{
		repo:        repo,
		userService: userService,
		db:          db,
	}
}

// GetByID retrieves a booking by its ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	return s.repo.GetByID(ctx, id)
}

// Create creates a new booking
func (s *Service) Create(ctx context.Context, booking *models.Booking) error {
	return s.repo.Create(ctx, booking)
}

// Update updates an existing booking
func (s *Service) Update(ctx context.Context, booking *models.Booking) error {
	return s.repo.Update(ctx, booking)
}

// Delete deletes a booking by its ID
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves all bookings for a user
func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]*models.Booking, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// GetByDateRange retrieves bookings within a date range
func (s *Service) GetByDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Booking, error) {
	return s.repo.GetByDateRange(ctx, userID, start, end)
}

// UpdateStatus updates the status of a booking
func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	booking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	booking.Status = status
	return s.repo.Update(ctx, booking)
}

// UpdatePaymentStatus updates the payment status of a booking
func (s *Service) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string) error {
	booking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	booking.PaymentStatus = status
	return s.repo.Update(ctx, booking)
}

// CountTotal returns the total number of bookings for a user
func (s *Service) CountTotal(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.repo.Count(ctx, map[string]interface{}{"user_id": userID})
}

// CountActive returns the number of active bookings for a user
func (s *Service) CountActive(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.repo.Count(ctx, map[string]interface{}{
		"user_id": userID,
		"status":  "confirmed",
	})
}

// CountUpcoming returns the number of upcoming bookings for a user
func (s *Service) CountUpcoming(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.repo.Count(ctx, map[string]interface{}{
		"user_id":    userID,
		"status":     "confirmed",
		"start_time": map[string]interface{}{"$gt": time.Now()},
	})
}

// CountThisMonth returns the number of bookings this month for a user
func (s *Service) CountThisMonth(ctx context.Context, userID uuid.UUID) (int64, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	return s.repo.Count(ctx, map[string]interface{}{
		"user_id":    userID,
		"start_time": map[string]interface{}{"$gte": startOfMonth, "$lte": endOfMonth},
	})
}

// GetRecent returns recent bookings for a user
func (s *Service) GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Booking, error) {
	bookings, err := s.repo.List(ctx, map[string]interface{}{
		"user_id": userID,
		"order":   "start_time DESC",
		"limit":   limit,
	})
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

// AddTeamMember adds a team member to a booking
func (s *Service) AddTeamMember(ctx context.Context, bookingID uuid.UUID, member map[string]interface{}) error {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	// Конвертуємо поточні дані команди
	var teamMembers []map[string]interface{}
	if booking.TeamMembers != nil {
		if err := json.Unmarshal(booking.TeamMembers, &teamMembers); err != nil {
			return err
		}
	}

	// Додаємо нового члена команди
	teamMembers = append(teamMembers, member)

	// Зберігаємо оновлені дані
	teamMembersJSON, err := json.Marshal(teamMembers)
	if err != nil {
		return err
	}

	booking.TeamMembers = teamMembersJSON
	return s.repo.Update(ctx, booking)
}

// RemoveTeamMember removes a team member from a booking
func (s *Service) RemoveTeamMember(ctx context.Context, bookingID uuid.UUID, memberID string) error {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	// Конвертуємо поточні дані команди
	var teamMembers []map[string]interface{}
	if err := json.Unmarshal(booking.TeamMembers, &teamMembers); err != nil {
		return err
	}

	// Видаляємо члена команди
	var updatedTeamMembers []map[string]interface{}
	for _, member := range teamMembers {
		if member["id"].(string) != memberID {
			updatedTeamMembers = append(updatedTeamMembers, member)
		}
	}

	// Зберігаємо оновлені дані
	teamMembersJSON, err := json.Marshal(updatedTeamMembers)
	if err != nil {
		return err
	}

	booking.TeamMembers = teamMembersJSON
	return s.repo.Update(ctx, booking)
}
