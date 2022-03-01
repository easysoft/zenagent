package hostAgentKvmService

import (
	"github.com/easysoft/zv/internal/comm/const"
	"github.com/easysoft/zv/internal/comm/domain"
	"github.com/libvirt/libvirt-go"
	"strings"
	"time"
)

type VmService struct {
	VmMapVar  map[string]domain.Vm
	TimeStamp int64

	LibvirtService *LibvirtService `inject:""`
}

func NewVmService() *VmService {
	s := VmService{}
	s.TimeStamp = time.Now().Unix()
	s.VmMapVar = map[string]domain.Vm{}

	return &s
}

func (s *VmService) GetVms() (vms []domain.Vm) {
	domains := s.LibvirtService.ListVm()

	for _, dom := range domains {
		vm := domain.Vm{}
		vm.Name, _ = dom.GetName()
		if strings.Index(vm.Name, "test-") != 0 {
			continue
		}

		vm.Status = consts.VmUnknown
		domainState, _, _ := dom.GetState()
		if domainState == libvirt.DOMAIN_RUNNING {
			vm.Status = consts.VmRunning
		} else if domainState == libvirt.DOMAIN_SHUTOFF || domainState == libvirt.DOMAIN_SHUTDOWN {
			vm.Status = consts.VmShutOff
		}

		vms = append(vms, vm)
	}

	return vms
}

func (s *VmService) UpdateVmMapAndDestroyTimeout(vms []domain.Vm) {
	names := map[string]bool{}

	for _, vm := range vms {
		name := vm.Name
		names[name] = true

		if _, ok := s.VmMapVar[name]; ok { // update status in map
			v := s.VmMapVar[name]
			v.Status = vm.Status
			s.VmMapVar[name] = v
		} else { // update time then add
			if vm.FirstDetectedTime.IsZero() {
				vm.FirstDetectedTime = time.Now()
			}
			s.VmMapVar[name] = vm
		}
	}

	keys := s.getKeys(s.VmMapVar)
	for _, key := range keys {
		if !names[key] { // remove vm in map but not found this time
			delete(s.VmMapVar, key)
			continue
		}

		// destroy and remove timeout vm
		v := s.VmMapVar[key]
		if time.Now().Unix()-v.FirstDetectedTime.Unix() > consts.WaitVmLifecycleTimeout { // timeout
			s.LibvirtService.DestroyVmByName(v.Name, true)
			delete(s.VmMapVar, key)
		}
	}
}

func (s *VmService) getKeys(m map[string]domain.Vm) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
