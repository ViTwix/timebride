package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// File модель для зберігання інформації про файли
type File struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null"`
	BookingID   *uuid.UUID     `gorm:"type:uuid"`
	FileName    string         `gorm:"type:text;not null"`
	FileType    string         `gorm:"type:text;not null"`
	FileSize    int64          `gorm:"not null"`
	Bucket      string         `gorm:"type:text;not null"`
	Key         string         `gorm:"type:text;not null"`
	CDNURL      string         `gorm:"type:text;not null"`
	Description string         `gorm:"type:text"`
	Tags        string         `gorm:"type:text"` // Зберігаємо як JSON рядок
	Metadata    string         `gorm:"type:text"` // Зберігаємо як JSON рядок
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName вказує ім'я таблиці для File
func (File) TableName() string {
	return "files"
}

// GetTags отримує теги як масив рядків
func (f *File) GetTags() []string {
	var tags []string
	json.Unmarshal([]byte(f.Tags), &tags)
	return tags
}

// SetTags встановлює теги з масиву рядків
func (f *File) SetTags(tags []string) {
	tagsJSON, _ := json.Marshal(tags)
	f.Tags = string(tagsJSON)
}

// GetMetadata отримує метадані як map
func (f *File) GetMetadata() map[string]any {
	var metadata map[string]any
	json.Unmarshal([]byte(f.Metadata), &metadata)
	return metadata
}

// SetMetadata встановлює метадані з map
func (f *File) SetMetadata(metadata map[string]any) {
	metadataJSON, _ := json.Marshal(metadata)
	f.Metadata = string(metadataJSON)
}

// BeforeCreate виконується перед створенням запису
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
