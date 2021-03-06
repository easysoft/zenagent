package model

import (
	"github.com/easysoft/zv/internal/comm/const"
)

type VmBacking struct {
	BaseModel

	Name string `json:"name"`
	Path string `json:"path"`
	Size int    `json:"size"`

	OsCategory consts.OsCategory `json:"osCategory"`
	OsType     consts.OsType     `json:"osType"`
	OsLang     consts.OsLang     `json:"osLang"`

	OsVersion string `json:"osVersion"`
	OsBuild   string `json:"osBuild"`
	OsBits    string `json:"osBits"`

	ResolutionHeight  int  `json:"resolutionHeight"`
	ResolutionWidth   int  `json:"resolutionWidth"`
	SuggestCpuCount   uint `json:"suggestCpuCount"`
	SuggestMemorySize uint `json:"suggestMemorySize"` // MB
	SuggestDiskSize   uint `json:"suggestDiskSize"`   // MB

	SysIsoId    uint `json:"sysIsoId"`
	DriverIsoId uint `json:"driverIsoId"`

	Browsers []Browser `gorm:"many2many:backing_browser_r;"`
}

func (VmBacking) TableName() string {
	return "biz_vm_baking"
}
