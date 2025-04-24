package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// EventType представляє тип події
type EventType string

// Constants for EventType
const (
	EventTypeWedding    EventType = "wedding"
	EventTypeEngagement EventType = "engagement"
	EventTypeCorporate  EventType = "corporate"
	EventTypeFamily     EventType = "family"
	EventTypePortrait   EventType = "portrait"
	EventTypeCommercial EventType = "commercial"
	EventTypeOther      EventType = "other"
)

func (et EventType) IsValid() bool {
	switch et {
	case EventTypeWedding, EventTypePortrait, EventTypeCorporate, EventTypeOther:
		return true
	default:
		return false
	}
}

func NewEventType(value string) (EventType, error) {
	et := EventType(value)
	if !et.IsValid() {
		return "", fmt.Errorf("invalid event type: %s", value)
	}
	return et, nil
}

// BookingStatus визначає статус бронювання
type BookingStatus string

const (
	BookingStatusDraft     BookingStatus = "draft"
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusBooked    BookingStatus = "booked"
	BookingStatusEditing   BookingStatus = "editing"
	BookingStatusReady     BookingStatus = "ready"
	BookingStatusArchived  BookingStatus = "archived"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusDone      BookingStatus = "done"
)

func (bs BookingStatus) IsValid() bool {
	switch bs {
	case BookingStatusDraft, BookingStatusPending, BookingStatusBooked, BookingStatusEditing, BookingStatusReady, BookingStatusArchived, BookingStatusCancelled, BookingStatusDone:
		return true
	default:
		return false
	}
}

func NewBookingStatus(value string) (BookingStatus, error) {
	bs := BookingStatus(value)
	if !bs.IsValid() {
		return "", fmt.Errorf("invalid booking status: %s", value)
	}
	return bs, nil
}

// PaymentStatus defines the payment status of a booking
type PaymentStatus string

// Constants for PaymentStatus
const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPartial  PaymentStatus = "partial"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// TeamPayment представляє оплату члену команди
type TeamPayment struct {
	TeamMemberID uuid.UUID `json:"team_member_id"`
	Amount       float64   `json:"amount"`
}

// Booking представляє бронювання в системі
type Booking struct {
	ID              uuid.UUID      `json:"id" gorm:"primarykey;type:uuid"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	ClientID        uuid.UUID      `json:"client_id" gorm:"type:uuid;not null"`
	Title           string         `json:"title"`
	EventType       EventType      `json:"event_type"`
	EventDate       time.Time      `json:"event_date"`
	Status          BookingStatus  `json:"status"`
	PaymentStatus   PaymentStatus  `json:"payment_status"`
	StartTime       time.Time      `json:"start_time"`
	EndTime         time.Time      `json:"end_time"`
	Description     string         `json:"description"`
	Location        string         `json:"location"`
	PackageName     string         `json:"package_name"`
	DeadlineDays    int            `json:"deadline_days"`
	PriceTotal      float64        `json:"price_total"`
	PriceExtra      float64        `json:"price_extra"`
	PricePrepayment float64        `json:"price_prepayment"`
	TeamMembers     datatypes.JSON `json:"team_members"`
	TeamPayments    datatypes.JSON `json:"team_payments"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       *time.Time     `json:"deleted_at,omitempty" gorm:"index"`

	// Зв'язки
	User   *User   `json:"-" gorm:"foreignKey:UserID"`
	Client *Client `json:"client" gorm:"foreignKey:ClientID"`
}

// BookingPublic представляє публічний вигляд бронювання
type BookingPublic struct {
	ID           uuid.UUID     `json:"id"`
	Title        string        `json:"title"`
	EventType    EventType     `json:"event_type"`
	EventDate    time.Time     `json:"event_date"`
	Location     string        `json:"location"`
	Status       BookingStatus `json:"status"`
	PackageName  string        `json:"package_name"`
	DeadlineDays int           `json:"deadline_days"`
}

// ToPublic конвертує Booking в BookingPublic
func (b *Booking) ToPublic() BookingPublic {
	return BookingPublic{
		ID:           b.ID,
		Title:        b.Title,
		EventType:    b.EventType,
		EventDate:    b.EventDate,
		Location:     b.Location,
		Status:       b.Status,
		PackageName:  b.PackageName,
		DeadlineDays: b.DeadlineDays,
	}
}

