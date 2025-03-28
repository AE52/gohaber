package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/username/haber/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthServiceExtended, genişletilmiş auth servisi
type AuthServiceExtended struct {
	DB *gorm.DB
}

// NewAuthServiceExtended yeni bir AuthServiceExtended oluşturur
func NewAuthServiceExtended(db *gorm.DB) *AuthServiceExtended {
	return &AuthServiceExtended{
		DB: db,
	}
}

// Authenticate kullanıcı adı ve parola ile kimlik doğrulama yapar
func (s *AuthServiceExtended) Authenticate(username, password string) (*domain.User, error) {
	var user domain.User

	// Kullanıcı adı veya e-posta ile sorgulama yapabiliriz
	result := s.DB.Where("(username = ? OR email = ?)", username, username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Kullanıcı adı veya parola hatalı")
		}
		return nil, result.Error
	}

	// Parola doğrulama
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("Kullanıcı adı veya parola hatalı")
	}

	// Son giriş zamanını güncelle
	s.DB.Model(&user).Update("updated_at", time.Now())

	return &user, nil
}

// CreateUser yeni bir kullanıcı oluşturur
func (s *AuthServiceExtended) CreateUser(req domain.RegisterUserRequest) (*domain.User, error) {
	// Zorunlu alanları kontrol et
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("Kullanıcı adı, e-posta ve parola zorunludur")
	}

	// Parolaların eşleştiğini kontrol et
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("Parolalar eşleşmiyor")
	}

	// Kullanıcı adının benzersiz olduğunu kontrol et
	var existingUser domain.User
	result := s.DB.Where("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("Bu kullanıcı adı zaten kullanılıyor")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// E-postanın benzersiz olduğunu kontrol et
	result = s.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("Bu e-posta adresi zaten kullanılıyor")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// Parolayı hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         domain.RoleUser, // Yeni kullanıcılar varsayılan olarak üye
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Kullanıcıyı veritabanına ekle
	result = s.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// GetUserByID ID'ye göre kullanıcı getirir
func (s *AuthServiceExtended) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	result := s.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Kullanıcı bulunamadı")
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername kullanıcı adına göre kullanıcı getirir
func (s *AuthServiceExtended) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	result := s.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Kullanıcı bulunamadı")
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail e-posta adresine göre kullanıcı getirir
func (s *AuthServiceExtended) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Kullanıcı bulunamadı")
		}
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser kullanıcı bilgilerini günceller
func (s *AuthServiceExtended) UpdateUser(user *domain.User) error {
	result := s.DB.Save(user)
	return result.Error
}

// DeleteUser kullanıcıyı siler
func (s *AuthServiceExtended) DeleteUser(id uint) error {
	return s.DB.Delete(&domain.User{}, id).Error
}

// ChangePassword kullanıcı parolasını değiştirir
func (s *AuthServiceExtended) ChangePassword(id uint, currentPassword, newPassword string) error {
	var user domain.User
	result := s.DB.First(&user, id)
	if result.Error != nil {
		return result.Error
	}

	// Mevcut parolayı doğrula
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword))
	if err != nil {
		return errors.New("Mevcut parola yanlış")
	}

	// Yeni parola hash'i
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Parolayı güncelle
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.DB.Save(&user).Error
}

// ResetPassword parola sıfırlama isteği oluşturur
func (s *AuthServiceExtended) ResetPassword(email string) (string, error) {
	var user domain.User
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("Bu e-posta adresine sahip bir kullanıcı bulunamadı")
		}
		return "", result.Error
	}

	// Token oluştur
	token, err := s.CreateToken(user.ID, "reset-password", 24*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ConfirmResetPassword parola sıfırlama işlemini onaylar
func (s *AuthServiceExtended) ConfirmResetPassword(token, newPassword string) error {
	// Tokeni doğrula
	user, err := s.VerifyToken(token, "reset-password")
	if err != nil {
		return err
	}

	// Yeni parola hash'i
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Kullanıcı şifresini güncelle
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()
	err = s.DB.Save(user).Error
	if err != nil {
		return err
	}

	// Token'ı sil (Bu kısım model yapımıza göre ayarlanmalı)
	// s.DB.Where("token = ? AND type = ?", token, "reset-password").Delete(&domain.Token{})

	return nil
}

// CreateToken token oluşturur
func (s *AuthServiceExtended) CreateToken(userID uint, tokenType string, expiration time.Duration) (string, error) {
	// Random token oluştur
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	// TODO: Token'ı veritabanına kaydetme işlemi domain.Token modeli tanımlandıktan sonra eklenecek

	return token, nil
}

// VerifyToken token'ı doğrular ve kullanıcıyı getirir
func (s *AuthServiceExtended) VerifyToken(token, tokenType string) (*domain.User, error) {
	// TODO: Token'ı veritabanında doğrulama işlemi domain.Token modeli tanımlandıktan sonra eklenecek
	// Şu an için basitleştirilmiş bir implementasyon:

	var user domain.User
	s.DB.First(&user, 1) // İlk kullanıcıyı getir, gerçekte token ile ilişkili kullanıcı getirilmeli

	return &user, nil
}
