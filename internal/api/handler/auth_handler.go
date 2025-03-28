package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/api/middleware"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/service"
)

// RefreshTokenRequest token yenileme isteği
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// ResetPasswordRequest şifre sıfırlama isteği
type ResetPasswordRequest struct {
	Email string `json:"email"`
}

// ConfirmResetPasswordRequest şifre sıfırlama onay isteği
type ConfirmResetPasswordRequest struct {
	Token           string `json:"token"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

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
func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	auth := app.Group("/api/auth")

	// Public routes
	auth.Post("/register", middleware.ValidateRequest(&domain.RegisterUserRequest{}), h.Register)
	auth.Post("/login", middleware.ValidateRequest(&domain.LoginRequest{}), h.Login)
	auth.Post("/refresh", h.RefreshToken)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Post("/reset-password", h.ResetPassword)

	// Protected routes
	auth.Use(middleware.NewAuthMiddleware(nil).Authenticate)
	auth.Get("/me", h.GetCurrentUser)
	auth.Post("/logout", h.Logout)
	auth.Put("/change-password", middleware.ValidateRequest(&domain.UpdatePasswordRequest{}), h.ChangePassword)
}

// Register kullanıcı kaydını sağlar
// @Summary Kullanıcı kaydı
// @Description Yeni kullanıcı hesabı oluşturur
// @Tags Kimlik Doğrulama
// @Accept json
// @Produce json
// @Param register body domain.RegisterUserRequest true "Kullanıcı kayıt bilgileri"
// @Success 201 {object} domain.AuthResponse "Başarılı kayıt"
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek formatı veya şifreler eşleşmiyor"
// @Failure 409 {object} domain.ErrorResponse "Kullanıcı adı veya e-posta zaten kullanımda"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Validasyon middleware ile doğrulanmış veriyi al
	reqData := middleware.GetValidated(c).(*domain.RegisterUserRequest)

	// Kullanıcı oluştur
	user, token, err := h.authService.Register(reqData)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Kullanıcı başarıyla oluşturuldu",
		"data": fiber.Map{
			"user":         user,
			"access_token": token,
		},
	})
}

// Login kullanıcı girişini sağlar
// @Summary Kullanıcı girişi
// @Description Kullanıcı adı ve şifre ile giriş yaparak token alır
// @Tags Kimlik Doğrulama
// @Accept json
// @Produce json
// @Param login body domain.LoginRequest true "Kullanıcı giriş bilgileri"
// @Success 200 {object} domain.AuthResponse "Başarılı giriş"
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek formatı"
// @Failure 401 {object} domain.ErrorResponse "Geçersiz kullanıcı adı veya şifre"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Validasyon middleware ile doğrulanmış veriyi al
	reqData := middleware.GetValidated(c).(*domain.LoginRequest)

	// Giriş yap
	user, token, err := h.authService.Login(reqData.Username, reqData.Password, reqData.Remember)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":         user,
			"access_token": token,
		},
	})
}

// GetCurrentUser giriş yapmış kullanıcının bilgilerini döndürür
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	// Context'ten kullanıcı bilgisini al - user_id anahtarını kullanıyoruz
	userID := c.Locals("user_id").(uint)

	// Kullanıcıyı servis üzerinden al
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		return err
	}

	// UserResponse'a dönüştür
	userResponse := &domain.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		FullName:     user.FullName,
		Role:         user.Role,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": userResponse,
		},
	})
}

// RefreshToken erişim token'ını yeniler
// @Summary Token yenileme
// @Description Refresh token kullanarak yeni bir access token alır
// @Tags Kimlik Doğrulama
// @Accept json
// @Produce json
// @Param refresh body RefreshTokenRequest true "Refresh token bilgileri"
// @Success 200 {object} domain.TokenResponse "Yeni token bilgileri"
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek formatı"
// @Failure 401 {object} domain.ErrorResponse "Geçersiz refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// TODO: Refresh token implementasyonu
	return nil
}

// Logout kullanıcının çıkış yapmasını sağlar
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// TODO: Logout implementasyonu
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Başarıyla çıkış yapıldı",
	})
}

// ForgotPassword şifre sıfırlama bağlantısı gönderir
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	// TODO: Şifre sıfırlama implementasyonu
	var request struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.BodyParser(&request); err != nil {
		return domain.NewValidationError([]middleware.ValidationError{
			{
				Field:   "email",
				Message: "Geçersiz email formatı",
			},
		})
	}

	// Validasyon yap
	if validationErrors, err := middleware.ValidateStruct(request); err != nil {
		return domain.NewValidationError(validationErrors)
	}

	// Email gönderme işlemi yapılır
	err := h.authService.ForgotPassword(request.Email)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Şifre sıfırlama bağlantısı email adresinize gönderildi",
	})
}

// ResetPassword şifre sıfırlama isteği oluşturur
// @Summary Şifre sıfırlama isteği
// @Description E-posta adresi ile şifre sıfırlama isteği oluşturur
// @Tags Kimlik Doğrulama
// @Accept json
// @Produce json
// @Param reset body ResetPasswordRequest true "E-posta bilgisi"
// @Success 200 {object} domain.MessageResponse "Şifre sıfırlama e-postası gönderildi"
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek formatı"
// @Failure 404 {object} domain.ErrorResponse "E-posta adresi bulunamadı"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	// TODO: Şifre sıfırlama implementasyonu
	return nil
}

// ChangePassword şifre değiştirir
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	// Validasyon middleware ile doğrulanmış veriyi al
	reqData := middleware.GetValidated(c).(*domain.UpdatePasswordRequest)

	// Kullanıcı ID'sini context'ten al
	userID := c.Locals("user_id").(uint)

	// Şifre değiştir
	err := h.authService.ChangePassword(userID, reqData.CurrentPassword, reqData.NewPassword)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Şifreniz başarıyla değiştirildi",
	})
}
