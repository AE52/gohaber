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
	Register(req domain.RegisterUserRequest) (*domain.User, error)
	Authenticate(username, password string) (*domain.User, error)
	ResetPassword(email string) (string, error)
	ConfirmResetPassword(token, newPassword string) error
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
func (s *AuthService) Register(req domain.RegisterUserRequest) (*domain.User, error) {
	// E-posta adresi zaten kayıtlı mı?
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrDuplicateEntry
	}

	// Kullanıcı adı zaten kayıtlı mı?
	existingUser, err = s.userRepo.GetByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, domain.ErrDuplicateEntry
	}

	// Şifreleri eşleşiyor mu?
	if req.Password != req.ConfirmPassword {
		return nil, &domain.ValidationError{
			Field:   "confirm_password",
			Message: "Şifreler eşleşmiyor",
		}
	}

	// Şifreyi hashle
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return user, nil
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
