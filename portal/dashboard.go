package portal

import (
	"fmt"
	"log"
	"math/rand"
	"net/netip"
	"slices"
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
		Name          string   `json:"name,omitempty"`
		Status        string   `json:"status,omitempty"`
		MaxMemory     string   `json:"max_memory,omitempty"`
		MaxCPUs       string   `json:"max_cp_us,omitempty"`
		CurrentMemory string   `json:"current_memory,omitempty"`
		CurrentCPU    string   `json:"current_cpu,omitempty"`
		MaxDisk       string   `json:"max_disk,omitempty"`
		MACAddress    string   `json:"mac_address,omitempty"`
		IPAddress     string   `json:"ip_address,omitempty"`
		ID            string   `json:"id,omitempty"`
		Tags          []string `json:"tags,omitempty"`
	}
	outputVMs := make([]VMOutputData, 0, len(vms))
	for i, vm := range vms {
		outputVMs[i].ID = fmt.Sprintf("%d", vm.Stats.VMID)

		outputVMs[i].ID = fmt.Sprintf("%d", vm.Stats.VMID)
		outputVMs[i].Name = vm.Config.Name
		outputVMs[i].Status = vm.Stats.Status
		outputVMs[i].MaxCPUs = fmt.Sprintf("%d", vm.Stats.CPUs)
		outputVMs[i].CurrentMemory = HumanFileSize(float64(vm.Stats.Memory))
		outputVMs[i].MaxMemory = HumanFileSize(float64(vm.Stats.MaxMemory))
		outputVMs[i].MaxDisk = HumanFileSize(float64(vm.Stats.MaxDiskBytes))
		outputVMs[i].CurrentCPU = fmt.Sprintf("%f%%", vm.Stats.CPU*100)
		outputVMs[i].MACAddress = vm.Config.MACAddress()
		outputVMs[i].IPAddress = ""
		outputVMs[i].Tags = vm.Config.TagList()

		for _, lease := range leases {
			if strings.EqualFold(lease.MAC.String(), vm.Config.MACAddress()) {
				outputVMs[i].IPAddress = lease.Address.String()
			}
		}
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
