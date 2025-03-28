package middleware

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/domain"
)

// Validasyon erroru için ileti formatlarını içerir
var validationErrorMessages = map[string]string{
	"required":  "{0} alanı zorunludur",
	"email":     "{0} geçerli bir e-posta adresi olmalıdır",
	"min":       "{0} en az {1} karakter uzunluğunda olmalıdır",
	"max":       "{0} en fazla {1} karakter uzunluğunda olmalıdır",
	"len":       "{0} tam olarak {1} karakter uzunluğunda olmalıdır",
	"eq":        "{0} {1} değerine eşit olmalıdır",
	"ne":        "{0} {1} değerine eşit olmamalıdır",
	"gt":        "{0} {1} değerinden büyük olmalıdır",
	"lt":        "{0} {1} değerinden küçük olmalıdır",
	"gte":       "{0} {1} değerine eşit veya büyük olmalıdır",
	"lte":       "{0} {1} değerine eşit veya küçük olmalıdır",
	"alpha":     "{0} sadece harflerden oluşmalıdır",
	"alphanum":  "{0} sadece harf ve rakamlardan oluşmalıdır",
	"numeric":   "{0} sadece rakamlardan oluşmalıdır",
	"lowercase": "{0} sadece küçük harflerden oluşmalıdır",
	"uppercase": "{0} sadece büyük harflerden oluşmalıdır",
}

// ValidationError validasyon hatası yapısı
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// getValidator validasyon instance'ını döndürür
func getValidator() *validator.Validate {
	validate := validator.New()

	// Struct alan etiketlerini özelleştir
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}

// ValidateStruct verilen struct'ı validate eder
func ValidateStruct(payload interface{}) ([]ValidationError, error) {
	validate := getValidator()

	err := validate.Struct(payload)
	if err == nil {
		return nil, nil
	}

	var validationErrors []ValidationError

	// Tüm validasyon hatalarını döngü ile işle
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		for _, err := range validationErr {
			var element ValidationError
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()

			// Kullanıcı dostu hata mesajı oluştur
			if msg, ok := validationErrorMessages[err.Tag()]; ok {
				msg = strings.Replace(msg, "{0}", err.Field(), -1)
				msg = strings.Replace(msg, "{1}", err.Param(), -1)
				element.Message = msg
			} else {
				element.Message = err.Error()
			}

			validationErrors = append(validationErrors, element)
		}
	}

	return validationErrors, err
}

// ValidateRequest istek validasyonu için middleware
func ValidateRequest(structPtr interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Yeni bir örnek oluştur (struct pointer için)
		payload := reflect.New(reflect.TypeOf(structPtr).Elem()).Interface()

		// JSON'ı parse et
		if err := c.BodyParser(payload); err != nil {
			return domain.NewValidationError([]ValidationError{
				{
					Field:   "body",
					Message: "Geçersiz JSON formatı: " + err.Error(),
				},
			})
		}

		// Validasyon yap
		if validationErrors, err := ValidateStruct(payload); err != nil {
			return domain.NewValidationError(validationErrors)
		}

		// Validate edilmiş modeli locale aktar
		c.Locals("validated", payload)

		return c.Next()
	}
}

// GetValidated validasyon yapılmış veriyi context'ten alır
func GetValidated(c *fiber.Ctx) interface{} {
	return c.Locals("validated")
}

// CustomValidator özel validasyon kuralı ekler
func CustomValidator(tag string, fn validator.Func, errMsg string) {
	validate := getValidator()
	validate.RegisterValidation(tag, fn)
	validationErrorMessages[tag] = errMsg
}