// CalculateProfit обчислює прибуток від бронювання
func (b *Booking) CalculateProfit() float64 {
	var teamPaymentsTotal float64
	var payments []TeamPayment
	_ = json.Unmarshal(b.TeamPayments, &payments)
	for _, payment := range payments {
		teamPaymentsTotal += payment.Amount
	}
	return b.PriceTotal - teamPaymentsTotal - b.PriceExtra
}

// CalculateLeftToPay обчислює суму, яку залишилось сплатити
func (b *Booking) CalculateLeftToPay() float64 {
	return b.PriceTotal - b.PricePrepayment
}

// BookingCreate структура для створення бронювання
type BookingCreate struct {
	UserID       string         `json:"user_id" validate:"required,uuid"`
	ClientID     string         `json:"client_id" validate:"required,uuid"`
	Title        string         `json:"title" validate:"required"`
	EventType    EventType      `json:"event_type" validate:"required"`
	EventDate    time.Time      `json:"event_date" validate:"required"`
	StartTime    time.Time      `json:"start_time" validate:"required"`
	EndTime      time.Time      `json:"end_time" validate:"required,gtfield=StartTime"`
	Amount       float64        `json:"amount" validate:"required,gte=0"`
	Prepayment   float64        `json:"prepayment" validate:"gte=0"`
	Currency     string         `json:"currency" validate:"required"`
	Description  string         `json:"description"`
	Location     string         `json:"location"`
	PackageName  string         `json:"package_name"`
	DeadlineDays int            `json:"deadline_days"`
	TeamMembers  datatypes.JSON `json:"team_members"`
	CustomFields datatypes.JSON `json:"custom_fields"`
}

// BookingUpdate структура для оновлення бронювання
type BookingUpdate struct {
	Title        *string         `json:"title,omitempty"`
	EventType    *EventType      `json:"event_type,omitempty"`
	EventDate    *time.Time      `json:"event_date,omitempty"`
	StartTime    *time.Time      `json:"start_time,omitempty"`
	EndTime      *time.Time      `json:"end_time,omitempty"`
	Status       *BookingStatus  `json:"status,omitempty"`
	Amount       *float64        `json:"amount,omitempty"`
	Prepayment   *float64        `json:"prepayment,omitempty"`
	Currency     *string         `json:"currency,omitempty"`
	Description  *string         `json:"description,omitempty"`
	Location     *string         `json:"location,omitempty"`
	PackageName  *string         `json:"package_name,omitempty"`
	DeadlineDays *int            `json:"deadline_days,omitempty"`
	TeamMembers  *datatypes.JSON `json:"team_members,omitempty"`
	CustomFields *datatypes.JSON `json:"custom_fields,omitempty"`
}

// BookingFilter представляє фільтр для пошуку бронювань
type BookingFilter struct {
	UserID    uuid.UUID     `json:"user_id,omitempty"`
	ClientID  uuid.UUID     `json:"client_id,omitempty"`
	Status    BookingStatus `json:"status,omitempty"`
	StartDate time.Time     `json:"start_date,omitempty"`
	EndDate   time.Time     `json:"end_date,omitempty"`
	Search    string        `json:"search,omitempty"`
	SortBy    string        `json:"sort_by,omitempty"`
	SortDesc  bool          `json:"sort_desc,omitempty"`
	Page      int           `json:"page,omitempty"`
	PageSize  int           `json:"page_size,omitempty"`
}

// Validate перевіряє коректність даних бронювання
func (b *Booking) Validate() error {
	if b.UserID == uuid.Nil {
		return ErrValidation{Field: "user_id", Message: "User ID is required"}
	}
	if b.ClientID == uuid.Nil {
		return ErrValidation{Field: "client_id", Message: "Client ID is required"}
	}
	if b.Title == "" {
		return ErrValidation{Field: "title", Message: "Title is required"}
	}
	if b.EventType == "" {
		return ErrValidation{Field: "event_type", Message: "Event type is required"}
	}
	if b.StartTime.After(b.EndTime) {
		return ErrValidation{Field: "end_time", Message: "End time must be after start time"}
	}
	return nil
}

// BeforeCreate - GORM hook for generating UUID before creation
func (b *Booking) BeforeCreate() error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
