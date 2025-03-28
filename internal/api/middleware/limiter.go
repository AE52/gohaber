package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/username/haber/internal/config"
)

// LimiterConfig rate limiter ayarları
type LimiterConfig struct {
	Max            int
	ExpireSeconds  int
	SkipSuccessful bool
	KeyGenerator   func(*fiber.Ctx) string
}

// NewRateLimiterMiddleware yeni bir rate limiter middleware oluşturur
func NewRateLimiterMiddleware(cfg config.IConfig) fiber.Handler {
	serverConfig := cfg.GetServer()
	expSeconds := 60   // Varsayılan 60 saniye
	maxRequests := 100 // Varsayılan 100 istek/dakika

	// Üretim ortamında daha sıkı limitler uygula
	if serverConfig.GetEnvironment() == "production" {
		maxRequests = 60 // Üretimde 60 istek/dakika
	}

	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: time.Duration(expSeconds) * time.Second,
		// IP tabanlı rate limit
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		// Başarılı API yanıtları için rate limit uygulanmasın
		SkipSuccessfulRequests: false,
		// Rate limitini aşanları engelle
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Çok fazla istek gönderdiniz, lütfen daha sonra tekrar deneyin",
				"code":  "TOO_MANY_REQUESTS",
			})
		},
		// İsteğe bağlı storage (varsayılan memory)
		// Storage: myCustomStorage{},
	})
}

// NewTokenBucketLimiter daha gelişmiş token bucket rate limiter
func NewTokenBucketLimiter(cfg config.IConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        30,               // 30 token
		Expiration: 30 * time.Second, // 30 saniyede bir yenilenir
		KeyGenerator: func(c *fiber.Ctx) string {
			// Auth token varsa ona göre rate limit uygula
			token := c.Get("Authorization")
			if token != "" {
				return "token:" + token
			}
			// Yoksa IP'ye göre limit uygula
			return "ip:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit aşıldı",
				"code":  "RATE_LIMITED",
			})
		},
	})
}

// NewPathBasedRateLimiter endpoint bazlı farklı rate limitleri uygular
func NewPathBasedRateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Auth endpointleri için daha sıkı limit
		if path == "/api/auth/login" || path == "/api/auth/register" {
			limiterHandler := limiter.New(limiter.Config{
				Max:        5,               // 5 istek
				Expiration: 1 * time.Minute, // 1 dakikada
				KeyGenerator: func(c *fiber.Ctx) string {
					return "auth:" + c.IP()
				},
			})
			return limiterHandler(c)
		}

		// Admin paneli için daha sıkı limit
		if len(c.Path()) >= 6 && c.Path()[:6] == "/admin" {
			limiterHandler := limiter.New(limiter.Config{
				Max:        30,              // 30 istek
				Expiration: 1 * time.Minute, // 1 dakikada
				KeyGenerator: func(c *fiber.Ctx) string {
					return "admin:" + c.IP()
				},
			})
			return limiterHandler(c)
		}

		// API için genel limit
		if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
			limiterHandler := limiter.New(limiter.Config{
				Max:        60,              // 60 istek
				Expiration: 1 * time.Minute, // 1 dakikada
				KeyGenerator: func(c *fiber.Ctx) string {
					return "api:" + c.IP()
				},
			})
			return limiterHandler(c)
		}

		// Diğer tüm istekler için daha yüksek limit
		limiterHandler := limiter.New(limiter.Config{
			Max:        100,             // 100 istek
			Expiration: 1 * time.Minute, // 1 dakikada
		})
		return limiterHandler(c)
	}
}
