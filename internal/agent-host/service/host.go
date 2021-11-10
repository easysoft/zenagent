package hostAgentService

import (
	hostKvmService "github.com/easysoft/zagent/internal/agent-host/service/kvm"
	agentConf "github.com/easysoft/zagent/internal/agent/conf"
	agentService "github.com/easysoft/zagent/internal/agent/service"
	testingService "github.com/easysoft/zagent/internal/agent/service/testing"
	agentZentaoService "github.com/easysoft/zagent/internal/agent/service/zentao"
	"github.com/easysoft/zagent/internal/comm/const"
	"github.com/easysoft/zagent/internal/comm/domain"
	_i118Utils "github.com/easysoft/zagent/internal/pkg/lib/i118"
	_logUtils "github.com/easysoft/zagent/internal/pkg/lib/log"
	"strings"
)

type HostService struct {
	VmService *hostKvmService.VmService `inject:""`

	JobService  *agentService.JobService   `inject:""`
	TestService *testingService.RunService `inject:""`

	ZentaoService agentZentaoService.ZentaoService `inject:""`
}

func NewHostService() *HostService {
	return &HostService{}
}

func (s *HostService) Check() {
	// is running，register busy
	if s.JobService.IsRunning() {
		s.Register(true)
		return
	}

	// no task to run, submit free
	if s.JobService.GetTaskSize() == 0 {
		s.Register(false)
		return
	}

	// has task to run，register busy, then run
	job := s.JobService.PeekJob()
	s.Register(true)

	s.JobService.StartTask()

	s.PassEnvsToContainerIfNeeded(&job)
	s.TestService.Run(&job)

	s.JobService.RemoveTask()
	s.JobService.EndTask()
}

func (s *HostService) Register(isBusy bool) {
	host := domain.HostNode{
		Node: domain.Node{Ip: agentConf.Inst.NodeIp, Port: agentConf.Inst.NodePort},
	}
	if isBusy {
		host.Status = consts.HostBusy
	} else {
		host.Status = consts.HostOnline
	}

	host.Vms = s.VmService.GetVms()
	s.VmService.UpdateVmMapAndDestroyTimeout(host.Vms)

	url := s.ZentaoService.GenUrl(agentConf.Inst.Server, "api.php/v1/host/register")
	resp, ok := s.ZentaoService.Post(url, host, false)

	if ok {
		_logUtils.Info(_i118Utils.I118Prt.Sprintf("success_to_register", agentConf.Inst.Server))
	} else {
		_logUtils.Info(_i118Utils.I118Prt.Sprintf("fail_to_register", agentConf.Inst.Server, resp))
	}
}

func (s *HostService) PassEnvsToContainerIfNeeded(build *domain.Build) {
	str := "docker run"
	if strings.Index(build.BuildCommands, "docker run") > -1 {
		newStr := str
		for _, env := range strings.Split(build.EnvVars, "\n") {
			newStr += " -e " + env + " "
		}

		build.BuildCommands = strings.Replace(build.BuildCommands, str, newStr, -1)
	}
}
