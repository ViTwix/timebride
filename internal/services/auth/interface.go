package auth

import (
	"context"

	"timebride/internal/models"
	"timebride/internal/types"
)

// IAuthService defines authentication service interface
type IAuthService interface {
	Register(ctx context.Context, email, password, name string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, *types.AuthTokens, error)
	Verify(ctx context.Context, token string) (*models.User, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*models.User, *types.AuthTokens, error)
	GenerateOAuthURL(provider string) (string, string, error)
	HandleOAuthCallback(ctx context.Context, provider, code, state, savedState string) (*models.User, *types.AuthTokens, error)
	GetJWTSecret() []byte
}
