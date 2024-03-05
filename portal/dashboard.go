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

	interfaces, err := opnsense.GetInterfaces()
	if err != nil {
		return err
	}

	usedIPs := []netip.Addr{}
	for _, lease := range leases {
		ip, err := lease.GetIP()
		if err != nil {
			return err
		}
		usedIPs = append(usedIPs, ip)
	}

	nextIPs := []string{}
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
		return subnet.Addr()
	}

	unusedIPs := []netip.Addr{}
	for ip := subnet.Addr(); subnet.Contains(ip); ip = ip.Next() {
		if slices.Index(usedIPs, ip) != -1 {
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
