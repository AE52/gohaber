package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/domain"
)

// ErrorHandler global hata yakalayıcı
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default 500 status code
		statusCode := fiber.StatusInternalServerError
		errorMessage := "İç sunucu hatası"
		errorCode := "INTERNAL_ERROR"
		var errorDetails interface{}

		var appError *domain.AppError
		if errors.As(err, &appError) {
			statusCode = appError.StatusCode
			errorMessage = appError.Message
			errorCode = string(appError.Code)
			errorDetails = appError.Details
		} else {
			// Spesifik hata tiplerine göre yanıt
			switch {
			case errors.Is(err, domain.ErrNotFound):
				statusCode = fiber.StatusNotFound
				errorMessage = err.Error()
				errorCode = "NOT_FOUND"
			case errors.Is(err, domain.ErrValidationFailed):
				statusCode = fiber.StatusBadRequest
				errorMessage = err.Error()
				errorCode = "VALIDATION_ERROR"
			case errors.Is(err, domain.ErrUnauthorized):
				statusCode = fiber.StatusUnauthorized
				errorMessage = err.Error()
				errorCode = "UNAUTHORIZED"
			case errors.Is(err, domain.ErrForbidden):
				statusCode = fiber.StatusForbidden
				errorMessage = err.Error()
				errorCode = "FORBIDDEN"
			case errors.Is(err, domain.ErrDuplicateEntry):
				statusCode = fiber.StatusConflict
				errorMessage = err.Error()
				errorCode = "DUPLICATE_ENTRY"
			case errors.Is(err, domain.ErrTooManyRequests):
				statusCode = fiber.StatusTooManyRequests
				errorMessage = err.Error()
				errorCode = "TOO_MANY_REQUESTS"
			}
		}

		// Validasyon hatalarını özel olarak işle
		var validationErrs ValidationError
		if errors.As(err, &validationErrs) {
			statusCode = fiber.StatusBadRequest
			errorMessage = "Validasyon hatası"
			errorCode = "VALIDATION_ERROR"
			errorDetails = validationErrs
		}

		// Hata yanıtını hazırla
		errResponse := fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    errorCode,
				"message": errorMessage,
			},
		}

		// Detaylar varsa ekle
		if errorDetails != nil {
			errResponse["error"].(fiber.Map)["details"] = errorDetails
		}

		// Development ortamında orijinal hata mesajını da gönder
		if c.App().Config().AppName == "development" && err.Error() != errorMessage {
			errResponse["error"].(fiber.Map)["debug"] = err.Error()
		}

		// Yanıtı döndür
		return c.Status(statusCode).JSON(errResponse)
	}
}
