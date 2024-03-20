package portal

import (
	"net/netip"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/NachoxMacho/opnlab/opnsense"
)

func ipInformation(c *fiber.Ctx) error {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	interfaces, err := getObjectsFromCache[opnsense.Interface](redisClient, "interfaces")
	if err != nil {
		return err
	}
	leases, err := getObjectsFromCache[opnsense.DHCPLease](redisClient, "leases")
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

	return c.Render("opnsense/nextip", fiber.Map{"NextIP": nextIPs})
}
