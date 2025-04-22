package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"

	"timebride/internal/config"
	"timebride/internal/models"
	"timebride/internal/repositories"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidOAuthState    = errors.New("invalid oauth state")
	ErrOAuthProviderError   = errors.New("oauth provider error")
	ErrProviderNotSupported = errors.New("auth provider not supported")
)

// AuthService представляє сервіс для автентифікації
type AuthService struct {
	userRepo       repositories.UserRepository
	config         *config.Config
	googleConfig   *oauth2.Config
	facebookConfig *oauth2.Config
	appleConfig    *oauth2.Config
}

// AuthTokens містить токени авторизації
type AuthTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// NewAuthService створює новий сервіс аутентифікації
func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config) *AuthService {
	// Налаштування OAuth для Google
	googleConfig := &oauth2.Config{
		ClientID:     cfg.Auth.Google.ClientID,
		ClientSecret: cfg.Auth.Google.ClientSecret,
		RedirectURL:  cfg.Auth.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Налаштування OAuth для Facebook
	facebookConfig := &oauth2.Config{
		ClientID:     cfg.Auth.Facebook.ClientID,
		ClientSecret: cfg.Auth.Facebook.ClientSecret,
		RedirectURL:  cfg.Auth.Facebook.RedirectURL,
		Scopes: []string{
			"email",
			"public_profile",
		},
		Endpoint: facebook.Endpoint,
	}

	// Налаштування для Apple Sign In
	// Apple використовує інший механізм, тому просто зберігаємо базову конфігурацію
	appleConfig := &oauth2.Config{
		ClientID:     cfg.Auth.Apple.ClientID,
		ClientSecret: cfg.Auth.Apple.ClientSecret,
		RedirectURL:  cfg.Auth.Apple.RedirectURL,
		Scopes: []string{
			"name",
			"email",
		},
		// Apple не має стандартного Endpoint в oauth2, налаштовуємо власні URL
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://appleid.apple.com/auth/authorize",
			TokenURL: "https://appleid.apple.com/auth/token",
		},
	}

	return &AuthService{
		userRepo:       userRepo,
		config:         cfg,
		googleConfig:   googleConfig,
		facebookConfig: facebookConfig,
		appleConfig:    appleConfig,
	}
}

