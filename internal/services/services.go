package services

import (
	"context"
	"io"
	"mime/multipart"
	"time"

	"timebride/internal/models"

	"github.com/google/uuid"
)

// UserService визначає інтерфейс для роботи з користувачами
type UserService interface {
	Register(ctx context.Context, email, password, name, role string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, user *models.User, currentPassword, newPassword string) error
}

// IAuthService визначає інтерфейс для аутентифікації
type IAuthService interface {
	Register(ctx context.Context, email, password, name string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, *AuthTokens, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*models.User, *AuthTokens, error)
	GenerateOAuthURL(provider string) (string, string, error)
	HandleOAuthCallback(ctx context.Context, provider, code, state, savedState string) (*models.User, *AuthTokens, error)
}

// StorageService визначає інтерфейс для роботи з файлами
type StorageService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID) (*models.File, error)
	DownloadFile(ctx context.Context, fileID uuid.UUID) (io.ReadCloser, string, error)
	DeleteFile(ctx context.Context, fileID uuid.UUID) error
	ListFiles(ctx context.Context, userID uuid.UUID) ([]*models.File, error)
	GetFileByID(ctx context.Context, id uuid.UUID) (*models.File, error)
	GetTotalFiles(ctx context.Context) (int64, error)
	GetTotalSize(ctx context.Context) (int64, error)
}

// TemplateService визначає інтерфейс для роботи з шаблонами
type TemplateService interface {
	Create(ctx context.Context, template *models.Template) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Template, error)
	Update(ctx context.Context, template *models.Template) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalTemplates(ctx context.Context) (int64, error)
	GetActiveTemplates(ctx context.Context) (int64, error)
}

// BookingService визначає інтерфейс для роботи з бронюваннями
type BookingService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	Create(ctx context.Context, booking *models.Booking) error
	Update(ctx context.Context, booking *models.Booking) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]*models.Booking, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*models.Booking, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string) error
	CountTotal(ctx context.Context, userID uuid.UUID) (int64, error)
	CountActive(ctx context.Context, userID uuid.UUID) (int64, error)
	CountUpcoming(ctx context.Context, userID uuid.UUID) (int64, error)
	CountThisMonth(ctx context.Context, userID uuid.UUID) (int64, error)
	GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Booking, error)
	AddTeamMember(ctx context.Context, bookingID uuid.UUID, member map[string]interface{}) error
	RemoveTeamMember(ctx context.Context, bookingID uuid.UUID, memberID string) error
}
