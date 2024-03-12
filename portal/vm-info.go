package portal

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/NachoxMacho/opnlab/proxmox"
)

func vmInfo(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}

	vm, err := proxmox.GetVMByID(id)
	if err != nil {
		return err
	}

	return c.Render("overview/vm-info", fiber.Map{"Data": vm}, "layouts/main")
}
