package repository

import (
	"errors"

	"github.com/username/haber/internal/domain"
	"gorm.io/gorm"
)

// IUserRepository kullanıcı işlemleri için repository interface
type IUserRepository interface {
	Create(user *domain.User) error
	GetByID(id uint) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	List(offset, limit int, filters map[string]interface{}) ([]*domain.User, int64, error)
	FindByResetToken(token string) (*domain.User, error)
	Search(query string, offset, limit int) ([]*domain.User, int64, error)
}

// UserRepository UserRepository'nin GORM implementasyonu
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository yeni bir IUserRepository instance'ı oluşturur
func NewUserRepository(db *Database) IUserRepository {
	return &UserRepository{db: db.DB}
}

// Create yeni bir kullanıcı oluşturur
func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetByID ID'ye göre kullanıcı getirir
func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: domain.ResourceUser,
				ID:           id,
			}
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername kullanıcı adına göre kullanıcı getirir
func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: domain.ResourceUser,
				ID:           username,
			}
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail email'e göre kullanıcı getirir
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: domain.ResourceUser,
				ID:           email,
			}
		}
		return nil, err
	}
	return &user, nil
}

// Update kullanıcı bilgilerini günceller
func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete kullanıcıyı siler (soft delete)
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

// List kullanıcıları listeler
func (r *UserRepository) List(offset, limit int, filters map[string]interface{}) ([]*domain.User, int64, error) {
	var users []*domain.User
	var count int64

	query := r.db.Model(&domain.User{})

	// Filtreleri uygula
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// FindByResetToken reset token'a göre kullanıcı bulur
func (r *UserRepository) FindByResetToken(token string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("reset_token = ? AND reset_token_expires_at > NOW()", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{
				ResourceType: domain.ResourceUser,
				ID:           token,
			}
		}
		return nil, err
	}
	return &user, nil
}

// Search kullanıcıları arar
func (r *UserRepository) Search(query string, offset, limit int) ([]*domain.User, int64, error) {
	var users []*domain.User
	var count int64

	err := r.db.Model(&domain.User{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("username LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").
		Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}
