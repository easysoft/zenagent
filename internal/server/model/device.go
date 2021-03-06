package model

import (
	"github.com/easysoft/zv/internal/comm/const"
	"github.com/easysoft/zv/internal/comm/domain"
)

type Device struct {
	BaseModel

	// from node register
	domain.DeviceInst

	// info to maintain
	Name   string
	Make   string
	Brand  string
	Series string

	CpuMake         string
	CpuModel        string
	Memory          int
	Storage         int
	BatteryCapacity int

	OsType    consts.OsDevice
	OsLevel   string
	OsVersion string
}

func NewDevice() Device {
	device := Device{}

	return device
}

func (Device) TableName() string {
	return "biz_device"
}
