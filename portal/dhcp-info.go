package portal

import (
	"fmt"
	"log"
	"math/rand"
	"net/netip"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/NachoxMacho/opnlab/opnsense"
)

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
