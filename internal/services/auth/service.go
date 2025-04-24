package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"timebride/internal/config"
	"timebride/internal/models"
	"timebride/internal/repositories"
	"timebride/internal/types"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
)

type authService struct {
	config       *config.Config
	userRepo     repositories.UserRepository
	refreshToken string
}

// NewAuthService creates a new auth service instance
func NewAuthService(cfg *config.Config, userRepo repositories.UserRepository) IAuthService {
	return &authService{
		config:   cfg,
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, email, password, name string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*models.User, *types.AuthTokens, error) {
	filter := map[string]interface{}{"email": email}
	users, err := s.userRepo.List(ctx, filter)
	if err != nil || len(users) == 0 {
		return nil, nil, ErrInvalidCredentials
	}
	user := users[0]

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *authService) Verify(ctx context.Context, token string) (*models.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return s.GetJWTSecret(), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshTokenString string) (*models.User, *types.AuthTokens, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.GetJWTSecret(), nil
	})

	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, ErrUserNotFound
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *authService) GenerateOAuthURL(provider string) (string, string, error) {
	// Implement OAuth URL generation based on provider
	return "", "", nil
}

func (s *authService) HandleOAuthCallback(ctx context.Context, provider, code, state, savedState string) (*models.User, *types.AuthTokens, error) {
	// Implement OAuth callback handling
	return nil, nil, nil
}

func (s *authService) GetJWTSecret() []byte {
	return []byte(s.config.JWT.Secret)
}

func (s *authService) generateTokens(user *models.User) (*types.AuthTokens, error) {
	now := time.Now()
	accessExp := now.Add(time.Duration(s.config.JWT.AccessExpirationMinutes) * time.Minute)
	refreshExp := now.Add(time.Duration(s.config.JWT.RefreshExpirationDays) * 24 * time.Hour)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": accessExp.Unix(),
		"iat": now.Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": refreshExp.Unix(),
		"iat": now.Unix(),
	})

	accessTokenString, err := accessToken.SignedString(s.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(s.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	return &types.AuthTokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExp,
	}, nil
}
