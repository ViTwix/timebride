package errors

import (
	"fmt"
	"net/http"
)

// ErrorType представляє тип помилки
type ErrorType string

const (
	// ErrorTypeValidation - помилки валідації
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeNotFound - ресурс не знайдено
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeUnauthorized - не авторизований доступ
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	// ErrorTypeForbidden - доступ заборонено
	ErrorTypeForbidden ErrorType = "forbidden"
	// ErrorTypeInternal - внутрішня помилка сервера
	ErrorTypeInternal ErrorType = "internal"
	// ErrorTypeBadRequest - неправильний запит
	ErrorTypeBadRequest ErrorType = "bad_request"
)

// AppError представляє помилку в додатку
type AppError struct {
	Type     ErrorType `json:"type"`
	Message  string    `json:"message"`
	Details  any       `json:"details,omitempty"`
	Internal error     `json:"-"`
	HTTPCode int       `json:"-"`
}

// Error реалізує інтерфейс error
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// NewValidationError створює нову помилку валідації
func NewValidationError(message string, details any) *AppError {
	return &AppError{
		Type:     ErrorTypeValidation,
		Message:  message,
		Details:  details,
		HTTPCode: http.StatusBadRequest,
	}
}

// NewNotFoundError створює нову помилку "не знайдено"
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:     ErrorTypeNotFound,
		Message:  message,
		HTTPCode: http.StatusNotFound,
	}
}

// NewUnauthorizedError створює нову помилку авторизації
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:     ErrorTypeUnauthorized,
		Message:  message,
		HTTPCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError створює нову помилку доступу
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:     ErrorTypeForbidden,
		Message:  message,
		HTTPCode: http.StatusForbidden,
	}
}

// NewInternalError створює нову внутрішню помилку
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:     ErrorTypeInternal,
		Message:  message,
		Internal: err,
		HTTPCode: http.StatusInternalServerError,
	}
}

// NewBadRequestError створює нову помилку неправильного запиту
func NewBadRequestError(message string, details any) *AppError {
	return &AppError{
		Type:     ErrorTypeBadRequest,
		Message:  message,
		Details:  details,
		HTTPCode: http.StatusBadRequest,
	}
}

// IsNotFound перевіряє чи помилка є типу "не знайдено"
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsValidation перевіряє чи помилка є помилкою валідації
func IsValidation(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeValidation
	}
	return false
}

// IsUnauthorized перевіряє чи помилка є помилкою авторизації
func IsUnauthorized(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeUnauthorized
	}
	return false
}
