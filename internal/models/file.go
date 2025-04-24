package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrEmptyFileName   = errors.New("file name cannot be empty")
	ErrInvalidFileSize = errors.New("file size must be greater than 0")
	ErrEmptyMimeType   = errors.New("mime type cannot be empty")
)

// FileType represents the type of file
type FileType string

const (
	FileTypeAvatar   FileType = "avatar"
	FileTypeDocument FileType = "document"
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
)

// File represents a file in the system
type File struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	BookingID   *uuid.UUID `gorm:"type:uuid" json:"booking_id,omitempty"`
	Name        string     `gorm:"not null" json:"name"`
	Path        string     `gorm:"not null" json:"path"`
	Size        int64      `gorm:"not null" json:"size"`
	ContentType string     `gorm:"not null" json:"content_type"`
	Type        FileType   `gorm:"not null" json:"type"`
	MimeType    string     `gorm:"not null" json:"mime_type"`
	URL         string     `gorm:"not null" json:"url"`
	PublicURL   string     `gorm:"not null" json:"public_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Зв'язки
	User    *User    `gorm:"foreignKey:UserID" json:"-"`
	Booking *Booking `json:"-" gorm:"foreignKey:BookingID"`
}

// FilePublic represents a public view of a file
type FilePublic struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	PublicURL string    `json:"public_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ToPublic converts File to FilePublic
func (f *File) ToPublic() FilePublic {
	return FilePublic{
		ID:        f.ID,
		Name:      f.Name,
		Size:      f.Size,
		MimeType:  f.MimeType,
		PublicURL: f.PublicURL,
		CreatedAt: f.CreatedAt,
	}
}

// Validate checks if the file is valid
func (f *File) Validate() error {
	if f.Name == "" {
		return ErrEmptyFileName
	}
	if f.Size <= 0 {
		return ErrInvalidFileSize
	}
	if f.MimeType == "" {
		return ErrEmptyMimeType
	}
	return nil
}

// IsImage checks if the file is an image
func (f *File) IsImage() bool {
	switch f.MimeType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

// IsVideo checks if the file is a video
func (f *File) IsVideo() bool {
	switch f.MimeType {
	case "video/mp4", "video/quicktime", "video/x-msvideo":
		return true
	default:
		return false
	}
}

// GetStorageKey returns the key for storing the file
func (f *File) GetStorageKey() string {
	if f.BookingID != nil {
		return f.UserID.String() + "/bookings/" + f.BookingID.String() + "/" + f.Name
	}
	return f.UserID.String() + "/files/" + f.Name
}

// GetHumanSize returns the file size in a human-readable format
func (f *File) GetHumanSize() string {
	const unit = 1024
	if f.Size < unit {
		return fmt.Sprintf("%d B", f.Size)
	}
	div, exp := int64(unit), 0
	for n := f.Size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(f.Size)/float64(div), "KMGTPE"[exp])
}

// BeforeCreate generates a new UUID for the file if not set
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (f *File) MarshalJSON() ([]byte, error) {
	type Alias File
	return json.Marshal(&struct {
		*Alias
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		Alias:     (*Alias)(f),
		CreatedAt: f.CreatedAt.UTC(),
		UpdatedAt: f.UpdatedAt.UTC(),
	})
}

// FileCreate represents data for creating a new file
type FileCreate struct {
	UserID    uuid.UUID  `json:"user_id" validate:"required"`
	BookingID *uuid.UUID `json:"booking_id,omitempty"`
	Name      string     `json:"name" validate:"required"`
	Size      int64      `json:"size" validate:"required"`
	MimeType  string     `json:"mime_type" validate:"required"`
	URL       string     `json:"url" validate:"required"`
}

// FileUpdate represents data for updating a file
type FileUpdate struct {
	Name      *string    `json:"name,omitempty"`
	BookingID *uuid.UUID `json:"booking_id,omitempty"`
	Size      *int64     `json:"size,omitempty"`
	MimeType  *string    `json:"mime_type,omitempty"`
	URL       *string    `json:"url,omitempty"`
}
