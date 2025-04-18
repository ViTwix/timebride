package models

import (
    "time"
    "encoding/json"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID       `json:"id"`
    Email        string         `json:"email"`
    PasswordHash string         `json:"-"`
    FullName     string         `json:"full_name"`
    Role         string         `json:"role"`
    Settings     json.RawMessage `json:"settings"`
    CustomFields json.RawMessage `json:"custom_fields"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
}

// Ролі користувачів
const (
    RoleAdmin      = "admin"
    RolePhotographer = "photographer"
    RoleVideographer = "videographer"
    RoleEditor      = "editor"
    RoleManager     = "manager"
)
