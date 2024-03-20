package portal

import (
	"github.com/gofiber/fiber/v2"

	"github.com/NachoxMacho/opnlab/opnsense"
)

func dashboard(c *fiber.Ctx) error {
	return c.Render("proxmox/overview", fiber.Map{}, "layouts/main")
}

func interfacesOverview(c *fiber.Ctx) error {

	interfaces, err := opnsense.GetInterfaces()
	if err != nil {
		return err
	}

	return c.Render("opnsense/interface", fiber.Map{"Interfaces": interfaces}, "layouts/main")
}
