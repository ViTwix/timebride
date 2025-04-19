package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Template struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primary_key"`
	UserID       uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	Name         string         `json:"name" gorm:"not null"`
	Description  string         `json:"description"`
	EventType    string         `json:"event_type" gorm:"not null"`
	FieldsTemplate datatypes.JSON `json:"fields_template" gorm:"type:jsonb;default:'{}'"`
	TeamTemplate  datatypes.JSON `json:"team_template" gorm:"type:jsonb;default:'[]'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (t *Template) BeforeCreate() error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
} 