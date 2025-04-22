package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Статуси оплати
const (
	PaymentStatusPending       = "pending"
	PaymentStatusPaid          = "paid"
	PaymentStatusPartial       = "partial"
	PaymentStatusAfterDelivery = "after_delivery"
)

// Статуси бронювання
const (
	BookingStatusActive    = "active"
	BookingStatusCompleted = "completed"
	BookingStatusArchived  = "archived"
	BookingStatusCancelled = "cancelled"
	BookingStatusPending   = "pending"
)

// Типи подій
const (
	EventTypeWedding    = "wedding"
	EventTypePortrait   = "portrait"
	EventTypeFamily     = "family"
	EventTypeEvent      = "event"
	EventTypeCommercial = "commercial"
)

// Booking представляє бронювання зйомки
type Booking struct {
	ID              uuid.UUID      `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID          uuid.UUID      `json:"user_id" db:"user_id" gorm:"type:uuid;not null"`
	Title           string         `json:"title" db:"title" gorm:"not null"`
	EventType       string         `json:"event_type" db:"event_type" gorm:"default:'event'"`
	EventDate       time.Time      `json:"event_date" db:"event_date" gorm:"not null"`
	StartTime       time.Time      `json:"start_time" db:"start_time" gorm:"not null"`
	EndTime         time.Time      `json:"end_time" db:"end_time" gorm:"not null"`
	Location        *string        `json:"location,omitempty" db:"location"`
	ClientName      string         `json:"client_name" db:"client_name" gorm:"not null"`
	ClientPhone     *string        `json:"client_phone,omitempty" db:"client_phone"`
	ClientEmail     string         `json:"client_email" db:"client_email" gorm:"not null"`
	InstagramHandle *string        `json:"instagram_handle,omitempty" db:"instagram_handle"`
	ContractURL     *string        `json:"contract_url,omitempty" db:"contract_url"`
	TotalPrice      float64        `json:"total_price" db:"total_price" gorm:"default:0"`
	Deposit         float64        `json:"deposit" db:"deposit" gorm:"default:0"`
	PaymentStatus   string         `json:"payment_status" db:"payment_status" gorm:"default:'pending'"`
	Status          string         `json:"status" db:"status" gorm:"default:'pending'"`
	Notes           *string        `json:"notes,omitempty" db:"notes"`
	TeamMembers     datatypes.JSON `json:"team_members" db:"team_members" gorm:"type:jsonb;default:'[]'"`
	CustomFields    datatypes.JSON `json:"custom_fields,omitempty" db:"custom_fields" gorm:"type:jsonb;default:'{}'"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// BookingCreate представляє дані для створення бронювання
type BookingCreate struct {
	UserID          uuid.UUID      `json:"user_id" db:"user_id"`
	Title           string         `json:"title" db:"title"`
	EventType       string         `json:"event_type" db:"event_type"`
	EventDate       time.Time      `json:"event_date" db:"event_date"`
	StartTime       time.Time      `json:"start_time" db:"start_time"`
	EndTime         time.Time      `json:"end_time" db:"end_time"`
	Location        *string        `json:"location,omitempty" db:"location"`
	ClientName      string         `json:"client_name" db:"client_name"`
	ClientPhone     *string        `json:"client_phone,omitempty" db:"client_phone"`
	ClientEmail     string         `json:"client_email" db:"client_email"`
	InstagramHandle *string        `json:"instagram_handle,omitempty" db:"instagram_handle"`
	ContractURL     *string        `json:"contract_url,omitempty" db:"contract_url"`
	TotalPrice      float64        `json:"total_price" db:"total_price"`
	Deposit         float64        `json:"deposit" db:"deposit"`
	PaymentStatus   string         `json:"payment_status" db:"payment_status"`
	Status          string         `json:"status" db:"status"`
	Notes           *string        `json:"notes,omitempty" db:"notes"`
	TeamMembers     datatypes.JSON `json:"team_members,omitempty" db:"team_members"`
	CustomFields    datatypes.JSON `json:"custom_fields,omitempty" db:"custom_fields"`
}

// BookingUpdate представляє дані для оновлення бронювання
type BookingUpdate struct {
	Title           *string        `json:"title,omitempty" db:"title"`
	EventType       *string        `json:"event_type,omitempty" db:"event_type"`
	EventDate       *time.Time     `json:"event_date,omitempty" db:"event_date"`
	StartTime       *time.Time     `json:"start_time,omitempty" db:"start_time"`
	EndTime         *time.Time     `json:"end_time,omitempty" db:"end_time"`
	Location        *string        `json:"location,omitempty" db:"location"`
	ClientName      *string        `json:"client_name,omitempty" db:"client_name"`
	ClientPhone     *string        `json:"client_phone,omitempty" db:"client_phone"`
	ClientEmail     *string        `json:"client_email,omitempty" db:"client_email"`
	InstagramHandle *string        `json:"instagram_handle,omitempty" db:"instagram_handle"`
	ContractURL     *string        `json:"contract_url,omitempty" db:"contract_url"`
	TotalPrice      *float64       `json:"total_price,omitempty" db:"total_price"`
	Deposit         *float64       `json:"deposit,omitempty" db:"deposit"`
	PaymentStatus   *string        `json:"payment_status,omitempty" db:"payment_status"`
	Status          *string        `json:"status,omitempty" db:"status"`
	Notes           *string        `json:"notes,omitempty" db:"notes"`
	TeamMembers     datatypes.JSON `json:"team_members,omitempty" db:"team_members"`
	CustomFields    datatypes.JSON `json:"custom_fields,omitempty" db:"custom_fields"`
}

// BookingWithUser представляє бронювання з даними про користувача
type BookingWithUser struct {
	Booking
	User UserPublic `json:"user"`
}

// BeforeCreate - GORM хук для генерації UUID перед створенням
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
