package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"timebride/internal/services/auth"
)

// Auth перевіряє JWT токен
func Auth(authService auth.IAuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Отримуємо токен з cookie
		tokenString := c.Cookies("token")
		if tokenString == "" {
			return c.Redirect("/login")
		}

		// Парсимо токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return authService.GetJWTSecret(), nil
		})
		if err != nil {
			return c.Redirect("/login")
		}

		// Перевіряємо валідність токена
		if !token.Valid {
			return c.Redirect("/login")
		}

		// Отримуємо claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Redirect("/login")
		}

		// Зберігаємо дані користувача в контексті
		c.Locals("user_id", claims["user_id"])
		c.Locals("email", claims["email"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
