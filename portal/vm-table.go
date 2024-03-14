package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/NachoxMacho/opnlab/opnsense"
	"github.com/NachoxMacho/opnlab/proxmox"
)

func vmTable(c *fiber.Ctx) error {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	var leases []opnsense.DHCPLease
	result, err := redisClient.Get(context.Background(), "leases").Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(result), &leases)
	if err != nil {
		return err
	}
	var vms []proxmox.VM
	result, err = redisClient.Get(context.Background(), "vms").Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(result), &vms)
	if err != nil {
		return err
	}

	slices.SortStableFunc(vms, func(a, b proxmox.VM) int {
		return strings.Compare(a.Config.Name, b.Config.Name)
	})

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
	outputVMs := make([]VMOutputData, len(vms))
	for i, vm := range vms {
		outputVMs[i].ID = fmt.Sprintf("%d", vm.Stats.VMID)

		outputVMs[i].ID = fmt.Sprintf("%d", vm.Stats.VMID)
		outputVMs[i].Name = vm.Config.Name
		outputVMs[i].Status = vm.Stats.Status
		outputVMs[i].MaxCPUs = fmt.Sprintf("%d", vm.Stats.CPUs)
		outputVMs[i].CurrentMemory = HumanFileSize(float64(vm.Stats.Memory))
		outputVMs[i].MaxMemory = HumanFileSize(float64(vm.Stats.MaxMemory))
		outputVMs[i].MaxDisk = HumanFileSize(float64(vm.Stats.MaxDiskBytes))
		outputVMs[i].CurrentCPU = fmt.Sprintf("%.2f%%", vm.Stats.CPU*100)
		outputVMs[i].MACAddress = vm.Config.MACAddress()
		outputVMs[i].IPAddress = ""
		outputVMs[i].Tags = vm.Config.TagList()

		for _, lease := range leases {
			if strings.EqualFold(lease.MAC, vm.Config.MACAddress()) {
				outputVMs[i].IPAddress = lease.Address.String()
			}
		}
	}
	return c.Render("overview/vm-table", fiber.Map{"VMs": outputVMs})
}
