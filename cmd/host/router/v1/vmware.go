package v1

type VmWareReq struct {
	VmId          string `json:"vmId"` // for destroy
	VmUniqueName  string `json:"vmName"`
	VmBackingName string `json:"VmBackingName"`

	VmProcessors uint `json:"vmProcessors"`
	VmMemory     uint `json:"vmMemory"`

	UserName string `json:"userName"`
	Password string `json:"password"`

	StartAfterCreated bool `json:"startAfterCreated"`
}

type VmWareResp struct {
	Code int `json:"code"`
	Msg  int `json:"msg"`

	VmId       string `json:"vmId"`
	Name       string `json:"name"`
	Mac        string `json:"mac"`
	VncAddress string `json:"vncAddress"`
}

type VmWareParam struct {
	Processors int `json:"processors"`
	Memory     int `json:"memory"`
}
