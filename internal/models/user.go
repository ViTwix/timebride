package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email        string         `json:"email" gorm:"unique;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	FullName     string         `json:"full_name" gorm:"not null"`
	CompanyName  string         `json:"company_name"`
	Phone        string         `json:"phone"`
	Role         string         `json:"role" gorm:"not null;default:'user'"`
	Settings     datatypes.JSON `json:"settings" gorm:"type:jsonb"`
	CreatedAt    time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt    *time.Time     `json:"-" gorm:"index"`
}

// UserSettings представляє налаштування користувача
type UserSettings struct {
	Theme            string            `json:"theme"`
	Language         string            `json:"language"`
	Notifications    bool              `json:"notifications"`
	DefaultCurrency  string            `json:"default_currency"`
	CustomFields     map[string]string `json:"custom_fields"`
	CalendarSettings CalendarSettings  `json:"calendar_settings"`
}

// CalendarSettings представляє налаштування календаря
type CalendarSettings struct {
	DefaultView string `json:"default_view"`
	StartOfWeek int    `json:"start_of_week"`
	WorkingDays []int  `json:"working_days"`
}

// PublicUser представляє публічну інформацію про користувача
type PublicUser struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	CompanyName string    `json:"company_name,omitempty"`
	Role        string    `json:"role"`
}

// ToPublic конвертує User в PublicUser
func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:          u.ID,
		Email:       u.Email,
		FullName:    u.FullName,
		CompanyName: u.CompanyName,
		Role:        u.Role,
	}
}

// GetSettings повертає налаштування користувача
func (u *User) GetSettings() (*UserSettings, error) {
	if u.Settings == nil {
		return &UserSettings{
			Theme:           "light",
			Language:        "uk",
			DefaultCurrency: "UAH",
			CustomFields:    make(map[string]string),
			CalendarSettings: CalendarSettings{
				DefaultView: "month",
				StartOfWeek: 1,
				WorkingDays: []int{1, 2, 3, 4, 5},
			},
		}, nil
	}

	var settings UserSettings
	if err := json.Unmarshal(u.Settings, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

// SetSettings встановлює налаштування користувача
func (u *User) SetSettings(settings *UserSettings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	u.Settings = datatypes.JSON(data)
	return nil
}

// BeforeCreate встановлює значення за замовчуванням перед створенням
func (u *User) BeforeCreate() error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	if u.Role == "" {
		u.Role = "user"
	}

	// Встановлюємо налаштування за замовчуванням
	if u.Settings == nil {
		settings := &UserSettings{
			Theme:           "light",
			Language:        "uk",
			DefaultCurrency: "UAH",
			CustomFields:    make(map[string]string),
			CalendarSettings: CalendarSettings{
				DefaultView: "month",
				StartOfWeek: 1,
				WorkingDays: []int{1, 2, 3, 4, 5},
			},
		}
		if err := u.SetSettings(settings); err != nil {
			return err
		}
	}
	return nil
}

// BeforeUpdate оновлює час модифікації перед оновленням
func (u *User) BeforeUpdate() error {
	u.UpdatedAt = time.Now()
	return nil
}
