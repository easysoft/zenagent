package serverService

import (
	"github.com/easysoft/zagent/internal/comm/const"
	_domain "github.com/easysoft/zagent/internal/pkg/domain"
	"github.com/easysoft/zagent/internal/server/model"
	"github.com/easysoft/zagent/internal/server/repo"
	"github.com/easysoft/zagent/internal/server/service/vendors"
	"time"
)

type HuaweiCloudVmService struct {
	HostRepo    *repo.HostRepo    `inject:""`
	BackingRepo *repo.BackingRepo `inject:""`
	VmRepo      *repo.VmRepo      `inject:""`

	VmCommonService       *VmCommonService               `inject:""`
	HistoryService        *HistoryService                `inject:""`
	HuaweiCloudEcsService *vendors.HuaweiCloudEcsService `inject:""`
}

func (s HuaweiCloudVmService) CreateRemote(hostId, backingId, queueId uint) (result _domain.RpcResp) {
	host := s.HostRepo.Get(hostId)
	backing := s.BackingRepo.Get(backingId)
	backing.Name = s.VmCommonService.genTmplName(backing)

	vm := model.Vm{
		HostId: host.ID, HostName: host.Name,
		Status:     consts.VmCreated,
		OsCategory: backing.OsCategory,
		OsType:     backing.OsType,
		OsVersion:  backing.OsVersion,
		OsLang:     backing.OsLang,
		BackingId:  backing.ID,
	}
	s.VmRepo.Save(&vm) // save vm to db, then update name with id
	vm.Name = s.VmCommonService.genVmName(backing, vm.ID)
	s.VmRepo.UpdateVmName(vm)

	ecsClient, err := s.HuaweiCloudEcsService.CreateEcsClient(host.CloudKey, host.CloudSecret, host.CloudRegion)
	imgClient, err := s.HuaweiCloudEcsService.CreateImgClient(host.CloudKey, host.CloudSecret, host.CloudRegion)
	vpcClient, err := s.HuaweiCloudEcsService.CreateVpcClient(host.CloudKey, host.CloudSecret, host.CloudRegion)

	if err != nil {
		result.Fail(err.Error())
		s.VmCommonService.SaveVmCreationResult(result.IsSuccess(), result.Msg, queueId, vm.ID, "", "", "")
		return
	}

	huaweiCloudService := vendors.NewHuaweiCloudEcsService()
	vm.CouldInstId, _, err = huaweiCloudService.CreateInst(vm.Name, backing.Name, ecsClient, imgClient, vpcClient)
	if err != nil {
		result.Fail(err.Error())
		s.VmCommonService.SaveVmCreationResult(result.IsSuccess(), result.Msg, queueId, vm.ID, "", "", "")
		return
	}

	result.Pass("")

	for i := 0; i < 60; i++ {
		<-time.After(1 * time.Second)

		_, result.Msg, vm.NodeIp, vm.MacAddress, err = huaweiCloudService.QueryInst(vm.CouldInstId, ecsClient)
		if err != nil {
			result.Fail(err.Error())
			s.VmCommonService.SaveVmCreationResult(result.IsSuccess(), result.Msg, queueId, vm.ID, vm.VncAddress, "", "")
			return
		}

		if vm.NodeIp != "" {
			break
		}
	}

	s.VmRepo.UpdateVmCloudInst(vm)

	vm.VncAddress, _ = huaweiCloudService.QueryVnc(vm.CouldInstId, ecsClient)
	s.VmCommonService.SaveVmCreationResult(result.IsSuccess(), result.Msg, queueId, vm.ID, vm.VncAddress, "", "")

	return
}

func (s HuaweiCloudVmService) DestroyRemote(vmId, queueId uint) (result _domain.RpcResp) {
	vm := s.VmRepo.GetById(vmId)
	host := s.HostRepo.Get(vm.HostId)

	status := consts.VmDestroy

	ecsClient, err := s.HuaweiCloudEcsService.CreateEcsClient(host.CloudKey, host.CloudSecret, host.CloudRegion)
	if err != nil {
		status = consts.VmFailDestroy
	} else {
		err = s.HuaweiCloudEcsService.RemoveInst(vm.CouldInstId, ecsClient)
		if err != nil {
			status = consts.VmFailDestroy
		}
	}

	s.VmRepo.UpdateStatusByCloudInstId([]string{vm.CouldInstId}, status)

	s.HistoryService.Create(consts.Vm, vmId, queueId, "", status.ToString())

	return
}
