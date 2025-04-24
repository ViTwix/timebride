package models

import (
	"encoding/json"
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
	if u.FullName == "" {
		return fmt.Errorf("full name is required")
	}

	settings, err := u.GetSettings()
	if err != nil {
		return fmt.Errorf("invalid settings: %v", err)
	}

	if settings.Language != "uk" && settings.Language != "en" {
		return fmt.Errorf("language must be either uk or en")
	}
	if settings.Theme != "light" && settings.Theme != "dark" {
		return fmt.Errorf("theme must be either light or dark")
	}
	if settings.DefaultCurrency != "UAH" && settings.DefaultCurrency != "USD" && settings.DefaultCurrency != "EUR" {
		return fmt.Errorf("currency must be UAH, USD or EUR")
	}
	return nil
}

// Validate перевіряє коректність даних члена команди
func (tm *TeamMember) Validate() error {
	if tm.Name == "" {
		return fmt.Errorf("name is required")
	}
	if tm.Role == "" {
		return fmt.Errorf("role is required")
	}

	var permissions struct {
		AccessLevel string `json:"access_level"`
	}
	if err := json.Unmarshal(tm.Permissions, &permissions); err != nil {
		return fmt.Errorf("invalid permissions format")
	}

	if permissions.AccessLevel != "full" && permissions.AccessLevel != "assigned_only" && permissions.AccessLevel != "assigned_with_finance" {
		return fmt.Errorf("access level must be full, assigned_only or assigned_with_finance")
	}
	return nil
}
