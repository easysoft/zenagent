package vmWareService

import (
	"fmt"
	v1 "github.com/easysoft/zv/cmd/host/router/v1"
	"github.com/easysoft/zv/internal/comm/domain"
	_logUtils "github.com/easysoft/zv/internal/pkg/lib/log"
	vmwareService "github.com/easysoft/zv/internal/server/service/vendors/vmware"
	"time"
)

type VmWareService struct {
	VmMapVar  map[string]domain.Vm
	TimeStamp int64
}

func NewVmWareService() *VmWareService {
	s := VmWareService{}
	s.TimeStamp = time.Now().Unix()
	s.VmMapVar = map[string]domain.Vm{}

	return &s
}

func (s *VmWareService) CreateVm(req *v1.VmWareReq, removeSameName bool) (id, macAddress string, err error) {
	client := vmwareService.NewVMWareService() // 8697
	err = client.Connect(fmt.Sprintf("https://127.0.0.1:8697"), req.UserName, req.Password)
	if err != nil {
		_logUtils.Errorf("Connect to vmware err %s", err.Error())
		return
	}

	// create machine
	vmInst, err := client.CreateVm(req.VmBackingName, req.VmUniqueName, req.VmProcessors, req.VmMemory)
	if err != nil {
		_logUtils.Errorf("Create vmware vm err %s", err.Error())
		return
	}

	id = vmInst.IdVM
	macAddress = vmInst.MacAddress

	return
}

func (s *VmWareService) DestroyVm(req *v1.VmWareReq, removeDiskImage bool) (err error) {
	client := vmwareService.NewVMWareService()
	err = client.Connect(fmt.Sprintf("https://127.0.0.1:8697"), req.UserName, req.Password)
	if err != nil {
		_logUtils.Errorf("Connect to vmware err %s", err.Error())
		return
	}

	client.DestroyVm(req.VmId)

	return
}
