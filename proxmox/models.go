package proxmox

type APINode struct {
	Status string `json:"status" redis:"status"`
	Name   string `json:"node" redis:"name"`
	// Size of the disk in bytes
	Disk           int     `json:"disk" redis:"disk"`
	CPU            float32 `json:"cpu" redis:"cpu"`
	MaxCPU         int     `json:"maxcpu" redis:"maxcpu"`
	Level          string  `json:"level" redis:"level"`
	ID             string  `json:"id" redis:"id"`
	SSLFingerprint string  `json:"ssl_fingerprint" redis:"ssl_fingerprint"`
	Uptime         int     `json:"uptime" redis:"uptime"`
	Type           string  `json:"type" redis:"type"`
	MaxDiskBytes   int     `json:"maxdisk" redis:"maxdisk"`
	Memory         int     `json:"mem" redis:"memory"`
}

type VM struct {
	Stats  APIVM
	Config APIVMConfig
}

type APIVM struct {
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	DiskWrite    int     `json:"diskwrite"`
	Memory       int     `json:"mem"`
	Disk         int     `json:"disk"`
	DiskRead     int     `json:"diskread"`
	CPUs         int     `json:"cpus"`
	NetIn        int     `json:"netin"`
	NetOut       int     `json:"netout"`
	Uptime       int     `json:"uptime"`
	MaxDiskBytes int     `json:"maxdisk"`
	PID          int     `json:"pid"`
	MaxMemory    int     `json:"maxmem"`
	VMID         int     `json:"vmid"`
	CPU          float32 `json:"cpu"`
}

type APIStorage struct {
	Enabled      string `json:"enabled"`
	Shared       string `json:"shared"`
	Total        string `json:"total"`
	Type         string `json:"type"`
	Active       string `json:"active"`
	Available    string `json:"avail"`
	Content      string `json:"content"`
	Name         string `json:"storage"`
	Used         string `json:"used"`
	UsedFraction string `json:"used_fraction"`
}
