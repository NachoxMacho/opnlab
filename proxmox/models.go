package proxmox

type APINode struct {
	Status string `json:"status"`
	Name   string `json:"node"`
	// Size of the disk in bytes
	Disk           int     `json:"disk"`
	CPU            float32 `json:"cpu"`
	MaxCPU         int     `json:"maxcpu"`
	Level          string  `json:"level"`
	ID             string  `json:"id"`
	SSLFingerprint string  `json:"ssl_fingerprint"`
	Uptime         int     `json:"uptime"`
	Type           string  `json:"type"`
	MaxDiskBytes   int     `json:"maxdisk"`
	Memory         int     `json:"mem"`
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
