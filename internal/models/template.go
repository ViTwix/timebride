package models

import (
	"time"

	"github.com/google/uuid"
)

// Template представляє шаблон документа або повідомлення
type Template struct {
	ID          uuid.UUID         `json:"id" gorm:"primarykey;type:uuid"`
	UserID      uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Subject     string            `json:"subject"`
	Content     string            `json:"content"`
	Variables   map[string]string `json:"variables" gorm:"type:jsonb;serializer:json"`
	IsDefault   bool              `json:"is_default"`
	IsActive    bool              `json:"is_active" gorm:"default:true"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// BeforeCreate генерує UUID для нового шаблону
func (t *Template) BeforeCreate() error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// BeforeUpdate оновлює час модифікації перед оновленням
func (t *Template) BeforeUpdate() error {
	t.UpdatedAt = time.Now()
	return nil
}
