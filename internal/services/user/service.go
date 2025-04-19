package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"timebride/internal/models"
	"timebride/internal/repositories"
)

type Service struct {
	repo repositories.UserRepository
}

func NewService(repo repositories.UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, email, password, fullName, companyName, role string) (*models.User, error) {
	// Хешуємо пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Створюємо користувача
	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		CompanyName:  companyName,
		Role:         role,
	}

	// Зберігаємо користувача
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*models.User, error) {
	// Отримуємо користувача за email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Перевіряємо пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByEmail повертає користувача за його email
func (s *Service) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) Update(ctx context.Context, user *models.User) error {
	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) UpdatePassword(ctx context.Context, user *models.User, currentPassword, newPassword string) error {
	// Перевіряємо поточний пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("invalid current password")
	}

	// Хешуємо новий пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Оновлюємо пароль
	user.PasswordHash = string(hashedPassword)
	return s.repo.Update(ctx, user)
}
