package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ============================================================================
// Common Types and Interfaces
// ============================================================================

// Error представляє помилку в системі
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// ValidationError представляє помилку валідації
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// BaseModel базова модель з ID
type BaseModel struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primarykey"`
}

// ModelWithTimestamps модель з часовими мітками
type ModelWithTimestamps struct {
	BaseModel
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt datatypes.Time `json:"-" gorm:"index"`
}

// ============================================================================
// Authentication Types
// ============================================================================

// AuthTokens представляє токени аутентифікації
type AuthTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// ============================================================================
// Team Types
// ============================================================================

// TeamRole тип ролі в команді
type TeamRole string

const (
	TeamRoleOwner  TeamRole = "owner"
	TeamRoleAdmin  TeamRole = "admin"
	TeamRoleMember TeamRole = "member"
	TeamRoleViewer TeamRole = "viewer"
)

// ============================================================================
// Booking Types
// ============================================================================

// BookingStatus статус бронювання
type BookingStatus string

const (
	BookingStatusDraft     BookingStatus = "draft"
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusBooked    BookingStatus = "booked"
	BookingStatusEditing   BookingStatus = "editing"
	BookingStatusReady     BookingStatus = "ready"
	BookingStatusArchived  BookingStatus = "archived"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusDone      BookingStatus = "done"
)

// ============================================================================
// List Options
// ============================================================================

// ListOptions базові опції для списків
type ListOptions struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
	SortBy   string `json:"sort_by"`
	SortDesc bool   `json:"sort_desc"`
}

