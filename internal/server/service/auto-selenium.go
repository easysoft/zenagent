package serverService

import (
	consts "github.com/easysoft/zv/internal/comm/const"
	_domain "github.com/easysoft/zv/internal/pkg/domain"
	"github.com/easysoft/zv/internal/server/model"
	"github.com/easysoft/zv/internal/server/repo"
	commonService "github.com/easysoft/zv/internal/server/service/common"
)

type SeleniumService struct {
	BuildRepo *repo.BuildRepo `inject:""`
	VmRepo    *repo.VmRepo    `inject:""`

	RpcService       *commonService.RpcService       `inject:""`
	QueueService     *QueueService                   `inject:""`
	HistoryService   *HistoryService                 `inject:""`
	WebSocketService *commonService.WebSocketService `inject:""`
}

func NewSeleniumService() *SeleniumService {
	return &SeleniumService{}
}

func (s SeleniumService) RunRemote(queue model.Queue) (result _domain.RpcResp) {
	vmId := queue.VmId
	vm := s.VmRepo.GetById(vmId)

	build := model.NewSeleniumBuildPo(queue, vm)
	s.BuildRepo.Save(&build)

	s.HistoryService.Create(consts.Build, build.ID, queue.ID, consts.ProgressCreated, consts.StatusCreated.ToString())

	build = s.BuildRepo.GetBuild(build.ID)

	result = s.RpcService.SeleniumTest(build)

	return
}
