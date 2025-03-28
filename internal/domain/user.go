package domain

import (
	"time"

	"gorm.io/gorm"
)

// User kullanıcı modelimiz
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	FullName     string         `gorm:"size:100;not null" json:"full_name"`
	Role         string         `gorm:"size:20;not null;default:user" json:"role"` // admin, editor, user
	ProfileImage string         `gorm:"size:255" json:"profile_image,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserRole tanımlı kullanıcı rolleri
const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleUser   = "user"
)

// RegisterUserRequest kullanıcı kaydı için gerekli alanlar
type RegisterUserRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FullName        string `json:"full_name" validate:"required"`
}

// LoginRequest kullanıcı girişi için gerekli alanlar
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Remember bool   `json:"remember"`
}

// UserResponse kullanıcı verilerinin API yanıtı
type UserResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	ProfileImage string    `json:"profile_image,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// UpdateUserRequest kullanıcı bilgilerini güncelleme isteği
type UpdateUserRequest struct {
	Username     string `json:"username,omitempty"`
	Email        string `json:"email,omitempty" validate:"omitempty,email"`
	FullName     string `json:"full_name,omitempty"`
	ProfileImage string `json:"profile_image,omitempty"`
	Role         string `json:"role,omitempty"`
}

// UpdatePasswordRequest şifre güncelleme isteği
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// CreateUserRequest kullanıcı oluşturma isteği (Admin tarafından)
type CreateUserRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FullName        string `json:"full_name" validate:"required"`
	Role            string `json:"role" validate:"required,oneof=admin editor user"`
}
