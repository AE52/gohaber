package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/pkg/auth"
)

// AuthMiddleware kimlik doğrulama işlemlerini yönetir
type AuthMiddleware struct {
	jwtAuth *auth.JWTAuth
}

// NewAuthMiddleware yeni bir AuthMiddleware oluşturur
func NewAuthMiddleware(jwtAuth *auth.JWTAuth) *AuthMiddleware {
	return &AuthMiddleware{
		jwtAuth: jwtAuth,
	}
}

// Authenticate, Protected metodu için bir alias - kimlik doğrulama gerektiren rotalar için
func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	return m.Protected()(c)
}

// Protected kimlik doğrulama gerektiren istekleri korur
func (m *AuthMiddleware) Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Token'ı al
		token := m.extractToken(c)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Kimlik doğrulama token'ı gereklidir",
			})
		}

		// Token'ı doğrula
		user, err := m.jwtAuth.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Geçersiz veya süresi dolmuş token",
			})
		}

		// Kullanıcıyı context'e ekle
		c.Locals("user", user)
		c.Locals("user_id", user.ID)
		c.Locals("user_role", user.Role)

		return c.Next()
	}
}

// RequireRole belirli rol gerektiren istekleri korur
func (m *AuthMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Önce kullanıcının giriş yapıp yapmadığını kontrol et
		err := m.Protected()(c)
		if err != nil {
			return err
		}

		// Kullanıcının rolünü al
		userRole := c.Locals("user_role").(string)

		// Rolü kontrol et
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		// Yetkisiz erişim
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bu işlem için yetkiniz bulunmuyor",
		})
	}
}

// extractToken istek header'ından token'ı çıkarır
func (m *AuthMiddleware) extractToken(c *fiber.Ctx) string {
	// Authorization header'ını al
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Bearer token'ı ayır
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
