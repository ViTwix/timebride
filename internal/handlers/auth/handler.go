package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"timebride/internal/models"
	"timebride/internal/services/auth"
	"timebride/internal/types"
)

// Handler реалізує обробку запитів аутентифікації
type Handler struct {
	authService auth.IAuthService
}

// NewHandler створює новий екземпляр обробника аутентифікації
func NewHandler(authService auth.IAuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

// ShowLoginPage відображає сторінку входу
func (h *Handler) ShowLoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Title": "Вхід",
	})
}

// ShowRegisterPage відображає сторінку реєстрації
func (h *Handler) ShowRegisterPage(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"Title": "Реєстрація",
	})
}

// HandleLogin обробляє запит на вхід
func (h *Handler) HandleLogin(c *fiber.Ctx) error {
	var input models.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	user, tokens, err := h.authService.Login(c.Context(), input.Email, input.Password)
	if err != nil {
		if c.XHR() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusUnauthorized).Render("auth/login", fiber.Map{
			"Title": "Вхід",
			"Error": err.Error(),
		})
	}

	// Встановлюємо токени в куки
	h.setAuthCookies(c, tokens)

	if c.XHR() {
		return c.JSON(fiber.Map{
			"message": "Успішний вхід",
			"user":    user,
			"tokens":  tokens,
		})
	}
	return c.Redirect("/app")
}

// HandleRegister обробляє запит на реєстрацію
func (h *Handler) HandleRegister(c *fiber.Ctx) error {
	var input models.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	user, err := h.authService.Register(c.Context(), input.Email, input.Password, input.FullName)
	if err != nil {
		if c.XHR() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).Render("auth/register", fiber.Map{
			"Title": "Реєстрація",
			"Error": err.Error(),
		})
	}

	// Автоматично логінимо користувача після реєстрації
	_, tokens, err := h.authService.Login(c.Context(), input.Email, input.Password)
	if err != nil {
		return err
	}

	// Встановлюємо токени в куки
	h.setAuthCookies(c, tokens)

	if c.XHR() {
		return c.JSON(fiber.Map{
			"message": "Успішна реєстрація",
			"user":    user,
			"tokens":  tokens,
		})
	}
	return c.Redirect("/app")
}

// HandleLogout обробляє запит на вихід
func (h *Handler) HandleLogout(c *fiber.Ctx) error {
	// Видаляємо куки
	h.clearAuthCookies(c)

	if c.XHR() {
		return c.JSON(fiber.Map{
			"message": "Успішний вихід",
		})
	}
	return c.Redirect("/login")
}

// HandleRefreshToken обробляє запит на оновлення токену
func (h *Handler) HandleRefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return fiber.ErrUnauthorized
	}

	_, tokens, err := h.authService.RefreshToken(c.Context(), refreshToken)
	if err != nil {
		return err
	}

	// Встановлюємо нові токени в куки
	h.setAuthCookies(c, tokens)

	return c.JSON(fiber.Map{
		"message": "Токени оновлено",
		"tokens":  tokens,
	})
}

// HandleOAuthRedirect обробляє редірект на OAuth провайдера
func (h *Handler) HandleOAuthRedirect(c *fiber.Ctx) error {
	provider := c.Params("provider")
	url, state, err := h.authService.GenerateOAuthURL(provider)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.Redirect(url)
}

// HandleOAuthCallback обробляє відповідь від OAuth провайдера
func (h *Handler) HandleOAuthCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	code := c.Query("code")
	state := c.Query("state")
	savedState := c.Cookies("oauth_state")

	user, tokens, err := h.authService.HandleOAuthCallback(c.Context(), provider, code, state, savedState)
	if err != nil {
		return err
	}

	// Видаляємо oauth_state куку
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	// Встановлюємо токени в куки
	h.setAuthCookies(c, tokens)

	if c.XHR() {
		return c.JSON(fiber.Map{
			"message": "Успішна OAuth автентифікація",
			"user":    user,
			"tokens":  tokens,
		})
	}
	return c.Redirect("/app")
}

// OAuthCallback обробляє callback від OAuth провайдера
func (h *Handler) OAuthCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	code := c.Query("code")
	state := c.Query("state")
	savedState := c.Cookies("oauth_state")

	user, tokens, err := h.authService.HandleOAuthCallback(c.Context(), provider, code, state, savedState)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Встановлюємо токени в куки
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Expires:  tokens.ExpiresAt,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  tokens.ExpiresAt.Add(24 * time.Hour * 30), // 30 днів
		HTTPOnly: true,
	})

	return c.JSON(user)
}

// OAuthRedirect перенаправляє на сторінку OAuth провайдера
func (h *Handler) OAuthRedirect(c *fiber.Ctx) error {
	provider := c.Params("provider")
	url, state, err := h.authService.GenerateOAuthURL(provider)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Зберігаємо state в куки
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
	})

	return c.Redirect(url)
}

// setAuthCookies встановлює токени в куки
func (h *Handler) setAuthCookies(c *fiber.Ctx, tokens *types.AuthTokens) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		MaxAge:   900, // 15 minutes
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		MaxAge:   604800, // 7 days
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})
}

// clearAuthCookies видаляє токени з кук
func (h *Handler) clearAuthCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})
}
