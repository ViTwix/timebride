package user

import (
	"context"

	"timebride/internal/models"

	"github.com/google/uuid"
)

// IUserService defines user management service interface
type IUserService interface {
	Register(ctx context.Context, email, password, name, role string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, user *models.User, currentPassword, newPassword string) error
	UpdateSettings(ctx context.Context, userID uuid.UUID, settings models.UserSettings) error
	GetSettings(ctx context.Context, userID uuid.UUID) (models.UserSettings, error)
	List(ctx context.Context) ([]*models.User, error)
}
