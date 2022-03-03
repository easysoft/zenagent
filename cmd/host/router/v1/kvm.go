package v1

import (
	"github.com/easysoft/zv/internal/comm/const"
)

type KvmReq struct {
	VmMacAddress string `json:"vmMacAddress"`
	VmTemplate   string `json:"vmTemplate"`
	VmBacking    string `json:"vmBacking"`

	VmUniqueName string `json:"vmUniqueName" example:"test-win10-x64-pro-zh_cn"`
	VmCpu        uint   `json:"vmCpu" example:"3"`
	VmMemorySize uint   `json:"vmMemorySize" example:"5120000"`
	VmDiskSize   uint   `json:"vmDiskSize" example:"30000"`
	//VmCdromSys    string `json:"vmCdromSys"`
	//VmCdromDriver string `json:"vmCdromDriver"`

	OsCategory consts.OsCategory `json:"osCategory" example:"windows"` // Enums consts.OsCategory
	OsType     consts.OsType     `json:"osType" example:"win10"`       // Enums consts.OsType
	OsVersion  string            `json:"osVersion" example:"x64-pro"`
	OsLang     consts.OsLang     `json:"osLang" example:"zh_cn"` // Enums consts.OsLang

	//StartAfterCreated bool `json:"startAfterCreated"`
}

type KvmReqClone struct {
	VmMacAddress string `json:"vmMacAddress"`
	VmUniqueName string `json:"vmUniqueName" example:"test-win10-x64-pro-zh_cn-clone1"`
	VmSrc        string `json:"vmSrc" example:"test-win10-x64-pro-zh_cn"`

	VmCpu        uint `json:"vmCpu" example:"3"`
	VmMemorySize uint `json:"vmMemorySize" example:"5120000"`
	VmDiskSize   uint `json:"vmDiskSize" example:"30000"`

	//StartAfterCreated bool `json:"startAfterCreated"`
}

type KvmResp struct {
	Name        string `json:"name"`
	MacAddress  string `json:"macAddress"`
	VncPort     string `json:"vncPort"`
	ImagePath   string `json:"imagePath"`
	BackingPath string `json:"backingPath"`
}
type KvmRespTempl struct {
	Name string `json:"name"`
	Type string `json:"type"`
	UUID string `json:"uuid"`

	CpuCoreNum  uint   `json:"cpuCoreNum"`
	MemoryValue uint   `json:"memoryValue"`
	MemoryUnit  string `json:"memoryUnit"`

	OsArch        string `json:"osArch"`
	MacAddress    string `json:"macAddress"`
	DiskFile      string `json:"diskFile"`
	BackingFile   string `json:"backingFile"`
	BackingFormat string `json:"backingFormat"`

	VncPost int `json:"memoryValue"`
}
