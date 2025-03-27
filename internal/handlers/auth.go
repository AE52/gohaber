package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/models"
	"github.com/username/haber/pkg/auth"
)

// AuthHandler kimlik doğrulama işleyicileri
type AuthHandler struct {
	authService models.AuthService
	jwtAuth     *auth.JWTAuth
}

// NewAuthHandler yeni bir AuthHandler oluşturur
func NewAuthHandler(authService models.AuthService, jwtAuth *auth.JWTAuth) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtAuth:     jwtAuth,
	}
}

// RegisterRoutes rotaları kayıt eder
func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	// Giriş ve kayıt sayfaları
	app.Get("/giris", h.LoginPage)
	app.Post("/giris", h.Login)
	app.Get("/kayit", h.RegisterPage)
	app.Post("/kayit", h.Register)
	app.Get("/cikis", h.Logout)

	// Şifre sıfırlama
	app.Get("/sifre-sifirla", h.ResetPasswordPage)
	app.Post("/sifre-sifirla", h.ResetPassword)
	app.Get("/sifre-sifirla/:token", h.ResetPasswordConfirmPage)
	app.Post("/sifre-sifirla/:token", h.ResetPasswordConfirm)
}

// LoginPage giriş sayfasını gösterir
func (h *AuthHandler) LoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Title": "Giriş Yap",
		"Error": c.Query("error"),
	})
}

// Login kullanıcı girişini işler
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Form verilerini al
	username := c.FormValue("username")
	password := c.FormValue("password")
	remember := c.FormValue("remember") == "1"

	// Kimlik doğrulama
	user, err := h.authService.Authenticate(username, password)
	if err != nil {
		return c.Redirect("/giris?error=Kullanıcı adı veya parola hatalı")
	}

	// JWT token oluştur
	accessToken, refreshToken, err := h.jwtAuth.GenerateTokens(user)
	if err != nil {
		return c.Redirect("/giris?error=Oturum oluşturulurken bir hata oluştu")
	}

	// Cookie ayarla
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    accessToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
	}

	// Beni hatırla seçeneği işaretlenmişse süreyi uzat
	if remember {
		cookie.Expires = time.Now().Add(24 * 7 * time.Hour) // 1 hafta
	}

	c.Cookie(&cookie)

	// Refresh token'ı sakla
	refreshCookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		Expires:  time.Now().Add(30 * 24 * time.Hour), // 30 gün
	}
	c.Cookie(&refreshCookie)

	// Kullanıcı rolüne göre yönlendirme
	if user.Role == models.RoleAdmin || user.Role == models.RoleEditor {
		return c.Redirect("/admin/panel")
	}

	// Normal kullanıcıları ana sayfaya yönlendir
	return c.Redirect("/")
}

// RegisterPage kayıt sayfasını gösterir
func (h *AuthHandler) RegisterPage(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"Title": "Üye Ol",
		"Error": c.Query("error"),
	})
}

// Register yeni kullanıcı kaydını işler
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Form verilerini al
	req := models.RegisterUserRequest{
		Username:        c.FormValue("username"),
		Email:           c.FormValue("email"),
		Password:        c.FormValue("password"),
		ConfirmPassword: c.FormValue("confirm_password"),
		FullName:        c.FormValue("full_name"),
	}

	// Kullanıcıyı oluştur
	user, err := h.authService.CreateUser(req)
	if err != nil {
		return c.Redirect("/kayit?error=" + err.Error())
	}

	// Oluşturulan kullanıcı ile giriş yap
	accessToken, refreshToken, err := h.jwtAuth.GenerateTokens(user)
	if err != nil {
		return c.Redirect("/giris?error=Kayıt başarılı, ancak otomatik giriş yapılamadı")
	}

	// Cookie ayarla
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    accessToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
	}
	c.Cookie(&cookie)

	// Refresh token'ı sakla
	refreshCookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		Expires:  time.Now().Add(30 * 24 * time.Hour), // 30 gün
	}
	c.Cookie(&refreshCookie)

	// Ana sayfaya yönlendir
	return c.Redirect("/")
}

// Logout kullanıcı çıkışını işler
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Token cookie'sini temizle
	c.ClearCookie("token", "refresh_token")

	// Ana sayfaya yönlendir
	return c.Redirect("/")
}

// ResetPasswordPage şifre sıfırlama sayfasını gösterir
func (h *AuthHandler) ResetPasswordPage(c *fiber.Ctx) error {
	return c.Render("auth/reset_password", fiber.Map{
		"Title":   "Şifre Sıfırla",
		"Error":   c.Query("error"),
		"Success": c.Query("success"),
	})
}

// ResetPassword şifre sıfırlama isteğini işler
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return c.Redirect("/sifre-sifirla?error=E-posta adresi gereklidir")
	}

	// Şifre sıfırlama tokeni oluştur
	_, err := h.authService.ResetPassword(email)
	if err != nil {
		return c.Redirect("/sifre-sifirla?error=" + err.Error())
	}

	// Gerçek bir uygulamada, e-posta gönderme işlemi burada yapılır
	// ...

	return c.Redirect("/sifre-sifirla?success=Şifre sıfırlama bağlantısı e-posta adresinize gönderildi")
}

// ResetPasswordConfirmPage şifre sıfırlama onay sayfasını gösterir
func (h *AuthHandler) ResetPasswordConfirmPage(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Redirect("/sifre-sifirla?error=Geçersiz token")
	}

	return c.Render("auth/reset_password_confirm", fiber.Map{
		"Title": "Yeni Şifre Belirle",
		"Token": token,
		"Error": c.Query("error"),
	})
}

// ResetPasswordConfirm şifre sıfırlama onayını işler
func (h *AuthHandler) ResetPasswordConfirm(c *fiber.Ctx) error {
	token := c.Params("token")
	newPassword := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	if token == "" {
		return c.Redirect("/sifre-sifirla?error=Geçersiz token")
	}

	if newPassword == "" {
		return c.Redirect("/sifre-sifirla/" + token + "?error=Şifre boş olamaz")
	}

	if newPassword != confirmPassword {
		return c.Redirect("/sifre-sifirla/" + token + "?error=Şifreler eşleşmiyor")
	}

	// Şifre sıfırlama
	err := h.authService.ConfirmResetPassword(token, newPassword)
	if err != nil {
		return c.Redirect("/sifre-sifirla/" + token + "?error=" + err.Error())
	}

	return c.Redirect("/giris?success=Şifreniz başarıyla değiştirildi, giriş yapabilirsiniz")
}
