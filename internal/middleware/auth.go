package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"timebride/internal/config"
	"timebride/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims містить дані JWT токена
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// AuthConfig містить налаштування для JWT
type AuthConfig struct {
	SecretKey     string
	TokenDuration time.Duration
	UserService   services.UserService
	AuthService   services.AuthService
}

// AuthMiddleware перевіряє JWT токен в заголовку Authorization
func AuthMiddleware(config AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Отримуємо токен з заголовка
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Перевіряємо формат токена
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Парсимо та валідуємо токен
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.SecretKey), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Додаємо дані користувача в контекст
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "email", claims.Email)
			ctx = context.WithValue(ctx, "role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GenerateToken генерує JWT токен
func GenerateToken(userID uuid.UUID, email string, role string, config AuthConfig) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(config.TokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.SecretKey))
}

// RequireRole перевіряє, чи має користувач необхідну роль
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value("role").(string)

			// Перевіряємо, чи має користувач одну з необхідних ролей
			hasRole := false
			for _, role := range roles {
				if role == userRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Auth перевіряє JWT токен з cookie і додає дані користувача до контексту
func Auth(c *fiber.Ctx) error {
	log.Println("Auth middleware started")

	// Список публічних маршрутів, які не потребують аутентифікації
	publicPaths := []string{
		"/login",
		"/register",
		"/api/auth/login",
		"/api/auth/register",
		"/static",
		"/images",
		"/css",
		"/js",
		"/fonts",
		"/favicon.ico",
	}

	// Перевіряємо, чи поточний шлях є публічним
	path := c.Path()
	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			log.Printf("Public path detected: %s, skipping auth", path)
			return c.Next()
		}
	}

	log.Printf("Checking auth for path: %s", path)

	// Отримуємо токен з cookie
	token := c.Cookies("session")
	if token == "" {
		log.Println("No token found in cookies, redirecting to login")
		return c.Redirect("/login")
	}

	log.Printf("Token found: %s", token[:10]+"...")

	// Отримуємо секретний ключ з контексту
	secretKey := c.Locals("jwt_secret")
	if secretKey == nil {
		log.Println("JWT secret not found in context, redirecting to login")
		return c.Redirect("/login")
	}

	secretKeyStr, ok := secretKey.(string)
	if !ok {
		log.Printf("JWT secret is not a string, got: %T", secretKey)
		return c.Redirect("/login")
	}

	// Парсимо JWT токен
	claims, err := ParseJWT(token, secretKeyStr)
	if err != nil {
		log.Printf("Error parsing JWT: %v", err)
		c.ClearCookie("session")
		return c.Redirect("/login")
	}

	// Отримуємо ID користувача з токена (використовуємо поле "sub", яке містить ID)
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		log.Printf("Failed to get user ID from claims (sub), got type: %T, value: %v", claims["sub"], claims["sub"])
		c.ClearCookie("session")
		return c.Redirect("/login")
	}

	log.Printf("User ID from token: %s (type: %T)", userIDStr, userIDStr)

	// Зберігаємо ID користувача і email в контексті
	c.Locals("user_id", userIDStr) // Зберігаємо як рядок
	log.Printf("Saved user_id to locals: %s (type: %T)", userIDStr, userIDStr)

	if email, ok := claims["email"].(string); ok {
		c.Locals("email", email)
	}
	if role, ok := claims["role"].(string); ok {
		c.Locals("role", role)
	}

	log.Println("Auth middleware completed successfully")
	return c.Next()
}

// ParseJWT перевіряє і повертає дані з JWT токена
func ParseJWT(tokenString string, secretKey string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Перевіряємо метод підпису
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Перевіряємо токен і отримуємо claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["exp"] != nil {
			exp := int64(claims["exp"].(float64))
			if time.Unix(exp, 0).Before(time.Now()) {
				return nil, fmt.Errorf("token expired")
			}
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// JWTMiddleware створює middleware для перевірки JWT токенів
func JWTMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Отримуємо токен з заголовка Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Перевіряємо формат Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format, must be 'Bearer {token}'",
			})
		}

		tokenString := parts[1]

		// Парсимо та перевіряємо токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Перевіряємо метод підпису
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(cfg.Auth.JWTSecret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid or expired token: %v", err),
			})
		}

		// Перевіряємо валідність токена
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Отримуємо claims з токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Витягуємо дані користувача з claims
		userID, ok := claims["sub"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}

		// Парсимо UUID
		uid, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}

		email, ok := claims["email"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email in token",
			})
		}

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid role in token",
			})
		}

		// Зберігаємо дані користувача в контексті для використання в обробниках
		c.Locals("userID", uid)
		c.Locals("email", email)
		c.Locals("role", role)

		return c.Next()
	}
}

// RoleAuthMiddleware створює middleware для перевірки ролі користувача
func RoleAuthMiddleware(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Отримуємо роль користувача з контексту (встановлену JWTMiddleware)
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User role not found in context",
			})
		}

		// Перевіряємо, чи має користувач необхідну роль
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: insufficient permissions",
			})
		}

		return c.Next()
	}
}
