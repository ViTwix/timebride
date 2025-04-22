package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Template struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name           string         `json:"name" gorm:"not null"`
	Description    string         `json:"description"`
	EventType      string         `json:"event_type" gorm:"not null"`
	FieldsTemplate datatypes.JSON `json:"fields_template" gorm:"type:jsonb;default:'{}'"`
	TeamTemplate   datatypes.JSON `json:"team_template" gorm:"type:jsonb;default:'[]'"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate - GORM хук для генерації UUID перед створенням
func (t *Template) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
