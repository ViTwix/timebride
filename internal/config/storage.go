package config

// StorageConfig містить налаштування сховища файлів
type StorageConfig struct {
	Provider  string `yaml:"provider"` // local, s3, etc.
	Path      string `yaml:"path"`     // локальний шлях або bucket
	MaxSizeGB int    `yaml:"max_size_gb"`
	Region    string `yaml:"region"`   // для cloud storage
	Endpoint  string `yaml:"endpoint"` // для custom endpoints
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

// GetStorageProvider повертає тип провайдера сховища
func (c *StorageConfig) GetStorageProvider() string {
	if c.Provider == "" {
		return "local"
	}
	return c.Provider
}

// GetStoragePath повертає шлях до сховища
func (c *StorageConfig) GetStoragePath() string {
	if c.Path == "" {
		return "./storage"
	}
	return c.Path
}

// GetMaxStorageSize повертає максимальний розмір сховища в байтах
func (c *StorageConfig) GetMaxStorageSize() int64 {
	return int64(c.MaxSizeGB) * 1024 * 1024 * 1024 // Convert GB to bytes
}

// IsCloudStorage перевіряє чи використовується хмарне сховище
func (c *StorageConfig) IsCloudStorage() bool {
	return c.Provider != "local"
}

// GetEndpoint повертає endpoint для хмарного сховища
func (c *StorageConfig) GetEndpoint() string {
	if c.Endpoint == "" && c.Region != "" {
		return "https://s3." + c.Region + ".amazonaws.com"
	}
	return c.Endpoint
}

// HasCredentials перевіряє чи налаштовані облікові дані
func (c *StorageConfig) HasCredentials() bool {
	return c.AccessKey != "" && c.SecretKey != ""
}
