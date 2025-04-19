package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"timebride/internal/config"
	"timebride/internal/models"
	"timebride/internal/repositories"
	"timebride/internal/services/user"
)

type Service struct {
	userRepo  repositories.UserRepository
	jwtConfig config.JWTConfig
}

func NewService(userRepo repositories.UserRepository, jwtConfig config.JWTConfig) *Service {
	return &Service{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

type Handler struct {
	userService *user.Service
}

func NewHandler(userService *user.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

// Register registers a new user
func (s *Service) Register(ctx context.Context, email, password, fullName, companyName string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		CompanyName:  companyName,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(s.jwtConfig.TokenDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *Service) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetByID(context.Background(), userID)
}
