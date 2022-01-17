package v1

type MultiPassReq struct {
	VmUniqueName string `json:"vmName"`
	VmMemory     uint   `json:"vmMemory"`
	Cpus         uint   `json:"cpus"`
	Disk         uint   `json:"disk"`
	ImgFrom      string `json:"imgFrom"`
	ImagePath    string `json:"imgPath"`

	UserName string `json:"userName"`
	Password string `json:"password"`
}

type MultiPassResp struct {
	Name    string `json:"name"`
	Cpus    string `json:"cpus"`
	Memory  string `json:"memory"`
	Network string `json:"network"`
	State   string `json:"state"`
	Image   string `json:"image"`
	IPv4    string `json:"ipv4"`
}

type MultiPassParam struct {
	Processors int `json:"processors"`
	Memory     int `json:"memory"`
}
