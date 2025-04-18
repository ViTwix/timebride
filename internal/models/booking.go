package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Booking struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description"`
	Status        string         `json:"status" gorm:"not null;default:'pending'"`
	EventType     string         `json:"event_type" gorm:"not null"`
	StartTime     time.Time      `json:"start_time" gorm:"not null"`
	EndTime       time.Time      `json:"end_time" gorm:"not null"`
	IsAllDay      bool           `json:"is_all_day" gorm:"default:false"`
	ClientName    string         `json:"client_name" gorm:"not null"`
	ClientPhone   string         `json:"client_phone"`
	ClientEmail   string         `json:"client_email"`
	Location      string         `json:"location"`
	Price         float64        `json:"price"`
	Currency      string         `json:"currency" gorm:"default:'USD'"`
	PaymentStatus string         `json:"payment_status" gorm:"default:'pending'"`
	TeamMembers   datatypes.JSON `json:"team_members" gorm:"type:jsonb;default:'[]'"`
	CustomFields  datatypes.JSON `json:"custom_fields" gorm:"type:jsonb;default:'{}'"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// Booking statuses
const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCanceled  = "canceled"
	BookingStatusCompleted = "completed"
	BookingStatusNoShow    = "no_show"
)

// Payment statuses
const (
	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusRefunded = "refunded"
	PaymentStatusCanceled = "canceled"
	PaymentStatusFailed   = "failed"
)

// Event types
const (
	EventTypeConsultation = "consultation"
	EventTypeMeeting      = "meeting"
	EventTypeService      = "service"
	EventTypeOther        = "other"
)

// SetCustomField встановлює значення кастомного поля для бронювання
func (b *Booking) SetCustomField(field CustomField) error {
	var fields map[string]CustomField
	if err := json.Unmarshal(b.CustomFields, &fields); err != nil {
		fields = make(map[string]CustomField)
	}
	fields[field.ID] = field

	data, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("помилка маршалингу кастомних полів: %w", err)
	}

	b.CustomFields = datatypes.JSON(data)
	return nil
}

// GetCustomField повертає значення кастомного поля бронювання
func (b *Booking) GetCustomField(fieldID string) (*CustomField, error) {
	var fields map[string]CustomField
	if err := json.Unmarshal(b.CustomFields, &fields); err != nil {
		return nil, fmt.Errorf("помилка анмаршалингу кастомних полів: %w", err)
	}

	field, ok := fields[fieldID]
	if !ok {
		return nil, fmt.Errorf("поле %s не знайдено", fieldID)
	}

	return &field, nil
}

// AddTeamMember додає учасника команди до бронювання
func (b *Booking) AddTeamMember(member map[string]interface{}) error {
	var members []map[string]interface{}
	if err := json.Unmarshal(b.TeamMembers, &members); err != nil {
		members = make([]map[string]interface{}, 0)
	}

	members = append(members, member)

	data, err := json.Marshal(members)
	if err != nil {
		return fmt.Errorf("помилка маршалингу учасників команди: %w", err)
	}

	b.TeamMembers = datatypes.JSON(data)
	return nil
}

// RemoveTeamMember видаляє учасника команди з бронювання
func (b *Booking) RemoveTeamMember(memberID string) error {
	var members []map[string]interface{}
	if err := json.Unmarshal(b.TeamMembers, &members); err != nil {
		return fmt.Errorf("помилка анмаршалингу учасників команди: %w", err)
	}

	// Створюємо новий слайс без видаленого елемента
	newMembers := make([]map[string]interface{}, 0, len(members))
	for _, member := range members {
		if id, ok := member["id"].(string); ok && id == memberID {
			continue
		}
		newMembers = append(newMembers, member)
	}

	data, err := json.Marshal(newMembers)
	if err != nil {
		return fmt.Errorf("помилка маршалингу учасників команди: %w", err)
	}

	b.TeamMembers = datatypes.JSON(data)
	return nil
}

func (b *Booking) BeforeCreate() error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
