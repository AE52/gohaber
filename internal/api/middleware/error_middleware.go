package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/domain"
)

// ErrorHandler uygulama genelinde hata yönetimini standartlaştırır
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Varsayılan hata kodu
		code := fiber.StatusInternalServerError
		message := "İç sunucu hatası"

		// Domain hatalarına göre özel kodlar
		var notFoundErr *domain.NotFoundError
		if errors.As(err, &notFoundErr) {
			code = fiber.StatusNotFound
			message = notFoundErr.Error()
		}

		var validationErr *domain.ValidationError
		if errors.As(err, &validationErr) {
			code = fiber.StatusBadRequest
			message = validationErr.Error()
		}

		var validationErrs domain.ValidationErrors
		if errors.As(err, &validationErrs) {
			code = fiber.StatusBadRequest
			message = validationErrs.Error()
		}

		var authErr *domain.AuthError
		if errors.As(err, &authErr) {
			code = fiber.StatusUnauthorized
			message = authErr.Error()
		}

		// Genel hata türlerini kontrol et
		if errors.Is(err, domain.ErrNotFound) {
			code = fiber.StatusNotFound
			message = "Kayıt bulunamadı"
		} else if errors.Is(err, domain.ErrInvalidInput) {
			code = fiber.StatusBadRequest
			message = "Geçersiz girdi"
		} else if errors.Is(err, domain.ErrDuplicateEntry) {
			code = fiber.StatusConflict
			message = "Kayıt zaten mevcut"
		} else if errors.Is(err, domain.ErrUnauthorized) {
			code = fiber.StatusUnauthorized
			message = "Yetkisiz erişim"
		} else if errors.Is(err, domain.ErrForbidden) {
			code = fiber.StatusForbidden
			message = "Erişim reddedildi"
		} else if errors.Is(err, domain.ErrValidationFailed) {
			code = fiber.StatusBadRequest
			message = "Doğrulama hatası"
		}

		// Fiber'ın kendi hatalarını kontrol et
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			message = fiberErr.Message
		}

		// Content type'ı kontrol et ve uygun yanıt döndür
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(code).JSON(fiber.Map{
			"error": message,
			"code":  code,
		})
	}
}
