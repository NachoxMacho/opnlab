package opnsense

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func Fetch(client *redis.Client) error {

	leases, err := GetDHCPLeases()
	if err != nil {
		return err
	}

	leasesEncoding, err := json.Marshal(leases)
	if err != nil {
		return err
	}

	client.Set(context.Background(), "leases", leasesEncoding, 5*time.Minute)
	client.Set(context.Background(), "leaseChangeTime", time.Now(), 5*time.Minute)

	interfaces, err := GetInterfaces()
	if err != nil {
		return err
	}

	interfacesEncoding, err := json.Marshal(interfaces)
	if err != nil {
		return err
	}

	client.Set(context.Background(), "interfaces", interfacesEncoding, 5*time.Minute)
	client.Set(context.Background(), "interfacesChangeTime", time.Now(), 5*time.Minute)

	return nil
}