// Register реєструє нового користувача з локальними обліковими даними
func (s *AuthService) Register(ctx context.Context, email, password, name string) (*models.User, error) {
	// Перевіряємо, чи існує вже користувач з таким email
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Хешуємо пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Конвертуємо дозволи в JSON
	permissions, err := json.Marshal(models.DefaultPermissions(models.RoleUser))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal permissions: %w", err)
	}

	// Створюємо нового користувача
	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		AuthProvider: models.AuthProviderLocal,
		Name:         name,
		Role:         models.RoleUser, // За замовчуванням звичайний користувач
		Permissions:  permissions,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Зберігаємо користувача в базі даних
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login автентифікує користувача за допомогою email та пароля
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, *AuthTokens, error) {
	// Отримуємо користувача за email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Перевіряємо, чи використовує користувач локальну аутентифікацію
	if user.AuthProvider != models.AuthProviderLocal {
		return nil, nil, fmt.Errorf("this account uses %s authentication", user.AuthProvider)
	}

	// Перевіряємо пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Генеруємо токени
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

// GenerateOAuthURL генерує URL для OAuth аутентифікації
func (s *AuthService) GenerateOAuthURL(provider string) (string, string, error) {
	// Генеруємо випадковий стан для безпеки
	state := generateRandomState()

	var authURL string
	switch provider {
	case models.AuthProviderGoogle:
		authURL = s.googleConfig.AuthCodeURL(state)
	case models.AuthProviderFacebook:
		authURL = s.facebookConfig.AuthCodeURL(state)
	case models.AuthProviderApple:
		authURL = s.appleConfig.AuthCodeURL(state)
	default:
		return "", "", ErrProviderNotSupported
	}

	return authURL, state, nil
}

// HandleOAuthCallback обробляє відповідь від OAuth провайдера
func (s *AuthService) HandleOAuthCallback(ctx context.Context, provider, code, state, savedState string) (*models.User, *AuthTokens, error) {
	// Перевіряємо стан для безпеки
	if state != savedState {
		return nil, nil, ErrInvalidOAuthState
	}

	var oauthConfig *oauth2.Config
	switch provider {
	case models.AuthProviderGoogle:
		oauthConfig = s.googleConfig
	case models.AuthProviderFacebook:
		oauthConfig = s.facebookConfig
	case models.AuthProviderApple:
		oauthConfig = s.appleConfig
	default:
		return nil, nil, ErrProviderNotSupported
	}

	// Обмінюємо код на токен
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("code exchange failed: %w", err)
	}

	// Отримуємо інформацію про користувача з провайдера
	userInfo, err := s.getUserInfoFromProvider(ctx, provider, token)
	if err != nil {
		return nil, nil, err
	}

	// Перевіряємо, чи існує вже користувач з таким email
	existingUser, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
	if err == nil && existingUser != nil {
		// Якщо користувач вже існує, перевіряємо провайдера
		if existingUser.AuthProvider != provider {
			// Користувач вже зареєстрований з іншим провайдером
			return nil, nil, fmt.Errorf("email already registered with %s", existingUser.AuthProvider)
		}

		// Успішний логін через OAuth, генеруємо токени
		tokens, err := s.generateTokens(existingUser)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
		}

		return existingUser, tokens, nil
	}

	// Конвертуємо дозволи в JSON
	permissions, err := json.Marshal(models.DefaultPermissions(models.RoleUser))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal permissions: %w", err)
	}

	// Створюємо нового користувача
	providerID := userInfo.ID
	newUser := &models.User{
		ID:           uuid.New(),
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		AuthProvider: provider,
		ProviderID:   &providerID,
		Role:         models.RoleUser,
		Permissions:  permissions,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if userInfo.Picture != "" {
		newUser.AvatarURL = &userInfo.Picture
	}

	// Зберігаємо користувача в базі даних
	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Генеруємо токени
	tokens, err := s.generateTokens(newUser)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return newUser, tokens, nil
}

// getUserInfoFromProvider отримує інформацію про користувача від OAuth провайдера
func (s *AuthService) getUserInfoFromProvider(ctx context.Context, provider string, token *oauth2.Token) (*OAuthUserInfo, error) {
	switch provider {
	case models.AuthProviderGoogle:
		return s.getGoogleUserInfo(ctx, token)
	case models.AuthProviderFacebook:
		return s.getFacebookUserInfo(ctx, token)
	case models.AuthProviderApple:
		return s.getAppleUserInfo(ctx, token)
	default:
		return nil, ErrProviderNotSupported
	}
}

// OAuthUserInfo представляє інформацію про користувача, отриману від OAuth провайдера
type OAuthUserInfo struct {
	ID      string
	Email   string
	Name    string
	Picture string
}

// getGoogleUserInfo отримує інформацію користувача з Google
func (s *AuthService) getGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	// Створюємо клієнт з токеном
	client := s.googleConfig.Client(ctx, token)

	// Запитуємо інформацію про користувача
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google API response: %s", resp.Status)
	}

	// Зчитуємо відповідь
	var googleUser struct {
		Sub     string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode Google user info: %w", err)
	}

	return &OAuthUserInfo{
		ID:      googleUser.Sub,
		Email:   googleUser.Email,
		Name:    googleUser.Name,
		Picture: googleUser.Picture,
	}, nil
}

// getFacebookUserInfo отримує інформацію користувача з Facebook
func (s *AuthService) getFacebookUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	// Створюємо клієнт з токеном
	client := s.facebookConfig.Client(ctx, token)

	// Запитуємо інформацію про користувача, включаючи email
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture.type(large)")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Facebook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Facebook API response: %s", resp.Status)
	}

	// Зчитуємо відповідь
	var fbUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&fbUser); err != nil {
		return nil, fmt.Errorf("failed to decode Facebook user info: %w", err)
	}

	return &OAuthUserInfo{
		ID:      fbUser.ID,
		Email:   fbUser.Email,
		Name:    fbUser.Name,
		Picture: fbUser.Picture.Data.URL,
	}, nil
}

