package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/repository"
	"github.com/username/haber/pkg/auth"
)

// IAuthService auth işlemleri için service interface
type IAuthService interface {
	Register(req *domain.RegisterUserRequest) (*domain.User, string, error)
	Authenticate(username, password string) (*domain.User, error)
	Login(username, password string, remember bool) (*domain.UserResponse, string, error)
	ResetPassword(email string) (string, error)
	ConfirmResetPassword(token, newPassword string) error
	ForgotPassword(email string) error
	ChangePassword(userID uint, currentPassword, newPassword string) error
	GetUserByID(id uint) (*domain.User, error)
}

// AuthService auth servisinin implementasyonu
type AuthService struct {
	userRepo repository.IUserRepository
}

// NewAuthService yeni bir AuthService oluşturur
func NewAuthService(userRepo repository.IUserRepository) IAuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register yeni kullanıcı kaydı yapar
func (s *AuthService) Register(req *domain.RegisterUserRequest) (*domain.User, string, error) {
	// E-posta adresi zaten kayıtlı mı?
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, "", domain.ErrDuplicateEntry
	}

	// Kullanıcı adı zaten kayıtlı mı?
	existingUser, err = s.userRepo.GetByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, "", domain.ErrDuplicateEntry
	}

	// Şifreleri eşleşiyor mu?
	if req.Password != req.ConfirmPassword {
		return nil, "", &domain.ValidationError{
			Field:   "confirm_password",
			Message: "Şifreler eşleşmiyor",
		}
	}

	// Şifreyi hashle
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	// Yeni kullanıcı oluştur
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		Role:         domain.RoleUser, // Varsayılan olarak normal kullanıcı
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Kullanıcıyı kaydet
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, "", err
	}

	// Burada gerçek uygulamada token üretilecek
	token := "sample-access-token-for-" + user.Username

	return user, token, nil
}

// Login kullanıcıyı giriş yapar ve token döndürür
func (s *AuthService) Login(username, password string, remember bool) (*domain.UserResponse, string, error) {
	// Önce kullanıcıyı doğrula
	user, err := s.Authenticate(username, password)
	if err != nil {
		return nil, "", err
	}

	// Kullanıcı yanıtı oluştur
	userResponse := &domain.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		FullName:     user.FullName,
		Role:         user.Role,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
	}

	// Burada gerçek uygulamada token üretilecek
	// Hatırla seçeneği varsa daha uzun süreli token verilebilir
	var token string
	if remember {
		token = "sample-long-term-token-for-" + user.Username
	} else {
		token = "sample-access-token-for-" + user.Username
	}

	return userResponse, token, nil
}

// Authenticate kullanıcı girişini doğrular
func (s *AuthService) Authenticate(username, password string) (*domain.User, error) {
	// Kullanıcıyı bul
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, &domain.AuthError{
			Message: "Kullanıcı adı veya şifre hatalı",
		}
	}

	// Şifreyi kontrol et
	if !auth.CheckPassword(user.PasswordHash, password) {
		return nil, &domain.AuthError{
			Message: "Kullanıcı adı veya şifre hatalı",
		}
	}

	return user, nil
}

// ForgotPassword şifre sıfırlama e-postası gönderir
func (s *AuthService) ForgotPassword(email string) error {
	// Kullanıcıyı e-posta adresine göre bul
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// Güvenlik için kullanıcı bulunamasa bile başarılı mesajı dön
		return nil
	}

	// Token oluştur
	token, err := s.ResetPassword(email)
	if err != nil {
		return err
	}

	// TODO: E-posta gönderme işlemi burada yapılacak
	// Örnek: emailService.SendPasswordResetEmail(user.Email, token)

	// Şimdilik sadece logla
	println("Şifre sıfırlama token'ı oluşturuldu: " + token + " kullanıcı: " + user.Email)

	return nil
}

// ChangePassword kullanıcının şifresini değiştirir
func (s *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// Kullanıcıyı bul
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Mevcut şifreyi doğrula
	if !auth.CheckPassword(user.PasswordHash, currentPassword) {
		return &domain.ValidationError{
			Field:   "current_password",
			Message: "Mevcut şifre hatalı",
		}
	}

	// Yeni şifreyi hashle
	passwordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Şifreyi güncelle
	user.PasswordHash = passwordHash
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

// ResetPassword şifre sıfırlama işlemini başlatır
func (s *AuthService) ResetPassword(email string) (string, error) {
	// Kullanıcıyı e-posta adresine göre bul
	_, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	// Sıfırlama token'ı oluştur
	token := generateResetToken()

	// TODO: Bu kısım kullanıcı modelinde reset_token ve expires_at eklendiğinde kullanılacak
	// user.ResetToken = token
	// user.ResetTokenExpiresAt = time.Now().Add(24 * time.Hour)
	// err = s.userRepo.Update(user)
	// if err != nil {
	// 	return "", err
	// }

	return token, nil
}

// ConfirmResetPassword şifre sıfırlama işlemini tamamlar
func (s *AuthService) ConfirmResetPassword(token, newPassword string) error {
	// TODO: Bu kısım kullanıcı modelinde reset_token ve expires_at eklendiğinde kullanılacak
	// user, err := s.userRepo.FindByResetToken(token)
	// if err != nil {
	// 	return err
	// }

	// Şifreyi hashle
	// passwordHash, err := auth.HashPassword(newPassword)
	// if err != nil {
	// 	return err
	// }

	// Kullanıcı şifresini güncelle ve token'ı temizle
	// user.PasswordHash = passwordHash
	// user.ResetToken = ""
	// user.ResetTokenExpiresAt = time.Time{}
	// return s.userRepo.Update(user)

	return nil
}

// GetUserByID kullanıcıyı ID'ye göre getirir
func (s *AuthService) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

// generateResetToken token oluşturur
func generateResetToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
