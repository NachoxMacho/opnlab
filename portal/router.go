package portal

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(group fiber.Router) error {

	group.Get("/", dashboard)

	return nil
}
