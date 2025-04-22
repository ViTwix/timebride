package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"timebride/internal/models"
	"timebride/internal/services"
)

// AuthHandler містить обробники запитів для автентифікації
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler створює новий обробник автентифікації
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest структура запиту на реєстрацію
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
}

// LoginRequest структура запиту на вхід
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest структура запиту на оновлення токену
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse структура відповіді з токенами
type AuthResponse struct {
	User         models.PublicUser `json:"user"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresAt    time.Time         `json:"expires_at"`
}

// Register обробляє запит на реєстрацію користувача
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валідація даних
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email, password and name are required",
		})
	}

	// Викликаємо сервіс для реєстрації
	user, err := h.authService.Register(c.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "Email already registered",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	// Генеруємо токени для автоматичного входу після реєстрації
	loginUser, tokens, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		// В цьому випадку користувач створений, але не можемо увійти
		// Повертаємо успіх реєстрації, але без токенів
		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully. Please login.",
			"user":    user.ToPublicUser(),
		})
	}

	// Повертаємо токени і дані користувача
	return c.Status(http.StatusCreated).JSON(AuthResponse{
		User:         loginUser.ToPublicUser(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	})
}

// Login обробляє запит на вхід користувача
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валідація даних
	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Викликаємо сервіс для логіну
	user, tokens, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}

	// Повертаємо токени і дані користувача
	return c.Status(http.StatusOK).JSON(AuthResponse{
		User:         user.ToPublicUser(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	})
}

// RefreshToken обробляє запит на оновлення токену
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валідація даних
	if req.RefreshToken == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	// Викликаємо сервіс для оновлення токену
	user, tokens, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired refresh token",
		})
	}

	// Повертаємо нові токени і дані користувача
	return c.Status(http.StatusOK).JSON(AuthResponse{
		User:         user.ToPublicUser(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	})
}

// OAuthRedirect генерує URL для OAuth автентифікації і перенаправляє користувача
func (h *AuthHandler) OAuthRedirect(c *fiber.Ctx) error {
	provider := c.Params("provider")

	// Перевіряємо, чи підтримується провайдер
	switch provider {
	case models.AuthProviderGoogle, models.AuthProviderFacebook, models.AuthProviderApple:
		// Продовжуємо
	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Unsupported OAuth provider",
		})
	}

	// Генеруємо URL для OAuth
	authURL, state, err := h.authService.GenerateOAuthURL(provider)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate OAuth URL",
		})
	}

	// Зберігаємо стан в сесії
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(time.Minute * 10), // Термін дії 10 хвилин
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "lax",
	})

	// Перенаправляємо користувача на URL автентифікації провайдера
	return c.Redirect(authURL, http.StatusFound)
}

// OAuthCallback обробляє відповідь від OAuth провайдера
func (h *AuthHandler) OAuthCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")

	// Перевіряємо, чи підтримується провайдер
	switch provider {
	case models.AuthProviderGoogle, models.AuthProviderFacebook, models.AuthProviderApple:
		// Продовжуємо
	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Unsupported OAuth provider",
		})
	}

	// Отримуємо код і стан з запиту
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid OAuth callback parameters",
		})
	}

	// Отримуємо збережений стан з cookie
	savedState := c.Cookies("oauth_state")
	if savedState == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing OAuth state",
		})
	}

	// Видаляємо cookie стану
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "lax",
	})

	// Перевіряємо стан для безпеки
	if state != savedState {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid OAuth state",
		})
	}

	// Обробляємо відповідь від OAuth провайдера
	user, tokens, err := h.authService.HandleOAuthCallback(c.Context(), provider, code, state, savedState)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate with provider: " + err.Error(),
		})
	}

	// Повертаємо токени і дані користувача
	return c.Status(http.StatusOK).JSON(AuthResponse{
		User:         user.ToPublicUser(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	})
}

// ShowLoginPage відображає сторінку входу
func (h *AuthHandler) ShowLoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Title": "Вхід в систему",
	})
}

// HandleLogin обробляє вхід користувача
func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	// Отримуємо дані з форми
	email := c.FormValue("email")
	password := c.FormValue("password")

	log.Printf("Спроба входу для користувача: %s", email)

	// Автентифікуємо користувача
	user, tokens, err := h.authService.Login(c.Context(), email, password)
	if err != nil {
		log.Printf("Помилка входу для %s: %v", email, err)
		// У випадку помилки повертаємося на сторінку логіну з повідомленням про помилку
		return c.Status(fiber.StatusUnauthorized).Render("login", fiber.Map{
			"Error": "Неправильний email або пароль",
			"Email": email,
		})
	}

	log.Printf("Успішний вхід для користувача: %s (ID: %s, Роль: %s)", email, user.ID.String(), user.Role)

	// Очищуємо старий cookie, якщо він є
	c.ClearCookie("session")

	// Встановлюємо токен в cookie
	cookie := fiber.Cookie{
		Name:     "session",
		Value:    tokens.AccessToken,
		Path:     "/",
		MaxAge:   int(time.Until(tokens.ExpiresAt).Seconds()),
		HTTPOnly: true,
		SameSite: "Lax",
	}
	c.Cookie(&cookie)

	// Безпечно логуємо перші символи токена
	tokenPreview := tokens.AccessToken
	if len(tokenPreview) > 20 {
		tokenPreview = tokenPreview[:20]
	}
	log.Printf("Встановлено cookie session з токеном (початок): %s...", tokenPreview)

	// Перенаправляємо на дашборд
	return c.Redirect("/app/dashboard")
}

// ShowRegisterPage відображає сторінку реєстрації
func (h *AuthHandler) ShowRegisterPage(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"Title": "Реєстрація",
	})
}

// HandleRegister обробляє реєстрацію користувача (переадресація на веб-інтерфейс)
func (h *AuthHandler) HandleRegister(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валідація даних
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Render("auth/register", fiber.Map{
			"Title": "Реєстрація",
			"Error": "Всі поля обов'язкові",
			"Email": req.Email,
			"Name":  req.Name,
		})
	}

	// Викликаємо сервіс для реєстрації
	_, err := h.authService.Register(c.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		errorMsg := "Помилка при реєстрації"
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			errorMsg = "Користувач з таким email вже існує"
		}

		return c.Render("auth/register", fiber.Map{
			"Title": "Реєстрація",
			"Error": errorMsg,
			"Email": req.Email,
			"Name":  req.Name,
		})
	}

	// Перенаправляємо на сторінку входу з повідомленням про успіх
	return c.Render("auth/login", fiber.Map{
		"Title":   "Вхід в систему",
		"Success": "Реєстрація успішна! Тепер ви можете увійти.",
		"Email":   req.Email,
	})
}

// GetJWTSecret повертає секретний ключ для JWT
func (h *AuthHandler) GetJWTSecret() string {
	return h.authService.GetJWTSecret()
}

// HandleLogout обробляє вихід користувача
func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	// Видаляємо cookie з токеном
	c.ClearCookie("session")

	// Перенаправляємо на сторінку входу
	return c.Redirect("/login")
}
