package models

import (
	"time"
)

// Client представляє клієнта в системі
type Client struct {
	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID    int64     `bun:"user_id" json:"userId"`
	FullName  string    `bun:"full_name" json:"fullName"`
	Email     string    `bun:"email" json:"email"`
	Phone     string    `bun:"phone" json:"phone"`
	Company   string    `bun:"company" json:"company"`
	Category  string    `bun:"category" json:"category"`
	Source    string    `bun:"source" json:"source"`
	Address   string    `bun:"address" json:"address"`
	Notes     string    `bun:"notes" json:"notes"`
	IsActive  bool      `bun:"is_active" json:"isActive"`
	Avatar    string    `bun:"avatar" json:"avatar"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp" json:"updatedAt"`
}

// TableName повертає назву таблиці для моделі Client
func (c *Client) TableName() string {
	return "clients"
}

// DefaultSources повертає список джерел клієнтів за замовчуванням
func DefaultSources() []string {
	return []string{
		"Сайт",
		"Рекомендація",
		"Соціальні мережі",
		"Виставка",
		"Партнери",
		"Інше",
	}
}

// DefaultCategories повертає список категорій клієнтів за замовчуванням
func DefaultCategories() []string {
	return []string{
		"Наречена",
		"Наречений",
		"Фотограф",
		"Відеограф",
		"Агенція",
		"Локація",
	}
}
