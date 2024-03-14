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

	return nil
}