// BookingListOptions опції для списку бронювань
type BookingListOptions struct {
	ListOptions
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// Booking представляє бронювання
type Booking struct {
	ModelWithTimestamps
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	ClientID        uuid.UUID      `json:"client_id" gorm:"type:uuid;not null"`
	Title           string         `json:"title"`
	EventType       string         `json:"event_type"`
	EventDate       time.Time      `json:"event_date"`
	StartTime       time.Time      `json:"start_time"`
	EndTime         time.Time      `json:"end_time"`
	Status          BookingStatus  `json:"status"`
	PaymentStatus   PaymentStatus  `json:"payment_status"`
	Location        string         `json:"location"`
	Description     string         `json:"description"`
	PackageName     string         `json:"package_name"`
	DeadlineDays    int            `json:"deadline_days"`
	Amount          float64        `json:"amount"`
	Currency        string         `json:"currency"`
	PriceTotal      float64        `json:"price_total"`
	PricePrepayment float64        `json:"price_prepayment"`
	PriceExtra      float64        `json:"price_extra"`
	TeamPayments    datatypes.JSON `json:"team_payments"`
	CustomFields    datatypes.JSON `json:"custom_fields"`
	ClientName      string         `json:"client_name"`
	ClientPhone     string         `json:"client_phone"`
	Type            string         `json:"type"`
	Category        string         `json:"category"`
	Metadata        datatypes.JSON `json:"metadata"`
}

// BookingCreate структура для створення бронювання
type BookingCreate struct {
	UserID       string         `json:"user_id" validate:"required,uuid"`
	ClientID     string         `json:"client_id" validate:"required,uuid"`
	StartTime    time.Time      `json:"start_time" validate:"required"`
	EndTime      time.Time      `json:"end_time" validate:"required,gtfield=StartTime"`
	Amount       float64        `json:"amount" validate:"required,gte=0"`
	Prepayment   float64        `json:"prepayment" validate:"gte=0"`
	Currency     string         `json:"currency" validate:"required"`
	Description  string         `json:"description"`
	Location     string         `json:"location"`
	TeamMembers  datatypes.JSON `json:"team_members"`
	CustomFields datatypes.JSON `json:"custom_fields"`
}

// BookingUpdate структура для оновлення бронювання
type BookingUpdate struct {
	StartTime    *time.Time      `json:"start_time,omitempty"`
	EndTime      *time.Time      `json:"end_time,omitempty"`
	Status       *BookingStatus  `json:"status,omitempty"`
	Amount       *float64        `json:"amount,omitempty"`
	Prepayment   *float64        `json:"prepayment,omitempty"`
	Currency     *string         `json:"currency,omitempty"`
	Description  *string         `json:"description,omitempty"`
	Location     *string         `json:"location,omitempty"`
	TeamMembers  *datatypes.JSON `json:"team_members,omitempty"`
	CustomFields *datatypes.JSON `json:"custom_fields,omitempty"`
}

// BookingWithClient розширена структура бронювання з даними клієнта
type BookingWithClient struct {
	Booking
	Client *Client `json:"client"`
}

// ============================================================================
// Payment Types
// ============================================================================

// PaymentStatus статус оплати
type PaymentStatus string

const (
	PaymentStatusNone       PaymentStatus = "none"
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusPaid       PaymentStatus = "paid"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCanceled   PaymentStatus = "canceled"
	PaymentStatusIncomplete PaymentStatus = "incomplete"
)

// ============================================================================
// Client Types
// ============================================================================

// ClientListOptions опції для списку клієнтів
type ClientListOptions struct {
	ListOptions
	Category string `json:"category"`
	Source   string `json:"source"`
}

// ClientListResult результат списку клієнтів
type ClientListResult struct {
	Total   int64    `json:"total"`
	Clients []Client `json:"clients"`
}

// Client представляє клієнта
type Client struct {
	ModelWithTimestamps
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	FullName     string         `json:"full_name"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	Notes        string         `json:"notes"`
	AvatarURL    string         `json:"avatar_url"`
	CustomFields datatypes.JSON `json:"custom_fields"`
	Status       string         `json:"status"`
	Source       string         `json:"source"`
	Tags         []string       `json:"tags" gorm:"-"`
	Categories   []string       `json:"categories" gorm:"-"`
	Metadata     datatypes.JSON `json:"metadata"`
}

// ============================================================================
// File Types
// ============================================================================

// File представляє файл
type File struct {
	ModelWithTimestamps
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	BookingID   uuid.UUID `json:"booking_id,omitempty" gorm:"type:uuid"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Type        FileType  `json:"type"`
	MimeType    string    `json:"mime_type"`
	URL         string    `json:"url"`
	PublicURL   string    `json:"public_url,omitempty"`
	Public      bool      `json:"public"`
}

// FileType визначає типи файлів
type FileType string

const (
	FileTypeAvatar   FileType = "avatar"
	FileTypeDocument FileType = "document"
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
)

// ============================================================================
// User Types
// ============================================================================

// User представляє користувача системи
type User struct {
	ModelWithTimestamps
	Email           string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash    string         `json:"-" gorm:"not null"`
	FullName        string         `json:"full_name" gorm:"not null"`
	CompanyName     string         `json:"company_name"`
	Phone           string         `json:"phone"`
	Settings        datatypes.JSON `json:"settings"`
	Language        string         `json:"language" gorm:"default:'uk'"`
	Theme           string         `json:"theme" gorm:"default:'light'"`
	DefaultCurrency string         `json:"default_currency" gorm:"default:'UAH'"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	PhoneNumber     string         `json:"phone_number"`
	Avatar          string         `json:"avatar"`
	Role            string         `json:"role" gorm:"not null;default:'user'"`
	Status          string         `json:"status"`
}

// UserRole визначає ролі користувачів
type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleManager  UserRole = "manager"
	UserRoleOperator UserRole = "operator"
	UserRoleClient   UserRole = "client"
)

// UserSettings налаштування користувача
type UserSettings struct {
	Language        string `json:"language"`
	Theme           string `json:"theme"`
	DefaultCurrency string `json:"default_currency"`
}

// ============================================================================
// Template Types
// ============================================================================

// Template представляє шаблон
type Template struct {
	ModelWithTimestamps
	UserID      uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Content     datatypes.JSON    `json:"content"`
	IsActive    bool              `json:"is_active" gorm:"default:true"`
	Type        string            `json:"type"`
	Subject     string            `json:"subject"`
	Variables   map[string]string `json:"variables" gorm:"type:jsonb;serializer:json"`
	IsDefault   bool              `json:"is_default"`
	Category    string            `json:"category"`
	Tags        []string          `json:"tags" gorm:"-"`
	Metadata    datatypes.JSON    `json:"metadata"`
}

// ============================================================================
// Price Types
// ============================================================================

// Price представляє ціну
type Price struct {
	ModelWithTimestamps
	UserID       string         `json:"user_id" gorm:"type:uuid;not null"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Amount       float64        `json:"amount"`
	Currency     string         `json:"currency"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CustomFields datatypes.JSON `json:"custom_fields"`
}

// TeamMember представляє члена команди
type TeamMember struct {
	ModelWithTimestamps
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	TeamID         uuid.UUID `json:"team_id" gorm:"type:uuid;not null"`
	Role           TeamRole  `json:"role"`
	JoinedAt       time.Time `json:"joined_at"`
	InvitedBy      uuid.UUID `json:"invited_by" gorm:"type:uuid"`
	TeamMemberName string    `json:"team_member_name"`
	AccessLevel    string    `json:"access_level"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
}

// PriceTemplate представляє шаблон ціни
type PriceTemplate struct {
	ModelWithTimestamps
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	EventType       string         `json:"event_type"`
	Currency        string         `json:"currency"`
	PriceTotal      float64        `json:"price_total"`
	PricePrepayment float64        `json:"price_prepayment"`
	DeadlineDays    int            `json:"deadline_days"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CustomFields    datatypes.JSON `json:"custom_fields"`
	Type            string         `json:"type"`
	Category        string         `json:"category"`
	Duration        int            `json:"duration"`
	Interval        string         `json:"interval"`
}
