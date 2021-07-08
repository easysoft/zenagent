package testing

import (
	_domain "github.com/easysoft/zagent/internal/pkg/domain"
	"github.com/easysoft/zagent/internal/server/model"
	"github.com/easysoft/zagent/internal/server/repo"
	serverService "github.com/easysoft/zagent/internal/server/service"
	commonService "github.com/easysoft/zagent/internal/server/service/common"
)

type SeleniumService struct {
	BuildRepo *repo.BuildRepo `inject:""`
	VmRepo    *repo.VmRepo    `inject:""`

	RpcService   *commonService.RpcService   `inject:""`
	QueueService *serverService.QueueService `inject:""`
}

func NewSeleniumService() *SeleniumService {
	return &SeleniumService{}
}

func (s SeleniumService) Start(queue model.Queue) (result _domain.RpcResp) {
	vmId := queue.VmId
	vm := s.VmRepo.GetById(vmId)

	build := model.NewSeleniumBuildPo(queue, vm)
	s.BuildRepo.Save(&build)

	build = s.BuildRepo.GetBuild(build.ID)

	result = s.RpcService.SeleniumTest(build)

	return
}
