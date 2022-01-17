package hostHandler

import (
	v1 "github.com/easysoft/zv/cmd/agent-host/router/v1"
	multiPassService "github.com/easysoft/zv/internal/agent-host/service/multipass"
	_const "github.com/easysoft/zv/internal/pkg/const"
	_httpUtils "github.com/easysoft/zv/internal/pkg/lib/http"
	"github.com/kataras/iris/v12"
)

type MultiPassCtrl struct {
	MultiPassService *multiPassService.MultiPassService `inject:""`
}

func NewMultiPassCtrl() *MultiPassCtrl {
	return &MultiPassCtrl{}
}

// List
// @summary 获取MultiPass虚拟机
// @Accept json
// @Produce json
// @Param task body v1.MultiPassReq true "MultiPass Request Object"
// @Success 200 {object} _httpUtils.Response{data=v1.MultiPassResp} "code = success? 1 : 0"
// @Router /api/v1/MultiPass/list [post]
func (c *MultiPassCtrl) List(ctx iris.Context) {
	domains, err := c.MultiPassService.GetVmList()
	if domains == nil || err != nil {
		ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "fail to find MultiPass vm", err))
		return
	}
	var mps []v1.MultiPassResp
	mp := v1.MultiPassResp{}
	for _, v := range domains {
		mp.Name = v.Name
		mp.Memory = v.Memory
		mp.Cpus = v.Cpus
		mp.IPv4 = v.IPv4
		mp.State = v.State
		mps = append(mps, mp)
	}

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to get MultiPass vm", mps))

	return
}

// Create
// @summary 创建MultiPass虚拟机
// @Accept json
// @Produce json
// @Param task body v1.MultiPassReq true "MultiPass Request Object"
// @Success 200 {object} _httpUtils.Response{data=v1.MultiPassResp} "code = success? 1 : 0"
// @Router /api/v1/MultiPass/create [post]
func (c *MultiPassCtrl) Create(ctx iris.Context) {
	req := v1.MultiPassReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, err.Error(), nil))
		return
	}
	domains, err := c.MultiPassService.CreateVm(&req, false)
	if err != nil {
		ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, "fail to create vm", err))
		return
	}

	vm := v1.MultiPassResp{
		Name: domains.Name, //VncPort:  strconv.Itoa(MultiPassVncPort),
	}

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to create MultiPass vm", vm))

	return
}

// Destroy
// @summary 摧毁MultiPass虚拟机
// @Accept json
// @Produce json
// @Param task body v1.MultiPassReq true "MultiPass Request Object"
// @Success 200 {object} _httpUtils.Response{data=int} "code = success? 1 : 0"
// @Router /api/v1/multipass/{name}/destroy [post]
func (c *MultiPassCtrl) Destroy(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	if name == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(_const.ResultFail, "vm name is empty", nil))
		return
	}

	c.MultiPassService.DestroyVm(name)

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to destroy MultiPass vm", name))
	return
}

// Reboot
// @summary 重启MultiPass虚拟机
// @Accept json
// @Produce json
// @Param name path string true "MultiPass Name"
// @Success 200 {object} _httpUtils.Response{} "code = success? 1 : 0"
// @Router /api/v1/MultiPass/{name}/reboot [post]
func (c *MultiPassCtrl) Reboot(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	if name == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(_const.ResultFail, "vm name is empty", nil))
		return
	}

	c.MultiPassService.RebootVmByName(name)

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to reboot vm", name))
	return
}

// Suspend
// @summary 暂停MultiPass虚拟机
// @Accept json
// @Produce json
// @Param name path string true "MultiPass Name"
// @Success 200 {object} _httpUtils.Response{} "code = success? 1 : 0"
// @Router /api/v1/MultiPass/{name}/suspend [post]
func (c *MultiPassCtrl) Suspend(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	if name == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(_const.ResultFail, "vm name is empty", nil))
		return
	}

	c.MultiPassService.SuspendVmByName(name)

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to suspend vm", name))
	return
}

// Resume
// @summary 恢复MultiPass虚拟机
// @Accept json
// @Produce json
// @Param name path string true "MultiPass Name"
// @Success 200 {object} _httpUtils.Response{} "code = success? 1 : 0"
// @Router /api/v1/MultiPass/{name}/resume [post]
func (c *MultiPassCtrl) Resume(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	if name == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(_const.ResultFail, "vm name is empty", nil))
		return
	}

	c.MultiPassService.ResumeVmByName(name)

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to resume vm", name))
	return
}

func (c *MultiPassCtrl) GetToken(ctx iris.Context) {
	port := ctx.URLParam("port")

	if port == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, "no port param", nil))
		return
	}

	ret := c.MultiPassService.GetToken(port)
	if ret.Token == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, "token not found", nil))
		return
	}

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success", ret))

	return
}
