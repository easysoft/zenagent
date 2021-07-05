package domain

import (
	commConst "github.com/easysoft/zagent/internal/comm/const"
	"time"
)

type Vm struct {
	Id        int
	BackingId int
	HostId    int

	Name        string
	Tmpl        string
	Backing     string
	ImagePath   string
	BackingPath string

	OsCategory commConst.OsCategory
	OsType     commConst.OsType
	OsVersion  string
	OsLang     commConst.OsLang

	Status            commConst.VmStatus
	DestroyAt         *time.Time
	FirstDetectedTime *time.Time

	PublicIp   string
	PublicPort int
	MacAddress string
	RpcPort    int
	SshPort    int
	VncPort    int
	WorkDir    string

	DefPath          string
	DiskSize         uint // M
	MemorySize       uint // M
	CdromSys         string
	CdromDriver      string
	ResolutionHeight int
	ResolutionWidth  int
}
