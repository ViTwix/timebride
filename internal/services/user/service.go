package user

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"timebride/internal/models"
	"timebride/internal/repositories"
)

type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo repositories.UserRepository) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Register реєструє нового користувача
func (s *userService) Register(ctx context.Context, email, password, name, role string) (*models.User, error) {
	// Хешуємо пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Створюємо користувача
	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     name,
		Role:         role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login виконує вхід користувача
func (s *userService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID отримує користувача за ID
func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetByEmail отримує користувача за email
func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	users, err := s.userRepo.List(ctx, map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return users[0], nil
}

// Update оновлює існуючого користувача
func (s *userService) Update(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

// Delete видаляє користувача
func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

// UpdatePassword оновлює пароль користувача
func (s *userService) UpdatePassword(ctx context.Context, user *models.User, currentPassword, newPassword string) error {
	// Перевіряємо поточний пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return err
	}

	// Хешуємо новий пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// UpdateSettings оновлює налаштування користувача
func (s *userService) UpdateSettings(ctx context.Context, userID uuid.UUID, settings models.UserSettings) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	user.Settings = datatypes.JSON(settingsJSON)
	return s.userRepo.Update(ctx, user)
}

// GetSettings отримує налаштування користувача
func (s *userService) GetSettings(ctx context.Context, userID uuid.UUID) (models.UserSettings, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return models.UserSettings{}, err
	}

	var settings models.UserSettings
	if err := json.Unmarshal(user.Settings, &settings); err != nil {
		return models.UserSettings{}, err
	}

	return settings, nil
}

// List отримує список користувачів
func (s *userService) List(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.List(ctx, nil)
}
