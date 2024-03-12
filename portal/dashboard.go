package portal

import (
	"fmt"
	"log"
	"math/rand"
	"net/netip"
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
		return strings.Compare(a.Config.Name, b.Config.Name)
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
	outputVMs := make([]VMOutputData,0,len(vms))
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
		slices.SortStableFunc(leases, func(a, b opnsense.DHCPLease) int {
			return a.Address.Compare(b.Address)
		})
	case "hostname":
		slices.SortStableFunc(leases, func(a, b opnsense.DHCPLease) int {
			return strings.Compare(strings.ToLower(a.Hostname), strings.ToLower(b.Hostname))
		})
	case "macaddress":
		slices.SortStableFunc(leases, func(a, b opnsense.DHCPLease) int {
			return strings.Compare(a.MAC.String(), b.MAC.String())
		})
	}

	if descending {
		slices.Reverse(leases)
	}

	interfaces, err := opnsense.GetInterfaces()
	if err != nil {
		return err
	}

	usedIPs := make([]netip.Addr, len(leases))
	for i, lease := range leases {
		usedIPs[i] = lease.Address
	}

	nextIPs := make([]string, 0, len(interfaces))
	for _, i := range interfaces {
		if i.Status == "down" {
			continue
		}
		if i.Status == "no carrier" {
			continue
		}
		if i.Device == "igb0" {
			continue
		}
		if strings.HasPrefix(i.Device, "lo") {
			continue
		}

		subnet, err := i.SubnetIPv4()
		if err != nil {
			return err
		}
		nextIP := getNewIP(subnet, usedIPs, true)
		nextIPs = append(nextIPs, i.Description+": "+nextIP.String())
	}

	return c.Render("opnsense/dhcp-table", fiber.Map{"Leases": leases, "NextIP": nextIPs, "SortType": sortType, "Descending": descending}, "layouts/main")
}

func getNewIP(subnet netip.Prefix, usedIPs []netip.Addr, randomize bool) netip.Addr {

	if subnet.IsSingleIP() {
		if slices.Contains(usedIPs, subnet.Addr()) {
			return subnet.Addr()
		}
		return netip.Addr{}
	}

	unusedIPs := []netip.Addr{}
	for ip := subnet.Addr(); subnet.Contains(ip); ip = ip.Next() {
		if slices.Contains(usedIPs, ip) {
			continue
		}

		if !randomize {
			return ip
		}

		unusedIPs = append(unusedIPs, ip)
	}

	if len(unusedIPs) == 0 {
		return netip.Addr{}
	}

	log.Println("UnusedIPs:" + fmt.Sprintf("%d", len(unusedIPs)))

	randIndex := rand.Intn(len(unusedIPs))

	return unusedIPs[randIndex]
}

func interfacesOverview(c *fiber.Ctx) error {

	interfaces, err := opnsense.GetInterfaces()
	if err != nil {
		return err
	}

	return c.Render("opnsense/interface", fiber.Map{"Interfaces": interfaces}, "layouts/main")
}
