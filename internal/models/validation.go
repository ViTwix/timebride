package models

import (
	"fmt"
)

// Validate перевіряє коректність даних користувача
func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if u.PasswordHash == "" {
		return fmt.Errorf("password is required")
	}
	if u.Role == "" {
		return fmt.Errorf("role is required")
	}
	return nil
}

// Validate перевіряє коректність даних бронювання
func (b *Booking) Validate() error {
	if b.Title == "" {
		return fmt.Errorf("title is required")
	}
	if b.Status == "" {
		return fmt.Errorf("status is required")
	}
	if b.EventDate.IsZero() {
		return fmt.Errorf("event date is required")
	}
	if b.StartTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	if b.EndTime.IsZero() {
		return fmt.Errorf("end time is required")
	}
	if b.EndTime.Before(b.StartTime) {
		return fmt.Errorf("end time cannot be before start time")
	}
	// Перевірка в минулому відключена, щоб можна було створювати бронювання для минулих подій
	// if b.EventDate.Before(time.Now()) {
	//     return fmt.Errorf("cannot create booking in the past")
	// }
	return nil
}
