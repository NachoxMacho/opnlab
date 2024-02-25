package proxmox

import "strings"

type APIVMConfig struct {
	Name    string `json:"name"`
	VMGenID string `json:"vmgenid"`
	Cores   string `json:"cores"`
	Net0    string `json:"net0"`
	SCSIHw  string `json:"scsihw"`
	Meta    string `json:"meta"`
	Sockets string `json:"sockets"`
	OSType  string `json:"ostype"`
	SCSI0   string `json:"scsi0"`
	Memory  string `json:"memory"`
	NUMA    string `json:"numa"`
	IDE2    string `json:"ide2"`
	SMBIOS1 string `json:"smbios1"`
	Boot    string `json:"boot"`
}

func (vm APIVMConfig) MACAddress() string {
	virtio := strings.Split(vm.Net0, ",")[0]

	return strings.Split(virtio, "=")[1]
}
