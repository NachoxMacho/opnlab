package proxmox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func Fetch(client *redis.Client) error {

	vms, err := GetVMs()
	if err != nil {
		return err
	}

	vmEncoding, err := json.Marshal(vms)
	if err != nil {
		return err
	}

	client.Set(context.Background(), "vms", vmEncoding, 5*time.Minute)

	nodes, err := GetNodes()
	if err != nil {
		return err
	}
	jsonEncoding, err := json.Marshal(nodes)
	if err != nil {
		return err
	}

	client.Set(context.Background(), "nodes", jsonEncoding, 5*time.Minute)

	return nil
}
