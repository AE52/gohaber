package config

// IConfig tüm konfigürasyon arayüzü
type IConfig interface {
	GetServer() IServerConfig
	GetDatabase() IDatabaseConfig
	GetRedis() IRedisConfig
	GetJWT() IJWTConfig
	GetMinIO() IMinIOConfig
	GetRateLimiter() IRateLimiterConfig
}

// IServerConfig sunucu ayarları arayüzü
type IServerConfig interface {
	GetPort() string
	GetTemplateDir() string
	GetStaticDir() string
	GetUploadsDir() string
	GetMaxUploadMB() int
	GetEnvironment() string
	GetAllowOrigins() string
}

// IDatabaseConfig veritabanı ayarları arayüzü
type IDatabaseConfig interface {
	GetHost() string
	GetPort() string
	GetUser() string
	GetPassword() string
	GetDBName() string
	GetSSLMode() string
	GetDSN() string
}

// IRedisConfig Redis önbellek ayarları arayüzü
type IRedisConfig interface {
	GetHost() string
	GetPort() string
	GetPassword() string
	GetDB() int
	GetRedisAddr() string
}

// IJWTConfig JWT ayarları arayüzü
type IJWTConfig interface {
	GetSecret() string
	GetAccessTokenExp() int
	GetRefreshTokenExp() int
}

// IMinIOConfig MinIO nesne depolama ayarları arayüzü
type IMinIOConfig interface {
	GetEndpoint() string
	GetAccessKeyID() string
	GetSecretAccessKey() string
	GetUseSSL() bool
	GetBucketName() string
	GetLocation() string
}

// IRateLimiterConfig Rate Limiter ayarları arayüzü
type IRateLimiterConfig interface {
	GetEnabled() bool
	GetMaxRequests() int
	GetExpireSeconds() int
	GetSkipSuccessful() bool
}
