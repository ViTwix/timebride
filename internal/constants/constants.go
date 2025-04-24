package constants

import "time"

// Загальні константи
const (
	// DefaultPageSize визначає розмір сторінки за замовчуванням
	DefaultPageSize = 20
	// MaxPageSize визначає максимальний розмір сторінки
	MaxPageSize = 100
	// DefaultTimeout визначає таймаут за замовчуванням
	DefaultTimeout = 30 * time.Second
)

// Константи для файлів
const (
	// MaxFileSize визначає максимальний розмір файлу (100MB)
	MaxFileSize = 100 * 1024 * 1024
	// AllowedImageTypes визначає дозволені типи зображень
	AllowedImageTypes = "image/jpeg,image/png,image/gif"
	// AllowedVideoTypes визначає дозволені типи відео
	AllowedVideoTypes = "video/mp4,video/quicktime"
)

// Константи для кешування
const (
	// DefaultCacheTTL визначає час життя кешу за замовчуванням
	DefaultCacheTTL = 15 * time.Minute
	// LongCacheTTL визначає довгий час життя кешу
	LongCacheTTL = 24 * time.Hour
	// ShortCacheTTL визначає короткий час життя кешу
	ShortCacheTTL = 5 * time.Minute
)

// Константи для JWT
const (
	// AccessTokenExpiration визначає час життя access токена
	AccessTokenExpiration = 15 * time.Minute
	// RefreshTokenExpiration визначає час життя refresh токена
	RefreshTokenExpiration = 7 * 24 * time.Hour
)

// Константи для валідації
const (
	// MinPasswordLength визначає мінімальну довжину пароля
	MinPasswordLength = 8
	// MaxPasswordLength визначає максимальну довжину пароля
	MaxPasswordLength = 72
	// MinNameLength визначає мінімальну довжину імені
	MinNameLength = 2
	// MaxNameLength визначає максимальну довжину імені
	MaxNameLength = 50
)

// Константи для rate limiting
const (
	// RateLimitWindow визначає вікно для rate limiting
	RateLimitWindow = time.Minute
	// RateLimitRequests визначає кількість дозволених запитів
	RateLimitRequests = 60
)
