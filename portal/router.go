package portal

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(group fiber.Router) error {

	group.Get("/proxmox/overview", dashboard)
	group.Get("/proxmox/overview/vms", vmTable)
	group.Get("/proxmox/vm/:id", vmInfo)
	group.Get("/opnsense/dhcp", dhcpOverview)
	group.Get("/opnsense/dhcp/nextip", ipInformation)
	group.Get("/opnsense/dhcp/leases", dhcpInfo)

	return nil
}
