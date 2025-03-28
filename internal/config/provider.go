package config

import (
	"sync"
)

// ConfigProvider, konfigürasyon için tek bir global provider sağlar
type ConfigProvider struct {
	config IConfig
	mu     sync.RWMutex
}

var (
	// Singleton pattern için tek bir provider instance'ı
	provider *ConfigProvider
	once     sync.Once
)

// GetProvider ConfigProvider singleton instance'ını döndürür
func GetProvider() *ConfigProvider {
	once.Do(func() {
		provider = &ConfigProvider{
			config: LoadConfig(),
		}
	})
	return provider
}

// GetConfig mevcut konfigürasyonu döndürür
func (p *ConfigProvider) GetConfig() IConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.config
}

// SetConfig yeni bir konfigürasyon ayarlar (test için kullanışlı)
func (p *ConfigProvider) SetConfig(config IConfig) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.config = config
}

// MockConfig test için mock bir config oluşturur
func MockConfig() IConfig {
	return &Config{
		Server: ServerConfig{
			Port:         "3000",
			TemplateDir:  "./web/templates",
			StaticDir:    "./web/static",
			UploadsDir:   "./web/uploads",
			MaxUploadMB:  10,
			Environment:  "test",
			AllowOrigins: "*",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "haberdb_test",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:          "test-secret",
			AccessTokenExp:  60,
			RefreshTokenExp: 168,
		},
		MinIO: MinIOConfig{
			Endpoint:        "localhost:9000",
			AccessKeyID:     "minioadmin",
			SecretAccessKey: "minioadmin",
			UseSSL:          false,
			BucketName:      "haber-test",
			Location:        "eu-west-1",
		},
	}
}
