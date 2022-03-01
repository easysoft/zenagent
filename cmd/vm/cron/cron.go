package vmCron

import (
	"fmt"
	consts "github.com/easysoft/zv/internal/comm/const"
	_cronUtils "github.com/easysoft/zv/internal/pkg/lib/cron"
	_dateUtils "github.com/easysoft/zv/internal/pkg/lib/date"
	_logUtils "github.com/easysoft/zv/internal/pkg/lib/log"
	vmAgentService "github.com/easysoft/zv/internal/vm/service"
	"github.com/kataras/iris/v12"
	"sync"
	"time"
)

type CronService struct {
	syncMap   sync.Map
	VmService *vmAgentService.VmService `inject:""`
}

func NewAgentCron() *CronService {
	inst := &CronService{}
	inst.Init()
	return inst
}

func (s *CronService) Init() {
	s.syncMap.Store("isRunning", false)
	s.syncMap.Store("lastCompletedTime", int64(0))

	_cronUtils.AddTask(
		"check",
		fmt.Sprintf("@every %ds", consts.AgentCheckInterval),
		func() {
			isRunning, _ := s.syncMap.Load("isRunning")
			lastCompletedTime, _ := s.syncMap.Load("lastCompletedTime")

			if isRunning.(bool) || time.Now().Unix()-lastCompletedTime.(int64) < consts.AgentCheckInterval {
				_logUtils.Infof("skip this iteration " + _dateUtils.DateTimeStr(time.Now()))
				return
			}
			s.syncMap.Store("isRunning", true)

			s.VmService.Check()

			s.syncMap.Store("isRunning", false)
			s.syncMap.Store("lastCompletedTime", time.Now().Unix())
		},
	)

	iris.RegisterOnInterrupt(func() {
		_cronUtils.Stop()
	})
}
