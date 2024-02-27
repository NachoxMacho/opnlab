package portal

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(group fiber.Router) error {

	group.Get("/", dashboard)
	group.Get("/overview/vms", vmTable)
	group.Get("/overview/vm/:id", vmInfo)
	group.Get("/overview/opnsense", dhcpInfo)

	return nil
}
