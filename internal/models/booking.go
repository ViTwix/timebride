package models

import (
    "encoding/json"
    "time"
    "github.com/google/uuid"
)

type Booking struct {
    ID           uuid.UUID       `json:"id"`
    UserID       uuid.UUID       `json:"user_id"`
    Title        string         `json:"title"`
    Status       string         `json:"status"`
    ClientName   string         `json:"client_name"`
    StartTime    time.Time      `json:"start_time"`
    EndTime      time.Time      `json:"end_time"`
    IsAllDay     bool           `json:"is_all_day"`
    CustomFields json.RawMessage `json:"custom_fields"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
}

// Статуси бронювання
const (
    BookingStatusPending   = "pending"
    BookingStatusConfirmed = "confirmed"
    BookingStatusCanceled  = "canceled"
    BookingStatusCompleted = "completed"
)
