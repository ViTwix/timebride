package models

// LoginInput структура для входу
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterInput структура для реєстрації
type RegisterInput struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	CompanyName string `json:"company_name"`
	Phone       string `json:"phone"`
}
