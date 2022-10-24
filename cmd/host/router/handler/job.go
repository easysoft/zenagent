package hostHandler

import (
	consts "github.com/easysoft/zv/internal/pkg/const"
	"github.com/easysoft/zv/internal/pkg/domain"
	agentService "github.com/easysoft/zv/internal/pkg/service"
	_httpUtils "github.com/easysoft/zv/pkg/lib/http"
	"github.com/kataras/iris/v12"
)

type JobCtrl struct {
	JobService *agentService.JobService `inject:""`
}

func NewJobCtrl() *JobCtrl {
	return &JobCtrl{}
}

func (c *JobCtrl) Add(ctx iris.Context) {
	build := domain.Build{}
	if err := ctx.ReadJSON(&build); err != nil {
		ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
		return
	}

	//size := c.JobService.GetTaskSize()
	//if size == 0 {
	c.JobService.AddTask(build)
	//} else {
	//	reply.Fail(fmt.Sprintf("already has %d jobs to be done.", size))
	//}

	ctx.JSON(_httpUtils.RespData(consts.ResultPass, "success to add job", nil))
	return
}
