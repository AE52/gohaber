package swagger

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// New returns Swagger handler
func New() fiber.Handler {
	return fiberSwagger.WrapHandler
}
