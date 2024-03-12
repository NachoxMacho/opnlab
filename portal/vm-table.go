package portal

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/NachoxMacho/opnlab/opnsense"
	"github.com/NachoxMacho/opnlab/proxmox"
)

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
