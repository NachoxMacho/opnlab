package portal

import (
	"fmt"
	"log"
	"math/rand"
	"net/netip"
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/NachoxMacho/opnlab/opnsense"
)

func dhcpOverview(c *fiber.Ctx) error {
	return c.Render("opnsense/dhcp-overview", fiber.Map{}, "layouts/main")
}

func dhcpInfo(c *fiber.Ctx) error {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	leases, err := getObjectsFromCache[opnsense.DHCPLease](redisClient, "leases")
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
			return strings.Compare(a.MAC, b.MAC)
		})
	}

	if descending {
		slices.Reverse(leases)
	}

	return c.Render("opnsense/dhcp-table", fiber.Map{"Leases": leases, "SortType": sortType, "Descending": descending})
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
