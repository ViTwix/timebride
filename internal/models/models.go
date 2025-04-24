package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSON is a custom type for handling JSONB fields
type JSON map[string]interface{}

// BaseModel містить спільні поля для всіх моделей
type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
}

// ModelWithTimestamps містить поля для відстеження часу
type ModelWithTimestamps struct {
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ModelWithSoftDelete contains common soft delete fields
type ModelWithSoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate генерує UUID перед створенням запису
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