// getAppleUserInfo отримує інформацію користувача з Apple
func (s *AuthService) getAppleUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUserInfo, error) {
	// Перевіряємо, чи контекст не скасований
	if ctx.Err() != nil {
		return nil, fmt.Errorf("context error: %w", ctx.Err())
	}

	// Apple не надає стандартний endpoint для отримання інформації про користувача
	// Основна інформація приходить в JWT токені
	// Декодуємо ідентифікаційний токен (id_token) для отримання інформації

	// Отримуємо id_token з відповіді
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token not found in Apple response")
	}

	// Парсимо JWT без перевірки підпису (Apple вже перевірив під час обміну коду на токен)
	// У реальному додатку вам потрібно правильно перевірити підпис JWT
	tokenParts := strings.Split(idToken, ".")
	if len(tokenParts) < 2 {
		return nil, errors.New("invalid JWT format from Apple")
	}

	// Декодуємо payload (другу частину JWT)
	payload, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	// Парсимо JSON-структуру
	var claims struct {
		Subject string `json:"sub"`
		Email   string `json:"email"`
	}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse JWT claims: %w", err)
	}

	// Apple передає ім'я користувача тільки при першій авторизації в полі user (в запиті),
	// тому ми використовуємо email як ім'я, якщо не маємо імені
	// У реальному додатку вам потрібно зберігати це ім'я, якщо воно приходить
	name := claims.Email
	if idx := strings.Index(name, "@"); idx > 0 {
		name = name[:idx] // Використовуємо частину до @ як ім'я
	}

	return &OAuthUserInfo{
		ID:    claims.Subject,
		Email: claims.Email,
		Name:  name,
		// Apple не надає URL для зображення профілю
	}, nil
}

// generateTokens генерує пару access і refresh токенів
func (s *AuthService) generateTokens(user *models.User) (*AuthTokens, error) {
	// Поточний час
	now := time.Now()

	// Час закінчення токена
	accessExpiresAt := now.Add(time.Minute * time.Duration(s.config.Auth.AccessTokenExpiryMinutes))
	refreshExpiresAt := now.Add(time.Hour * 24 * time.Duration(s.config.Auth.RefreshTokenExpiryDays))

	// Claims для JWT
	accessClaims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   accessExpiresAt.Unix(),
		"iat":   now.Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"exp":   refreshExpiresAt.Unix(),
		"iat":   now.Unix(),
		"token": generateRandomString(32),
	}

	// Створюємо токени
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Підписуємо токени
	accessTokenString, err := accessToken.SignedString([]byte(s.config.Auth.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.Auth.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	// Формуємо результат
	return &AuthTokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiresAt,
	}, nil
}

// RefreshToken оновлює токен доступу за допомогою refresh токена
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (*models.User, *AuthTokens, error) {
	// Парсимо refresh токен
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		// Перевіряємо алгоритм підпису
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Auth.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, nil, errors.New("invalid refresh token")
	}

	// Отримуємо claims з токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, errors.New("invalid token claims")
	}

	// Перевіряємо термін дії
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, nil, errors.New("invalid token expiration")
	}

	// Перевіряємо, чи не сплив термін дії
	if time.Now().Unix() > int64(exp) {
		return nil, nil, errors.New("refresh token expired")
	}

	// Отримуємо ID користувача
	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, nil, errors.New("invalid user ID in token")
	}

	// Парсимо UUID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Отримуємо користувача з бази даних
	user, err := s.userRepo.GetByID(ctx, uid)
	if err != nil {
		return nil, nil, ErrUserNotFound
	}

	// Генеруємо нові токени
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

// generateRandomState генерує випадковий рядок для стану OAuth
func generateRandomState() string {
	return generateRandomString(32)
}

// generateRandomString генерує випадковий рядок заданої довжини
func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		// У випадку помилки використовуємо поточний час як fallback
		// Це не так безпечно, але краще ніж нічого
		for i := range b {
			b[i] = byte(time.Now().Nanosecond() % 256)
			time.Sleep(time.Nanosecond)
		}
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// GetJWTSecret returns the JWT secret key
func (s *AuthService) GetJWTSecret() string {
	return s.config.Auth.JWTSecret
}
