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

	// Зв'язуємо змінні середовища з конфігурацією
	viper.BindEnv("storage.backblaze.account_id", "B2_ACCOUNT_ID")
	viper.BindEnv("storage.backblaze.application_key", "B2_APPLICATION_KEY")
	viper.BindEnv("storage.backblaze.bucket", "B2_BUCKET")
	viper.BindEnv("storage.backblaze.bucket_id", "B2_BUCKET_ID")

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

	return &config, nil
}
