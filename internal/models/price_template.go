package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PriceTemplate представляє шаблон цін підрядника
type PriceTemplate struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primarykey"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name         string         `json:"name" gorm:"not null"`
	EventType    EventType      `json:"event_type" gorm:"not null"`
	Currency     string         `json:"currency" gorm:"default:'UAH'"`
	Price        float64        `json:"price" gorm:"type:decimal(10,2)"`
	Deposit      float64        `json:"deposit" gorm:"type:decimal(10,2)"`
	Description  string         `json:"description"`
	Duration     time.Duration  `json:"duration"`
	TeamPayments datatypes.JSON `json:"team_payments" gorm:"type:jsonb;default:'[]'"`
	DeadlineDays int            `json:"deadline_days" gorm:"default:180"`
	Settings     datatypes.JSON `json:"settings" gorm:"type:jsonb;default:'{}'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at,omitempty" gorm:"index"`

	// Зв'язки
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// PriceTemplatePublic представляє публічну інформацію про шаблон цін
type PriceTemplatePublic struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	EventType    EventType `json:"event_type"`
	Currency     string    `json:"currency"`
	Price        float64   `json:"price"`
	Deposit      float64   `json:"deposit"`
	Description  string    `json:"description"`
	DeadlineDays int       `json:"deadline_days"`
}

// ToPublic конвертує PriceTemplate в PriceTemplatePublic
func (pt *PriceTemplate) ToPublic() PriceTemplatePublic {
	return PriceTemplatePublic{
		ID:           pt.ID,
		Name:         pt.Name,
		EventType:    pt.EventType,
		Currency:     pt.Currency,
		Price:        pt.Price,
		Deposit:      pt.Deposit,
		Description:  pt.Description,
		DeadlineDays: pt.DeadlineDays,
	}
}

// BeforeCreate - GORM хук для генерації UUID перед створенням
func (pt *PriceTemplate) BeforeCreate() error {
	if pt.ID == uuid.Nil {
		pt.ID = uuid.New()
	}
	return nil
}

// Validate перевіряє коректність даних шаблону ціни
func (pt *PriceTemplate) Validate() error {
	if pt.Name == "" {
		return ErrValidation{Field: "name", Message: "Name is required"}
	}
	if pt.Price < 0 {
		return ErrValidation{Field: "price", Message: "Price cannot be negative"}
	}
	if pt.Deposit < 0 {
		return ErrValidation{Field: "deposit", Message: "Deposit cannot be negative"}
	}
	if pt.Deposit > pt.Price {
		return ErrValidation{Field: "deposit", Message: "Deposit cannot be greater than price"}
	}
	return nil
}
