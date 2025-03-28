package service

import (
	"time"

	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/repository"
	"github.com/username/haber/pkg/auth"
)

// IUserService kullanıcı işlemleri için servis arayüzü
type IUserService interface {
	GetUserByID(id uint) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	ListUsers(offset, limit int, filters map[string]interface{}) ([]*domain.User, int64, error)
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
	ChangePassword(id uint, currentPassword, newPassword string) error
	UpdatePassword(id uint, currentPassword, newPassword string) error
	CheckPassword(hashedPassword, password string) bool
	SearchUsers(query string, offset, limit int) ([]*domain.User, int64, error)
}

// UserService kullanıcı servis implementasyonu
type UserService struct {
	userRepo repository.IUserRepository
}

// NewUserService yeni bir UserService oluşturur
func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID kullanıcıyı ID'ye göre getirir
func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

// GetUserByUsername kullanıcıyı kullanıcı adına göre getirir
func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.userRepo.GetByUsername(username)
}

// GetUserByEmail kullanıcıyı e-posta adresine göre getirir
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(email)
}

// ListUsers kullanıcıları listeler
func (s *UserService) ListUsers(offset, limit int, filters map[string]interface{}) ([]*domain.User, int64, error) {
	return s.userRepo.List(offset, limit, filters)
}

// CreateUser yeni bir kullanıcı oluşturur
func (s *UserService) CreateUser(user *domain.User) error {
	// Şifreyi hashle
	hashedPassword, err := auth.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return s.userRepo.Create(user)
}

// UpdateUser kullanıcıyı günceller
func (s *UserService) UpdateUser(user *domain.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

// DeleteUser kullanıcıyı siler
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

// ChangePassword kullanıcı şifresini değiştirir
func (s *UserService) ChangePassword(id uint, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Mevcut şifreyi kontrol et
	if !auth.CheckPassword(user.PasswordHash, currentPassword) {
		return &domain.ValidationError{
			Field:   "current_password",
			Message: "Mevcut şifre hatalı",
		}
	}

	// Yeni şifreyi hashle
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

// UpdatePassword kullanıcı şifresini değiştirir
func (s *UserService) UpdatePassword(id uint, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Mevcut şifreyi kontrol et
	if !auth.CheckPassword(user.PasswordHash, currentPassword) {
		return &domain.ValidationError{
			Field:   "current_password",
			Message: "Mevcut şifre hatalı",
		}
	}

	// Yeni şifreyi hashle
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

// CheckPassword şifre doğrulaması yapar
func (s *UserService) CheckPassword(hashedPassword, password string) bool {
	return auth.CheckPassword(hashedPassword, password)
}

// SearchUsers kullanıcıları arar
func (s *UserService) SearchUsers(query string, offset, limit int) ([]*domain.User, int64, error) {
	return s.userRepo.Search(query, offset, limit)
}
