package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Типи аутентифікації користувачів
const (
	AuthProviderLocal    = "local"
	AuthProviderGoogle   = "google"
	AuthProviderApple    = "apple"
	AuthProviderFacebook = "facebook"
)

// Ролі користувачів
const (
	RoleOwner     = "owner"
	RoleAdmin     = "admin"
	RoleRetoucher = "retoucher"
	RoleEditor    = "editor"
	RoleOperator  = "operator"
	RoleAssistant = "assistant"
	RoleUser      = "user"
)

// Permissions - типи дозволів для користувача
type Permissions struct {
	ViewFinancials bool `json:"view_financials"`
	EditProjects   bool `json:"edit_projects"`
	ViewProjects   bool `json:"view_projects"`
	ManageTeam     bool `json:"manage_team"`
}

// User представляє користувача системи
type User struct {
	ID              uuid.UUID      `json:"id" db:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email           string         `json:"email" db:"email" gorm:"uniqueIndex;not null"`
	PasswordHash    string         `json:"-" db:"password_hash"`
	AuthProvider    string         `json:"auth_provider" db:"auth_provider" gorm:"default:'local'"`
	ProviderID      *string        `json:"provider_id,omitempty" db:"provider_id"`
	Name            string         `json:"name" db:"name" gorm:"not null"`
	AvatarURL       *string        `json:"avatar_url,omitempty" db:"avatar_url"`
	IsEmailVerified bool           `json:"is_email_verified" db:"is_email_verified" gorm:"default:false"`
	Is2FAEnabled    bool           `json:"is_2fa_enabled" db:"is_2fa_enabled" gorm:"default:false"`
	TwoFASecret     *string        `json:"-" db:"two_fa_secret"`
	Phone           *string        `json:"phone,omitempty" db:"phone"`
	Country         *string        `json:"country,omitempty" db:"country"`
	City            *string        `json:"city,omitempty" db:"city"`
	Role            string         `json:"role" db:"role" gorm:"default:'user'"`
	Language        string         `json:"language" db:"language" gorm:"default:'uk'"`
	Settings        datatypes.JSON `json:"settings,omitempty" gorm:"type:jsonb;default:'{}'"`
	ParentID        *uuid.UUID     `json:"parent_id,omitempty" db:"parent_id" gorm:"type:uuid"`
	Permissions     datatypes.JSON `json:"permissions" gorm:"type:jsonb;default:'{}'"`
	CustomFields    datatypes.JSON `json:"custom_fields,omitempty" gorm:"type:jsonb;default:'{}'"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// Settings представляє кастомні налаштування користувача
type Settings struct {
	Notifications *NotificationSettings `json:"notifications,omitempty"`
	Preferences   *UserPreferences      `json:"preferences,omitempty"`
	Theme         *ThemeSettings        `json:"theme,omitempty"`
}

// Value - реалізація інтерфейсу Valuer для Settings
func (s Settings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan - реалізація інтерфейсу Scanner для Settings
func (s *Settings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &s)
}

// NotificationSettings представляє налаштування сповіщень
type NotificationSettings struct {
	Email             bool `json:"email"`
	PushNotifications bool `json:"push_notifications"`
	SMS               bool `json:"sms"`
}

// UserPreferences представляє загальні налаштування користувача
type UserPreferences struct {
	DefaultView    string `json:"default_view"`
	FirstDayOfWeek int    `json:"first_day_of_week"`
}

// ThemeSettings представляє налаштування теми
type ThemeSettings struct {
	Mode         string `json:"mode"` // light, dark, system
	PrimaryColor string `json:"primary_color"`
	AccentColor  string `json:"accent_color"`
}

// PublicUser - публічна інформація про користувача (без конфіденційних даних)
type PublicUser struct {
	ID     uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
	Role   string    `json:"role"`
}

// UserPublic забезпечує сумісність зі старим кодом
type UserPublic = PublicUser

// ToPublicUser конвертує User в PublicUser
func (u *User) ToPublicUser() PublicUser {
	var avatar string
	if u.AvatarURL != nil {
		avatar = *u.AvatarURL
	}

	return PublicUser{
		ID:     u.ID,
		Email:  u.Email,
		Name:   u.Name,
		Avatar: avatar,
		Role:   u.Role,
	}
}

// DefaultPermissions повертає стандартні дозволи в залежності від ролі
func DefaultPermissions(role string) Permissions {
	switch role {
	case RoleOwner, RoleAdmin:
		return Permissions{
			ViewFinancials: true,
			EditProjects:   true,
			ViewProjects:   true,
			ManageTeam:     true,
		}
	case RoleEditor, RoleRetoucher:
		return Permissions{
			ViewFinancials: false,
			EditProjects:   true,
			ViewProjects:   true,
			ManageTeam:     false,
		}
	case RoleOperator, RoleAssistant:
		return Permissions{
			ViewFinancials: false,
			EditProjects:   false,
			ViewProjects:   true,
			ManageTeam:     false,
		}
	default:
		return Permissions{
			ViewFinancials: false,
			EditProjects:   false,
			ViewProjects:   true,
			ManageTeam:     false,
		}
	}
}

// BeforeCreate - GORM хук для генерації UUID перед створенням
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Value - реалізація інтерфейсу Valuer для Permissions
func (p Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan - реалізація інтерфейсу Scanner для Permissions
func (p *Permissions) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &p)
}
