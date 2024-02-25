package portal

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(group fiber.Router) error {

	group.Get("/", dashboard)
	group.Get("/overview/vms", vmTable)
	group.Get("/overview/vm/:id", vmInfo)

	return nil
}
