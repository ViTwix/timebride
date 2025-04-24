package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TeamMember представляє члена команди підрядника
type TeamMember struct {
	BaseModel           `json:",inline"`
	ModelWithTimestamps `json:",inline"`
	UserID              uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name                string         `json:"name" gorm:"not null"`
	Email               string         `json:"email"`
	Role                string         `json:"role" gorm:"not null"`
	Permissions         datatypes.JSON `json:"permissions" gorm:"type:jsonb;default:'{}'"`
	Settings            datatypes.JSON `json:"settings" gorm:"type:jsonb;default:'{}'"`

	// Зв'язки
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// TeamMemberPublic представляє публічну інформацію про члена команди
type TeamMemberPublic struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

// ToPublic конвертує TeamMember в TeamMemberPublic
func (tm *TeamMember) ToPublic() TeamMemberPublic {
	return TeamMemberPublic{
		ID:    tm.ID,
		Name:  tm.Name,
		Email: tm.Email,
		Role:  tm.Role,
	}
}

// BeforeCreate - GORM хук для генерації UUID перед створенням
func (tm *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if tm.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return nil
}
