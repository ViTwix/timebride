package models

import (
	"time"

	"github.com/google/uuid"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	Exp    int64     `json:"exp"`
	Iat    int64     `json:"iat"`
}

// NewJWTClaims creates a new JWT claims object
func NewJWTClaims(user *User, expirationTime time.Duration) *JWTClaims {
	now := time.Now()
	return &JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		Exp:    now.Add(expirationTime).Unix(),
		Iat:    now.Unix(),
	}
}
