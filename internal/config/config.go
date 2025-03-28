package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config yapımız için tüm ayarları tutan struct
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	MinIO    MinIOConfig
}

// ServerConfig sunucu ayarları
type ServerConfig struct {
	Port         string
	TemplateDir  string
	StaticDir    string
	UploadsDir   string
	MaxUploadMB  int
	Environment  string
	AllowOrigins string
}

// DatabaseConfig veritabanı ayarları
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig Redis önbellek ayarları
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig JWT ayarları
type JWTConfig struct {
	Secret          string
	AccessTokenExp  int // dakika cinsinden
	RefreshTokenExp int // saat cinsinden
}

// MinIOConfig MinIO nesne depolama ayarları
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	Location        string
}

// LoadConfig çevre değişkenlerinden yapılandırmayı yükler
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "3000"),
			TemplateDir:  getEnv("TEMPLATE_DIR", "./web/templates"),
			StaticDir:    getEnv("STATIC_DIR", "./web/static"),
			UploadsDir:   getEnv("UPLOADS_DIR", "./web/uploads"),
			MaxUploadMB:  getEnvAsInt("MAX_UPLOAD_MB", 10),
			Environment:  getEnv("ENVIRONMENT", "development"),
			AllowOrigins: getEnv("ALLOW_ORIGINS", "*"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "haberdb"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "change-this-secret-in-production"),
			AccessTokenExp:  getEnvAsInt("JWT_ACCESS_TOKEN_EXP", 60),   // 60 dakika
			RefreshTokenExp: getEnvAsInt("JWT_REFRESH_TOKEN_EXP", 168), // 7 gün
		},
		MinIO: MinIOConfig{
			Endpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKeyID:     getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretAccessKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			UseSSL:          getEnvAsBool("MINIO_USE_SSL", false),
			BucketName:      getEnv("MINIO_BUCKET_NAME", "haber"),
			Location:        getEnv("MINIO_LOCATION", "eu-west-1"),
		},
	}
}

// GetDSN PostgreSQL bağlantı dizesini döndürür
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Europe/Istanbul",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetRedisAddr Redis adresini döndürür
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Yardımcı fonksiyonlar
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
