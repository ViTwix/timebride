package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config містить всі налаштування програми
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Storage  StorageConfig  `yaml:"storage"`
}

// ServerConfig містить налаштування сервера
type ServerConfig struct {
	Address        string        `yaml:"address"`
	CorsOrigins    []string      `yaml:"cors_origins"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	MaxHeaderBytes int           `yaml:"max_header_bytes"`
}

// DatabaseConfig містить налаштування бази даних
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}

// JWTConfig містить налаштування JWT
type JWTConfig struct {
	Secret                  string        `yaml:"secret"`
	AccessExpirationMinutes time.Duration `yaml:"access_expiration_minutes"`
	RefreshExpirationDays   time.Duration `yaml:"refresh_expiration_days"`
	Issuer                  string        `yaml:"issuer"`
	Audience                string        `yaml:"audience"`
}

// Load завантажує конфігурацію з .env файлу та змінних середовища
func Load() (*Config, error) {
	// Завантажуємо .env файл, якщо він існує
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Address:        getEnv("SERVER_ADDRESS", ":3000"),
			CorsOrigins:    []string{getEnv("CORS_ORIGINS", "*")},
			ReadTimeout:    time.Duration(getEnvInt("SERVER_READ_TIMEOUT", 60)) * time.Second,
			WriteTimeout:   time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT", 60)) * time.Second,
			MaxHeaderBytes: 1 << 20, // 1MB
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "timebride"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:                  getEnv("JWT_SECRET", "your-secret-key"),
			AccessExpirationMinutes: time.Duration(getEnvInt("JWT_ACCESS_EXPIRATION_MINUTES", 15)),
			RefreshExpirationDays:   time.Duration(getEnvInt("JWT_REFRESH_EXPIRATION_DAYS", 7)),
			Issuer:                  getEnv("JWT_ISSUER", "timebride"),
			Audience:                getEnv("JWT_AUDIENCE", "timebride-api"),
		},
		Storage: StorageConfig{
			Provider:  getEnv("STORAGE_PROVIDER", "local"),
			Path:      getEnv("STORAGE_PATH", "./storage"),
			MaxSizeGB: getEnvInt("STORAGE_MAX_SIZE_GB", 100),
			Region:    getEnv("STORAGE_REGION", ""),
			Endpoint:  getEnv("STORAGE_ENDPOINT", ""),
			AccessKey: getEnv("STORAGE_ACCESS_KEY", ""),
			SecretKey: getEnv("STORAGE_SECRET_KEY", ""),
		},
	}, nil
}

// getEnv отримує значення змінної середовища або повертає значення за замовчуванням
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt отримує числове значення змінної середовища або повертає значення за замовчуванням
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
