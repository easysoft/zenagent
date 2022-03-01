package domain

type MultiPass struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	IPv4        string `json:"ipv4"`
	Release     string `json:"release"`
	ImageHash   string `json:"imageHash"`
	Load        string `json:"load"`
	DiskUsage   string `json:"diskUsage"`
	MemoryUsage string `json:"memoryUsage"`
	Mounts      string `json:"mounts"`
	Cpus        string `json:"cpus"`
	Memory      string `json:"memory"`
	Network     string `json:"network"`
	Image       string `json:"image"`
}
