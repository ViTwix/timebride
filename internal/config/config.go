package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config містить всі налаштування програми
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

// ServerConfig містить налаштування сервера
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	Environment  string `mapstructure:"environment"`
	AllowOrigins string `mapstructure:"allow_origins"`
}

// DatabaseConfig містить налаштування бази даних
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// JWTConfig містить налаштування JWT
type JWTConfig struct {
	SecretKey     string        `mapstructure:"secret_key"`
	TokenDuration time.Duration `mapstructure:"token_duration"`
}

// StorageConfig містить налаштування сховища файлів
type StorageConfig struct {
	Provider  string          `mapstructure:"provider"`
	Backblaze BackblazeConfig `mapstructure:"backblaze"`
	CDN       CDNConfig       `mapstructure:"cdn"`
}

// BackblazeConfig містить налаштування Backblaze B2
type BackblazeConfig struct {
	AccountID      string `mapstructure:"account_id"`
	ApplicationKey string `mapstructure:"application_key"`
	Bucket         string `mapstructure:"bucket"`
	BucketID       string `mapstructure:"bucket_id"`
	Endpoint       string `mapstructure:"endpoint"`
	Region         string `mapstructure:"region"`
}

// CDNConfig містить налаштування CDN
type CDNConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Domain   string `mapstructure:"domain"`
	Protocol string `mapstructure:"protocol"`
}

// AuthConfig містить налаштування автентифікації
type AuthConfig struct {
	JWTSecret                string        `mapstructure:"jwt_secret"`
	AccessTokenExpiryMinutes int           `mapstructure:"access_token_expiry_minutes"`
	RefreshTokenExpiryDays   int           `mapstructure:"refresh_token_expiry_days"`
	Google                   OAuthProvider `mapstructure:"google"`
	Facebook                 OAuthProvider `mapstructure:"facebook"`
	Apple                    OAuthProvider `mapstructure:"apple"`
}

// OAuthProvider містить налаштування OAuth провайдера
type OAuthProvider struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

// Load завантажує конфігурацію з файлу
func Load() (*Config, error) {
	// Спочатку завантажуємо .env файл, перевіряючи різні можливі шляхи
	envPaths := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
		"./config/.env",
		"../config/.env",
		"../../config/.env",
	}

	var envLoaded bool
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		return nil, godotenv.Load() // спробуємо завантажити з поточної директорії як останній варіант
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../..")
	viper.AddConfigPath("../../config")

	// Встановлюємо значення за замовчуванням
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.allow_origins", "*")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.token_duration", "24h")
	viper.SetDefault("storage.backblaze.endpoint", "https://s3.eu-central-003.backblazeb2.com")
	viper.SetDefault("storage.backblaze.region", "eu-central-003")

	// Значення за замовчуванням для автентифікації
	viper.SetDefault("auth.access_token_expiry_minutes", 15)
	viper.SetDefault("auth.refresh_token_expiry_days", 7)

	// Зв'язуємо змінні середовища з конфігурацією
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")

	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.host", "SERVER_HOST")

	viper.BindEnv("storage.backblaze.account_id", "B2_ACCOUNT_ID")
	viper.BindEnv("storage.backblaze.application_key", "B2_APPLICATION_KEY")
	viper.BindEnv("storage.backblaze.bucket", "B2_BUCKET")
	viper.BindEnv("storage.backblaze.bucket_id", "B2_BUCKET_ID")

	// OAuth провайдери
	viper.BindEnv("auth.jwt_secret", "JWT_SECRET")
	viper.BindEnv("auth.google.client_id", "GOOGLE_CLIENT_ID")
	viper.BindEnv("auth.google.client_secret", "GOOGLE_CLIENT_SECRET")
	viper.BindEnv("auth.google.redirect_url", "GOOGLE_REDIRECT_URL")
	viper.BindEnv("auth.facebook.client_id", "FACEBOOK_CLIENT_ID")
	viper.BindEnv("auth.facebook.client_secret", "FACEBOOK_CLIENT_SECRET")
	viper.BindEnv("auth.facebook.redirect_url", "FACEBOOK_REDIRECT_URL")
	viper.BindEnv("auth.apple.client_id", "APPLE_CLIENT_ID")
	viper.BindEnv("auth.apple.client_secret", "APPLE_CLIENT_SECRET")
	viper.BindEnv("auth.apple.redirect_url", "APPLE_REDIRECT_URL")

	// Завантажуємо змінні середовища
	viper.AutomaticEnv()

	// Завантажуємо основний конфіг
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Якщо JWT секрет не встановлено, використовуємо значення з JWTConfig
	if config.Auth.JWTSecret == "" {
		config.Auth.JWTSecret = config.JWT.SecretKey
	}

	return &config, nil
}
