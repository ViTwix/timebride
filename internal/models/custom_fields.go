package models

import (
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

// CustomField представляє кастомне поле
type CustomField struct {
	ID         string      `json:"id"`
	Label      string      `json:"label"`
	Type       string      `json:"type"`
	Value      interface{} `json:"value"`
	Required   bool        `json:"required,omitempty"`
	Options    []string    `json:"options,omitempty"`    // Для select полів
	Validation string      `json:"validation,omitempty"` // Регулярний вираз для валідації
}

// SetCustomField встановлює значення кастомного поля
func (u *User) SetCustomField(field CustomField) error {
	var fields map[string]CustomField
	if err := json.Unmarshal(u.CustomFields, &fields); err != nil {
		fields = make(map[string]CustomField)
	}
	fields[field.ID] = field

	data, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("error marshaling custom fields: %w", err)
	}

	u.CustomFields = datatypes.JSON(data)
	return nil
}

// GetCustomField повертає значення кастомного поля
func (u *User) GetCustomField(fieldID string) (*CustomField, error) {
	var fields map[string]CustomField
	if err := json.Unmarshal(u.CustomFields, &fields); err != nil {
		return nil, fmt.Errorf("error unmarshaling custom fields: %w", err)
	}

	field, ok := fields[fieldID]
	if !ok {
		return nil, fmt.Errorf("field %s not found", fieldID)
	}

	return &field, nil
}
