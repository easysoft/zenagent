package serverService

import (
	"github.com/easysoft/zagent/internal/comm/domain"
	"github.com/easysoft/zagent/internal/server/repo"
)

type BuildService struct {
	BuildRepo *repo.BuildRepo `inject:""`
	QueueRepo *repo.QueueRepo `inject:""`

	QueueService *QueueService `inject:""`
}

func NewBuildService() *BuildService {
	return &BuildService{}
}

func (s BuildService) SaveResult(build domain.Build) (count int) {
	s.BuildRepo.SaveResult(build)

	po := s.BuildRepo.GetBuild(build.ID)

	s.QueueService.SetQueueResult(po.QueueId, po.Progress, po.Status)

	return
}
