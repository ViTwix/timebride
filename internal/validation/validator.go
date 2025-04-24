package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"timebride/internal/errors"
)

var (
	validate *validator.Validate
	// emailRegex визначає регулярний вираз для валідації email
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	// phoneRegex визначає регулярний вираз для валідації телефону
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func init() {
	validate = validator.New()

	// Реєструємо кастомні валідатори
	_ = validate.RegisterValidation("uuid", validateUUID)
	_ = validate.RegisterValidation("phone", validatePhone)
	_ = validate.RegisterValidation("email", validateEmail)
}

// ValidateStruct валідує структуру
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.NewValidationError("validation failed", validationErrors)
		}
		return errors.NewInternalError("validation error", err)
	}
	return nil
}

// validateUUID перевіряє чи string є валідним UUID
func validateUUID(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	_, err := uuid.Parse(field)
	return err == nil
}

// validatePhone перевіряє чи string є валідним номером телефону
func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

// validateEmail перевіряє чи string є валідним email
func validateEmail(fl validator.FieldLevel) bool {
	return emailRegex.MatchString(fl.Field().String())
}

// ValidationRules містить правила валідації для різних полів
var ValidationRules = struct {
	Password string
	Email    string
	Phone    string
	Name     string
}{
	Password: "required,min=8,max=72",
	Email:    "required,email",
	Phone:    "required,phone",
	Name:     "required,min=2,max=50",
}
