package portal

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/NachoxMacho/opnlab/opnsense"
	"github.com/NachoxMacho/opnlab/proxmox"
)

func dashboard(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{}, "layouts/main")
}

func vmTable(c *fiber.Ctx) error {

	vms, err := proxmox.GetVMs()
	if err != nil {
		return err
	}
	slices.SortStableFunc(vms, func(a, b proxmox.VM) int {
		if a.Config.Name > b.Config.Name {
			return 1
		} else if a.Config.Name < b.Config.Name {
			return -1
		}
		return 0
	})

	leases, err := opnsense.GetDHCPLeases()
	if err != nil {
		return err
	}

	type VMOutputData struct {
		Name          string
		Status        string
		MaxMemory     string
		MaxCPUs       string
		CurrentMemory string
		CurrentCPU    string
		MaxDisk       string
		MACAddress    string
		IPAddress     string
		ID            string
		Tags          []string
	}
	outputVMs := []VMOutputData{}
	for _, vm := range vms {
		o := VMOutputData{
			ID:            fmt.Sprintf("%d", vm.Stats.VMID),
			Name:          vm.Config.Name,
			Status:        vm.Stats.Status,
			MaxCPUs:       fmt.Sprintf("%d", vm.Stats.CPUs),
			CurrentMemory: HumanFileSize(float64(vm.Stats.Memory)),
			MaxMemory:     HumanFileSize(float64(vm.Stats.MaxMemory)),
			MaxDisk:       HumanFileSize(float64(vm.Stats.MaxDiskBytes)),
			CurrentCPU:    fmt.Sprintf("%f%%", vm.Stats.CPU*100),
			MACAddress:    vm.Config.MACAddress(),
			IPAddress:     "",
			Tags:          vm.Config.TagList(),
		}

		for _, lease := range leases {
			if strings.EqualFold(lease.MACAddress, vm.Config.MACAddress()) {
				o.IPAddress = lease.Address
			}
		}

		outputVMs = append(outputVMs, o)
	}
	return c.Render("overview/vm-table", fiber.Map{"VMs": outputVMs})
}

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

func dhcpInfo(c *fiber.Ctx) error {

	leases, err := opnsense.GetDHCPLeases()
	if err != nil {
		return err
	}

	sortType := c.Query("sort", "address")
	descending := c.Query("order", "asc") != "asc"

	switch sortType {

	case "address":
		sort.Slice(leases, func(i, j int) bool {
			ip1, err := leases[i].GetIP()
			if err != nil {
				return false
			}

			ip2, err := leases[j].GetIP()
			if err != nil {
				return true
			}
			return ip1.Compare(ip2) < 0 != descending
		})
	case "hostname":
		sort.Slice(leases, func(i, j int) bool {
			return strings.ToLower(leases[i].Hostname) < strings.ToLower(leases[j].Hostname) != descending
		})

	case "macaddress":
		sort.Slice(leases, func(i, j int) bool {
			return strings.ToLower(leases[i].MACAddress) < strings.ToLower(leases[j].MACAddress) != descending
		})
	}

	nextIP, err := opnsense.GetNextAvailableIP()
	if err != nil {
		return err
	}

	return c.Render("opnsense/dhcp-table", fiber.Map{"Leases": leases, "NextIP": nextIP.String(), "SortType": sortType, "Descending": descending}, "layouts/main")
}
