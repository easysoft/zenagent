package vmHandler

import (
	"github.com/easysoft/zv/internal/agent/service"
	"github.com/easysoft/zv/internal/comm/domain"
	"github.com/easysoft/zv/internal/pkg/lib/http"
	"github.com/kataras/iris/v12"
)

type JobCtrl struct {
	JobService *agentService.JobService `inject:""`
}

func NewJobCtrl() *JobCtrl {
	return &JobCtrl{}
}

// Add
// @summary 创建任务
// @Accept json
// @Produce json
// @Param task body domain.Build true "Build Request Object"
// @Success 200 {object} _httpUtils.Response{} "code = success? 1 : 0"
// @Router /api/v1/vmware/create [post]
func (c *JobCtrl) Add(ctx iris.Context) {
	build := domain.Build{}
	if err := ctx.ReadJSON(&build); err != nil {
		ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, err.Error(), nil))
		return
	}

	//size := c.JobService.GetTaskSize()
	//if size == 0 {
	c.JobService.AddTask(build)
	//} else {
	//	reply.Fail(fmt.Sprintf("already has %d jobs to be done.", size))
	//}

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to add job", nil))
	return
}
