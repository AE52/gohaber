package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/service"
)

// AuthHandler kimlik doğrulama işleyicileri
type AuthHandler struct {
	authService service.IAuthService
}

// NewAuthHandler yeni bir AuthHandler oluşturur
func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRoutes rotaları kayıt eder
func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	// Kimlik doğrulama rotaları
	auth := router.Group("/auth")
	auth.Post("/login", h.Login)
	auth.Post("/register", h.Register)
	auth.Post("/refresh", h.RefreshToken)
	auth.Post("/reset-password", h.ResetPassword)
	auth.Post("/reset-password/confirm", h.ConfirmResetPassword)
}

// Login kullanıcı girişini sağlar
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	user, err := h.authService.Authenticate(req.Username, req.Password)
	if err != nil {
		return err
	}

	// TODO: Token oluşturma mekanizması eklendikten sonra accessToken ve refreshToken değerleri ayarlanacak
	accessToken := "sample-access-token"
	refreshToken := "sample-refresh-token"

	return c.JSON(domain.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

// Register kullanıcı kaydını sağlar
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	if req.Password != req.ConfirmPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Şifreler eşleşmiyor")
	}

	user, err := h.authService.Register(req)
	if err != nil {
		return err
	}

	// TODO: Token oluşturma mekanizması eklendikten sonra accessToken ve refreshToken değerleri ayarlanacak
	accessToken := "sample-access-token"
	refreshToken := "sample-refresh-token"

	return c.Status(fiber.StatusCreated).JSON(domain.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

// RefreshToken erişim token'ını yeniler
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	// TODO: Token yenileme sistemi eklenince bu kısım güncellenecek
	accessToken := "new-sample-access-token"
	refreshToken := "new-sample-refresh-token"

	return c.JSON(domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

// ResetPassword şifre sıfırlama isteği oluşturur
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	type ResetPasswordRequest struct {
		Email string `json:"email"`
	}

	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	token, err := h.authService.ResetPassword(req.Email)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Şifre sıfırlama bağlantısı e-posta adresinize gönderildi",
		"token":   token, // Gerçek uygulamada bu token client'a gönderilmemeli, sadece e-posta ile iletilmeli
	})
}

// ConfirmResetPassword şifre sıfırlama işlemini tamamlar
func (h *AuthHandler) ConfirmResetPassword(c *fiber.Ctx) error {
	type ConfirmResetPasswordRequest struct {
		Token           string `json:"token"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	var req ConfirmResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	if req.NewPassword != req.ConfirmPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Şifreler eşleşmiyor")
	}

	err := h.authService.ConfirmResetPassword(req.Token, req.NewPassword)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Şifreniz başarıyla güncellendi",
	})
}
