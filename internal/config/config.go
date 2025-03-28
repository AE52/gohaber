package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config yapımız için tüm ayarları tutan struct
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	MinIO       MinIOConfig
	RateLimiter RateLimiterConfig
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

// RateLimiterConfig Rate Limiter ayarları
type RateLimiterConfig struct {
	Enabled        bool
	MaxRequests    int
	ExpireSeconds  int
	SkipSuccessful bool
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
		RateLimiter: RateLimiterConfig{
			Enabled:        getEnvAsBool("RATE_LIMITER_ENABLED", true),
			MaxRequests:    getEnvAsInt("RATE_LIMITER_MAX_REQUESTS", 100),
			ExpireSeconds:  getEnvAsInt("RATE_LIMITER_EXPIRE_SECONDS", 60),
			SkipSuccessful: getEnvAsBool("RATE_LIMITER_SKIP_SUCCESSFUL", false),
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

// IConfig implentasyonu için getter metotları
func (c *Config) GetServer() IServerConfig {
	return &c.Server
}

func (c *Config) GetDatabase() IDatabaseConfig {
	return &c.Database
}

func (c *Config) GetRedis() IRedisConfig {
	return &c.Redis
}

func (c *Config) GetJWT() IJWTConfig {
	return &c.JWT
}

func (c *Config) GetMinIO() IMinIOConfig {
	return &c.MinIO
}

func (c *Config) GetRateLimiter() IRateLimiterConfig {
	return &c.RateLimiter
}

// IServerConfig implentasyonu için getter metotları
func (c *ServerConfig) GetPort() string {
	return c.Port
}

func (c *ServerConfig) GetTemplateDir() string {
	return c.TemplateDir
}

func (c *ServerConfig) GetStaticDir() string {
	return c.StaticDir
}

func (c *ServerConfig) GetUploadsDir() string {
	return c.UploadsDir
}

func (c *ServerConfig) GetMaxUploadMB() int {
	return c.MaxUploadMB
}

func (c *ServerConfig) GetEnvironment() string {
	return c.Environment
}

func (c *ServerConfig) GetAllowOrigins() string {
	return c.AllowOrigins
}

// IDatabaseConfig implentasyonu için getter metotları
func (c *DatabaseConfig) GetHost() string {
	return c.Host
}

func (c *DatabaseConfig) GetPort() string {
	return c.Port
}

func (c *DatabaseConfig) GetUser() string {
	return c.User
}

func (c *DatabaseConfig) GetPassword() string {
	return c.Password
}

func (c *DatabaseConfig) GetDBName() string {
	return c.DBName
}

func (c *DatabaseConfig) GetSSLMode() string {
	return c.SSLMode
}

// IRedisConfig implentasyonu için getter metotları
func (c *RedisConfig) GetHost() string {
	return c.Host
}

func (c *RedisConfig) GetPort() string {
	return c.Port
}

func (c *RedisConfig) GetPassword() string {
	return c.Password
}

func (c *RedisConfig) GetDB() int {
	return c.DB
}

// IJWTConfig implentasyonu için getter metotları
func (c *JWTConfig) GetSecret() string {
	return c.Secret
}

func (c *JWTConfig) GetAccessTokenExp() int {
	return c.AccessTokenExp
}

func (c *JWTConfig) GetRefreshTokenExp() int {
	return c.RefreshTokenExp
}

// IMinIOConfig implentasyonu için getter metotları
func (c *MinIOConfig) GetEndpoint() string {
	return c.Endpoint
}

func (c *MinIOConfig) GetAccessKeyID() string {
	return c.AccessKeyID
}

func (c *MinIOConfig) GetSecretAccessKey() string {
	return c.SecretAccessKey
}

func (c *MinIOConfig) GetUseSSL() bool {
	return c.UseSSL
}

func (c *MinIOConfig) GetBucketName() string {
	return c.BucketName
}

func (c *MinIOConfig) GetLocation() string {
	return c.Location
}

// IRateLimiterConfig implementasyonu için getter metotları
func (c *RateLimiterConfig) GetEnabled() bool {
	return c.Enabled
}

func (c *RateLimiterConfig) GetMaxRequests() int {
	return c.MaxRequests
}

func (c *RateLimiterConfig) GetExpireSeconds() int {
	return c.ExpireSeconds
}

func (c *RateLimiterConfig) GetSkipSuccessful() bool {
	return c.SkipSuccessful
}
