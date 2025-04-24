package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"timebride/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository handles database operations for users
type UserRepository interface {
	Repository[models.User]

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetByDomain retrieves a user by their domain
	GetByDomain(ctx context.Context, domain string) (*models.User, error)

	// GetSubUsers retrieves all users under a specific admin
	GetSubUsers(ctx context.Context, adminID uuid.UUID) ([]*models.User, error)
}

type userRepository struct {
	baseRepository[models.User]
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		baseRepository: baseRepository[models.User]{db: db},
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByDomain(ctx context.Context, domain string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "domain = ?", domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetSubUsers(ctx context.Context, adminID uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).Where("parent_admin_id = ?", adminID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) List(ctx context.Context, filter map[string]interface{}) ([]*models.User, error) {
	var users []*models.User
	query := r.db.WithContext(ctx)

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
