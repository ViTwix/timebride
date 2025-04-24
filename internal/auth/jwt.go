package auth

import (
	"time"

	"github.com/google/uuid"

	"timebride/internal/models"
)

// JWTConfig містить налаштування для JWT токенів
type JWTConfig struct {
	Secret                  string        `yaml:"secret"`
	AccessExpirationMinutes time.Duration `yaml:"access_expiration_minutes"`
	RefreshExpirationDays   time.Duration `yaml:"refresh_expiration_days"`
	Issuer                  string        `yaml:"issuer"`
	Audience                string        `yaml:"audience"`
}

// GetAccessExpiration повертає тривалість дії access токена
func (c *JWTConfig) GetAccessExpiration() time.Duration {
	return c.AccessExpirationMinutes * time.Minute
}

// GetRefreshExpiration повертає тривалість дії refresh токена
func (c *JWTConfig) GetRefreshExpiration() time.Duration {
	return c.RefreshExpirationDays * 24 * time.Hour
}

// GetSecret повертає секретний ключ для підпису токенів
func (c *JWTConfig) GetSecret() []byte {
	return []byte(c.Secret)
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	Exp    int64     `json:"exp"`
	Iat    int64     `json:"iat"`
}

// NewJWTClaims creates a new JWT claims object
func NewJWTClaims(user *models.User, expirationTime time.Duration) *JWTClaims {
	now := time.Now()
	return &JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		Exp:    now.Add(expirationTime).Unix(),
		Iat:    now.Unix(),
	}
}
