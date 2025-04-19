package models

import (
	"time"
)

// Activity представляє активність користувача в системі
type Activity struct {
	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID      int64     `bun:"user_id" json:"userId"`
	EntityType  string    `bun:"entity_type" json:"entityType"`  // Тип сутності (client, booking, etc)
	EntityID    int64     `bun:"entity_id" json:"entityId"`      // ID сутності
	Action      string    `bun:"action" json:"action"`           // Тип дії (create, update, delete)
	Description string    `bun:"description" json:"description"` // Опис дії
	CreatedAt   time.Time `bun:"created_at,default:current_timestamp" json:"createdAt"`
}

// TableName повертає назву таблиці для моделі Activity
func (a *Activity) TableName() string {
	return "activities"
}

// NewCreateActivity створює нову активність типу "створення"
func NewCreateActivity(userID int64, entityType string, entityID int64, description string) Activity {
	return Activity{
		UserID:      userID,
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      "create",
		Description: description,
		CreatedAt:   time.Now(),
	}
}

// NewUpdateActivity створює нову активність типу "оновлення"
func NewUpdateActivity(userID int64, entityType string, entityID int64, description string) Activity {
	return Activity{
		UserID:      userID,
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      "update",
		Description: description,
		CreatedAt:   time.Now(),
	}
}

// NewDeleteActivity створює нову активність типу "видалення"
func NewDeleteActivity(userID int64, entityType string, entityID int64, description string) Activity {
	return Activity{
		UserID:      userID,
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      "delete",
		Description: description,
		CreatedAt:   time.Now(),
	}
}
