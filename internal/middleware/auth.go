package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/models"
	"github.com/username/haber/pkg/auth"
)

// AuthMiddleware kimlik doğrulama işlemleri için middleware
type AuthMiddleware struct {
	jwtAuth *auth.JWTAuth
}

// NewAuthMiddleware yeni bir AuthMiddleware oluşturur
func NewAuthMiddleware(jwtAuth *auth.JWTAuth) *AuthMiddleware {
	return &AuthMiddleware{
		jwtAuth: jwtAuth,
	}
}

// Protected yetkilendirme kontrolü yapar
func (m *AuthMiddleware) Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Token'ı al
		token := extractToken(c)
		if token == "" {
			return c.Redirect("/giris?error=Oturum açmanız gerekiyor")
		}

		// Token'ı doğrula
		claims, err := m.jwtAuth.ValidateToken(token)
		if err != nil {
			// Token süresi dolmuşsa, yenileme tokeni ile yeni bir token almayı dene
			refreshToken := c.Cookies("refresh_token")
			if refreshToken != "" {
				newAccessToken, err := m.jwtAuth.RefreshAccessToken(refreshToken)
				if err == nil {
					// Yeni access token'ı cookie'ye ekle
					c.Cookie(&fiber.Cookie{
						Name:     "token",
						Value:    newAccessToken,
						Path:     "/",
						HTTPOnly: true,
						Secure:   c.Protocol() == "https",
					})

					// Yeni token'ı doğrula
					claims, err = m.jwtAuth.ValidateToken(newAccessToken)
					if err == nil {
						// Kullanıcı bilgisini locale ekle
						c.Locals("user", claims)
						return c.Next()
					}
				}
			}
			return c.Redirect("/giris?error=Oturumunuzun süresi dolmuş")
		}

		// Kullanıcı bilgisini locale ekle
		c.Locals("user", claims)
		return c.Next()
	}
}

// AdminOnly sadece admin rolü kontrolü yapar
func (m *AuthMiddleware) AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Önce kimlik doğrulama kontrolü yap
		err := m.Protected()(c)
		if err != nil {
			return err
		}

		// Kullanıcı bilgisini al
		claims, ok := c.Locals("user").(*auth.JWTCustomClaims)
		if !ok {
			return c.Redirect("/giris?error=Oturum bilgisi alınamadı")
		}

		// Admin rolü kontrolü
		if claims.Role != models.RoleAdmin {
			return c.Status(fiber.StatusForbidden).Render("error", fiber.Map{
				"Title":        "Erişim Reddedildi",
				"ErrorMessage": "Bu sayfaya erişim için gerekli izniniz yok.",
			})
		}

		return c.Next()
	}
}

// AdminOrEditorOnly admin veya editör rolü kontrolü yapar
func (m *AuthMiddleware) AdminOrEditorOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Önce kimlik doğrulama kontrolü yap
		err := m.Protected()(c)
		if err != nil {
			return err
		}

		// Kullanıcı bilgisini al
		claims, ok := c.Locals("user").(*auth.JWTCustomClaims)
		if !ok {
			return c.Redirect("/giris?error=Oturum bilgisi alınamadı")
		}

		// Admin veya editör rolü kontrolü
		if claims.Role != models.RoleAdmin && claims.Role != models.RoleEditor {
			return c.Status(fiber.StatusForbidden).Render("error", fiber.Map{
				"Title":        "Erişim Reddedildi",
				"ErrorMessage": "Bu sayfaya erişim için gerekli izniniz yok.",
			})
		}

		return c.Next()
	}
}

// CurrentUser güncel kullanıcıyı her istekte kullanılabilir hale getirir
func (m *AuthMiddleware) CurrentUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractToken(c)
		if token == "" {
			// Token yoksa anonim kullanıcı
			c.Locals("user", nil)
			return c.Next()
		}

		// Token'ı doğrula
		claims, err := m.jwtAuth.ValidateToken(token)
		if err != nil {
			// Token geçersizse anonim kullanıcı
			c.Locals("user", nil)
			return c.Next()
		}

		// Kullanıcı bilgisini locale ekle
		c.Locals("user", claims)
		return c.Next()
	}
}

// GetUser isteğe bağlı kullanıcı bilgisini getirir
func GetUser(c *fiber.Ctx) (*auth.JWTCustomClaims, error) {
	user := c.Locals("user")
	if user == nil {
		return nil, errors.New("kullanıcı bulunamadı")
	}

	claims, ok := user.(*auth.JWTCustomClaims)
	if !ok {
		return nil, errors.New("geçersiz kullanıcı bilgisi")
	}

	return claims, nil
}

// Token çıkarma yardımcı fonksiyonu
func extractToken(c *fiber.Ctx) string {
	// Önce cookie'den kontrol et
	token := c.Cookies("token")
	if token != "" {
		return token
	}

	// Authorization header'dan kontrol et
	auth := c.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}

	return ""
}
