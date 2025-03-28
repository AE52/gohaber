package domain

import (
	"errors"
	"fmt"
	"net/http"
)

// Genel domain hataları
var (
	ErrNotFound           = errors.New("kayıt bulunamadı")
	ErrInvalidInput       = errors.New("geçersiz girdi")
	ErrDuplicateEntry     = errors.New("kayıt zaten mevcut")
	ErrUnauthorized       = errors.New("yetkisiz erişim")
	ErrForbidden          = errors.New("erişim reddedildi")
	ErrValidationFailed   = errors.New("doğrulama hatası")
	ErrInternalError      = errors.New("iç sunucu hatası")
	ErrBadRequest         = errors.New("geçersiz istek")
	ErrTimeout            = errors.New("istek zaman aşımına uğradı")
	ErrServiceUnavailable = errors.New("servis kullanılamıyor")
	ErrTooManyRequests    = errors.New("çok fazla istek")
	ErrDatabaseError      = errors.New("veritabanı hatası")
)

// ErrorCode özel hata kodları
type ErrorCode string

// Hata kodları
const (
	ErrorCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrorCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrorCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrorCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrorCodeDuplicate       ErrorCode = "DUPLICATE_ENTRY"
	ErrorCodeInternal        ErrorCode = "INTERNAL_ERROR"
	ErrorCodeTimeout         ErrorCode = "TIMEOUT"
	ErrorCodeUnavailable     ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
	ErrorCodeDatabase        ErrorCode = "DATABASE_ERROR"
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

// AppError uygulama genelinde kullanılan hata yapısı
type AppError struct {
	Err        error       // Orijinal hata
	Code       ErrorCode   // Hata kodu
	Message    string      // Kullanıcı dostu mesaj
	StatusCode int         // HTTP durum kodu
	Details    interface{} // Ek detaylar
}

// Error hata mesajını döndürür
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "bilinmeyen hata"
}

// Is hata türünü kontrol eder
func (e *AppError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

// Unwrap altta yatan hatayı döndürür
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError yeni bir AppError oluşturur
func NewAppError(err error, code ErrorCode, message string, statusCode int, details interface{}) *AppError {
	return &AppError{
		Err:        err,
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    details,
	}
}

// NewNotFoundError yeni bir kayıt bulunamadı hatası oluşturur
func NewNotFoundError(resourceType ResourceType, id interface{}) *AppError {
	msg := fmt.Sprintf("%s bulunamadı: %v", resourceType, id)
	return NewAppError(ErrNotFound, ErrorCodeNotFound, msg, http.StatusNotFound, nil)
}

// NewValidationError yeni bir doğrulama hatası oluşturur
func NewValidationError(validationErrors interface{}) *AppError {
	return NewAppError(ErrValidationFailed, ErrorCodeValidation, "Doğrulama hatası", http.StatusBadRequest, validationErrors)
}

// NewAuthError yeni bir kimlik doğrulama hatası oluşturur
func NewAuthError(message string) *AppError {
	return NewAppError(ErrUnauthorized, ErrorCodeUnauthorized, message, http.StatusUnauthorized, nil)
}

// NewForbiddenError yeni bir erişim reddedildi hatası oluşturur
func NewForbiddenError(message string) *AppError {
	return NewAppError(ErrForbidden, ErrorCodeForbidden, message, http.StatusForbidden, nil)
}

// NewDuplicateError yeni bir kayıt zaten mevcut hatası oluşturur
func NewDuplicateError(resourceType ResourceType, field string, value interface{}) *AppError {
	msg := fmt.Sprintf("%s için %s değeri zaten kullanılıyor: %v", resourceType, field, value)
	return NewAppError(ErrDuplicateEntry, ErrorCodeDuplicate, msg, http.StatusConflict, nil)
}

// NewInternalError yeni bir iç sunucu hatası oluşturur
func NewInternalError(err error) *AppError {
	return NewAppError(err, ErrorCodeInternal, "İç sunucu hatası", http.StatusInternalServerError, nil)
}

// NewDatabaseError yeni bir veritabanı hatası oluşturur
func NewDatabaseError(err error) *AppError {
	return NewAppError(err, ErrorCodeDatabase, "Veritabanı işlemi başarısız", http.StatusInternalServerError, nil)
}

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
