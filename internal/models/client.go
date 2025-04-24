package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// NotificationPreferences визначає налаштування сповіщень для клієнта
type NotificationPreferences struct {
	Email    bool `json:"email"`
	SMS      bool `json:"sms"`
	Telegram bool `json:"telegram"`
}

// ClientSettings представляє налаштування клієнта
type ClientSettings struct {
	Notifications bool              `json:"notifications"`
	Language      string            `json:"language"`
	Preferences   map[string]string `json:"preferences"`
}

// Client представляє клієнта в системі
type Client struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	FullName     string         `json:"full_name"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	Notes        string         `json:"notes"`
	Settings     datatypes.JSON `json:"settings"`
	Status       string         `json:"status"`
	CustomFields datatypes.JSON `json:"custom_fields"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
	User         *User          `json:"-" gorm:"foreignKey:UserID"`
	Bookings     []*Booking     `json:"-" gorm:"foreignKey:ClientID"`
	Avatar       string         `json:"avatar"`
}

// ClientPublic представляє публічний вигляд клієнта
type ClientPublic struct {
	ID       uuid.UUID       `json:"id"`
	FullName string          `json:"full_name"`
	Phone    string          `json:"phone"`
	Email    string          `json:"email"`
	Settings *ClientSettings `json:"settings"`
}

// ClientListResult представляє результат пагінованого списку клієнтів
type ClientListResult struct {
	Items      []*Client `json:"items"`
	TotalItems int64     `json:"total_items"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
}

// ClientListOptions опції для списку клієнтів
type ClientListOptions struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
	SortBy   string `json:"sort_by"`
	SortDesc bool   `json:"sort_desc"`
	Category string `json:"category"`
	Source   string `json:"source"`
}

// ToPublic конвертує Client в ClientPublic
func (c *Client) ToPublic() ClientPublic {
	settings := &ClientSettings{}
	if c.Settings != nil {
		_ = json.Unmarshal(c.Settings, settings)
	}

	return ClientPublic{
		ID:       c.ID,
		FullName: c.FullName,
		Phone:    c.Phone,
		Email:    c.Email,
		Settings: settings,
	}
}

// MarshalJSON реалізує інтерфейс json.Marshaler
func (c *Client) MarshalJSON() ([]byte, error) {
	type Alias Client
	settings := &ClientSettings{}
	if c.Settings != nil {
		_ = json.Unmarshal(c.Settings, settings)
	}

	return json.Marshal(&struct {
		*Alias
		Settings *ClientSettings `json:"settings"`
	}{
		Alias:    (*Alias)(c),
		Settings: settings,
	})
}

// UnmarshalJSON реалізує інтерфейс json.Unmarshaler
func (c *Client) UnmarshalJSON(data []byte) error {
	type Alias Client
	aux := &struct {
		*Alias
		Settings *ClientSettings `json:"settings"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Settings != nil {
		settingsJSON, err := json.Marshal(aux.Settings)
		if err != nil {
			return err
		}
		c.Settings = settingsJSON
	}

	return nil
}

// ToJSON конвертує ClientSettings в JSON string
func (cs ClientSettings) ToJSON() string {
	jsonBytes, _ := json.Marshal(cs)
	return string(jsonBytes)
}

// Validate перевіряє коректність даних клієнта
func (c *Client) Validate() error {
	if c.FullName == "" {
		return ErrValidation{Field: "full_name", Message: "Full name is required"}
	}
	return nil
}

// BeforeCreate генерує UUID для нового клієнта
func (c *Client) BeforeCreate() error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// ErrNotFound повертається, коли запис не знайдено
var ErrNotFound = errors.New("not found")

// CreateClientInput вхідні дані для створення клієнта
type CreateClientInput struct {
	UserID   uuid.UUID      `json:"user_id" validate:"required"`
	FullName string         `json:"full_name" validate:"required"`
	Email    string         `json:"email" validate:"email"`
	Phone    string         `json:"phone"`
	Notes    string         `json:"notes"`
	Settings ClientSettings `json:"settings"`
}

// UpdateClientInput вхідні дані для оновлення клієнта
type UpdateClientInput struct {
	ID       uuid.UUID      `json:"id" validate:"required"`
	UserID   uuid.UUID      `json:"user_id" validate:"required"`
	FullName string         `json:"full_name" validate:"required"`
	Email    string         `json:"email" validate:"email"`
	Phone    string         `json:"phone"`
	Notes    string         `json:"notes"`
	Settings ClientSettings `json:"settings"`
}
