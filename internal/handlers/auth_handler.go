package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"timebride/internal/middleware"
	"timebride/internal/models"
	"timebride/internal/services/auth"
	"timebride/internal/services/user"

	"github.com/google/uuid"
)

// AuthHandler обробляє запити аутентифікації
type AuthHandler struct {
	userService *user.Service
	authConfig  middleware.AuthConfig
	authService *auth.Service
}

// NewAuthHandler створює новий екземпляр AuthHandler
func NewAuthHandler(userService *user.Service, authConfig middleware.AuthConfig, authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authConfig:  authConfig,
		authService: authService,
	}
}

// RegisterRequest містить дані для реєстрації
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	CompanyName string `json:"company_name"`
	Role        string `json:"role"`
}

// LoginRequest містить дані для входу
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse містить дані відповіді аутентифікації
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// RegisterInput структура для реєстрації
type RegisterInput struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	CompanyName string `json:"company_name"`
	Phone       string `json:"phone"`
}

// LoginInput структура для входу
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register обробляє реєстрацію нового користувача
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Створюємо користувача
	user, err := h.userService.Register(r.Context(), req.Email, req.Password, req.FullName, req.CompanyName, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генеруємо JWT токен
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, h.authConfig)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Відправляємо відповідь
	response := AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login обробляє вхід користувача
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Перевіряємо облікові дані
	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Генеруємо JWT токен
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, h.authConfig)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Відправляємо відповідь
	response := AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Me повертає інформацію про поточного користувача
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("user_id").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ShowLoginPage відображає сторінку входу
func (h *AuthHandler) ShowLoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Title": "Вхід",
	})
}

// HandleLogin обробляє запит на вхід
func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var input models.LoginInput
	if err := c.BodyParser(&input); err != nil {
		log.Printf("Error parsing login input: %v", err)
		return c.Render("auth/login", fiber.Map{
			"Title": "Вхід",
			"Error": "Неправильний формат даних",
		})
	}

	log.Printf("Login attempt for email: %s", input.Email)
	user, err := h.userService.Login(c.Context(), input.Email, input.Password)
	if err != nil {
		log.Printf("Login failed for email %s: %v", input.Email, err)
		return c.Render("auth/login", fiber.Map{
			"Title": "Вхід",
			"Error": "Неправильний email або пароль",
		})
	}

	log.Printf("User authenticated successfully: %s (%s)", user.Email, user.ID.String())
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, h.authConfig)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return c.Render("auth/login", fiber.Map{
			"Title": "Вхід",
			"Error": "Помилка генерації токена",
		})
	}

	// Set cookie with longer expiration time
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		MaxAge:   86400 * 30, // 30 days in seconds
		Secure:   c.Protocol() == "https",
	})

	// Log successful login
	log.Printf("Successfully set session cookie for user %s", user.Email)
	log.Printf("Redirecting to dashboard...")

	// Redirect to dashboard
	return c.Redirect("/")
}

// ShowRegisterPage відображає сторінку реєстрації
func (h *AuthHandler) ShowRegisterPage(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"Title": "Реєстрація",
	})
}

// HandleRegister обробляє запит на реєстрацію
func (h *AuthHandler) HandleRegister(c *fiber.Ctx) error {
	var input models.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Render("auth/register", fiber.Map{
			"Title": "Реєстрація",
			"Error": "Неправильний формат даних",
		})
	}

	// Перевіряємо, чи вже існує користувач з таким email
	existingUser, err := h.userService.GetByEmail(c.Context(), input.Email)
	if err == nil && existingUser != nil {
		return c.Render("auth/register", fiber.Map{
			"Title": "Реєстрація",
			"Error": "Користувач з таким email вже існує",
		})
	}

	// Створюємо нового користувача
	user, err := h.userService.Register(c.Context(), input.Email, input.Password, input.FullName, input.CompanyName, "user")
	if err != nil {
		return c.Render("auth/register", fiber.Map{
			"Title": "Реєстрація",
			"Error": "Помилка при реєстрації: " + err.Error(),
		})
	}

	// Генеруємо токен
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, h.authConfig)
	if err != nil {
		return c.Render("auth/register", fiber.Map{
			"Title": "Реєстрація",
			"Error": "Помилка генерації токена",
		})
	}

	// Встановлюємо токен в cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    token,
		HTTPOnly: true,
	})

	return c.Redirect("/dashboard")
}

// HandleLogout обробляє запит на вихід
func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	c.ClearCookie("session")
	return c.Redirect("/login")
}

// GetJWTSecret повертає секретний ключ для JWT
func (h *AuthHandler) GetJWTSecret() string {
	return h.authConfig.SecretKey
}
