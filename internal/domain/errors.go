package domain

import (
	"errors"
	"fmt"
)

// Genel domain hataları
var (
	ErrNotFound         = errors.New("kayıt bulunamadı")
	ErrInvalidInput     = errors.New("geçersiz girdi")
	ErrDuplicateEntry   = errors.New("kayıt zaten mevcut")
	ErrUnauthorized     = errors.New("yetkisiz erişim")
	ErrForbidden        = errors.New("erişim reddedildi")
	ErrValidationFailed = errors.New("doğrulama hatası")
	ErrInternalError    = errors.New("iç sunucu hatası")
)

// ResourceType kaynak türleri
type ResourceType string

// Kaynak türleri sabitleri
const (
	ResourceUser     ResourceType = "Kullanıcı"
	ResourceArticle  ResourceType = "Makale"
	ResourceCategory ResourceType = "Kategori"
	ResourceTag      ResourceType = "Etiket"
	ResourceComment  ResourceType = "Yorum"
	ResourceMedia    ResourceType = "Medya"
	ResourceSetting  ResourceType = "Ayar"
	ResourceAdSpace  ResourceType = "Reklam Alanı"
)

// NotFoundError kaynak bulunamadı hatası
type NotFoundError struct {
	ResourceType ResourceType
	ID           interface{}
	Slug         string
}

// Error hata mesajını döndürür
func (e *NotFoundError) Error() string {
	if e.Slug != "" {
		return fmt.Sprintf("%s bulunamadı: %s", e.ResourceType, e.Slug)
	}
	return fmt.Sprintf("%s bulunamadı: %v", e.ResourceType, e.ID)
}

// Is hata türünü kontrol eder
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

// ValidationError doğrulama hatası
type ValidationError struct {
	Field   string
	Message string
}

// Error hata mesajını döndürür
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Is hata türünü kontrol eder
func (e *ValidationError) Is(target error) bool {
	return target == ErrValidationFailed
}

// ValidationErrors doğrulama hataları listesi
type ValidationErrors []ValidationError

// Error hata mesajını döndürür
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "doğrulama hataları"
	}

	return fmt.Sprintf("doğrulama hatası: %s", e[0].Error())
}

// Is hata türünü kontrol eder
func (e ValidationErrors) Is(target error) bool {
	return target == ErrValidationFailed
}

// AuthError kimlik doğrulama hatası
type AuthError struct {
	Message string
}

// Error hata mesajını döndürür
func (e *AuthError) Error() string {
	if e.Message == "" {
		return "kimlik doğrulama hatası"
	}
	return e.Message
}

// Is hata türünü kontrol eder
func (e *AuthError) Is(target error) bool {
	return target == ErrUnauthorized
}
