package kvmService

import (
	agentConf "github.com/easysoft/zagent/internal/agent/conf"
	"github.com/easysoft/zagent/internal/comm/const"
	"github.com/easysoft/zagent/internal/comm/domain"
	_httpUtils "github.com/easysoft/zagent/internal/pkg/lib/http"
	_i118Utils "github.com/easysoft/zagent/internal/pkg/lib/i118"
	_logUtils "github.com/easysoft/zagent/internal/pkg/lib/log"
	"time"
)

type VmService struct {
	VmMap     map[string]domain.Vm
	TimeStamp int64

	LibvirtService *LibvirtService `inject:""`
}

func NewVmService() *VmService {
	s := VmService{}
	s.TimeStamp = time.Now().Unix()
	s.VmMap = map[string]domain.Vm{}

	return &s
}

func (s *VmService) Register(isBusy bool) {
	node := domain.VmNode{
		Node: domain.Node{Name: agentConf.Inst.NodeName, WorkDir: agentConf.Inst.WorkDir,
			Ip: agentConf.Inst.NodeIp, Port: agentConf.Inst.NodePort,
		},
	}

	if isBusy {
		node.ServiceStatus = consts.ServiceBusy
	} else {
		node.ServiceStatus = consts.ServiceActive
	}

	url := _httpUtils.GenUrl(agentConf.Inst.Server, "vms/register")
	resp, ok := _httpUtils.Post(url, node)

	if ok {
		_logUtils.Info(_i118Utils.I118Prt.Sprintf("success_to_register", agentConf.Inst.Server))
	} else {
		_logUtils.Info(_i118Utils.I118Prt.Sprintf("fail_to_register", agentConf.Inst.Server, resp))
	}
}

func (s *VmService) UpdateVmsStatus(vms []domain.Vm) {
	names := map[string]bool{}

	for _, vm := range vms {
		name := vm.Name
		names[name] = true

		if _, ok := s.VmMap[name]; ok { // update status in map
			v := s.VmMap[name]
			v.Status = vm.Status
			s.VmMap[name] = v
		} else { // update time then add
			if vm.FirstDetectedTime.IsZero() {
				vm.FirstDetectedTime = time.Now()
			}
			s.VmMap[name] = vm
		}
	}

	keys := s.getKeys(s.VmMap)
	for _, key := range keys {
		if !names[key] { // remove vm in map but not found this time
			delete(s.VmMap, key)
			continue
		}

		// destroy and remove timeout vm
		v := s.VmMap[key]
		if time.Now().Unix()-v.FirstDetectedTime.Unix() > consts.VmTimeout*60 { // timeout
			s.LibvirtService.DestroyVmByName(v.Name, true)
			delete(s.VmMap, key)
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