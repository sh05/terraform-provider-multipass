package common

// MultipassInstance represents a Multipass VM instance
type MultipassInstance struct {
	Name      string                 `json:"name"`
	State     string                 `json:"state"`
	IPv4      []string               `json:"ipv4,omitempty"`
	Release   string                 `json:"release,omitempty"`
	ImageHash string                 `json:"image_hash,omitempty"`
	Load      []float64              `json:"load,omitempty"`
	DiskUsage string                 `json:"disk_usage,omitempty"`
	Memory    map[string]interface{} `json:"memory,omitempty"`
	Mounts    map[string]interface{} `json:"mounts,omitempty"`
}

// MultipassInstanceList represents the list response from multipass list
type MultipassInstanceList struct {
	List []MultipassInstance `json:"list"`
}

// MultipassInstanceInfo represents detailed info from multipass info
type MultipassInstanceInfo struct {
	Info   map[string]MultipassInstance `json:"info"`
	Errors []string                     `json:"errors,omitempty"`
}

// LaunchOptions represents options for launching a new instance
type LaunchOptions struct {
	Name      string
	Image     string
	CPU       string
	Memory    string
	Disk      string
	CloudInit string
}
