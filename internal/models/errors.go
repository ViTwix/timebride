package models

// ErrValidation представляє помилку валідації
type ErrValidation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ErrValidation) Error() string {
	return e.Message
}

// IsValidationError перевіряє чи є помилка помилкою валідації
func IsValidationError(err error) bool {
	_, ok := err.(ErrValidation)
	return ok
}

// NewValidationError створює нову помилку валідації
func NewValidationError(field, message string) error {
	return ErrValidation{
		Field:   field,
		Message: message,
	}
}
