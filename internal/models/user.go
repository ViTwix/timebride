package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type User struct {
	ID                 uuid.UUID       `json:"id" gorm:"type:uuid;primary_key"`
	Email              string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash       string         `json:"-" gorm:"not null"`
	FullName           string         `json:"full_name" gorm:"not null"`
	CompanyName        string         `json:"company_name"`
	Phone              string         `json:"phone"`
	Role               string         `json:"role" gorm:"not null;default:'user'"`
	Language           string         `json:"language" gorm:"default:'en'"`
	Timezone           string         `json:"timezone" gorm:"default:'UTC'"`
	Theme              string         `json:"theme" gorm:"default:'light'"`
	SubscriptionPlan   string         `json:"subscription_plan" gorm:"default:'free'"`
	GoogleCalendarToken string        `json:"-"`
	TelegramChatID     string         `json:"-"`
	Settings           datatypes.JSON `json:"settings" gorm:"type:jsonb;default:'{}'"`
	TeamSettings       datatypes.JSON `json:"team_settings" gorm:"type:jsonb;default:'{}'"`
	CustomFields       datatypes.JSON `json:"custom_fields" gorm:"type:jsonb;default:'{}'"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// Ролі користувачів
const (
	RoleAdmin      = "admin"
	RolePhotographer = "photographer"
	RoleVideographer = "videographer"
	RoleEditor      = "editor"
	RoleManager     = "manager"
)

func (u *User) BeforeCreate() error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
