package hostHandler

import (
	v1 "github.com/easysoft/zv/cmd/host/router/v1"
	vmWareService "github.com/easysoft/zv/internal/host/service/vmware"
	_httpUtils "github.com/easysoft/zv/internal/pkg/lib/http"
	"github.com/kataras/iris/v12"
)

type VmWareCtrl struct {
	VmWareService *vmWareService.VmWareService `inject:""`
}

func NewVmWareCtrl() *VmWareCtrl {
	return &VmWareCtrl{}
}

// Create
// @summary 创建VmWare虚拟机
// @Accept json
// @Produce json
// @Param task body v1.VmWareReq true "VmWare Request Object"
// @Success 200 {object} _httpUtils.Response{data=v1.VmWareResp} "code = success? 1 : 0"
// @Router /api/v1/vmware/create [post]
func (c *VmWareCtrl) Create(ctx iris.Context) {
	req := v1.VmWareReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, err.Error(), nil))
		return
	}

	id, macAddress, err := c.VmWareService.CreateVm(&req, true)
	if err != nil {
		ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "fail to create VmWare vm", err))
		return
	}

	vm := v1.VmWareResp{
		VmId: id,
		Name: req.VmUniqueName,
		Mac:  macAddress,
		//VncPort:  strconv.Itoa(VmWareVncPort),
	}

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to create VmWare vm", vm))

	return
}

// Destroy
// @summary 摧毁VmWare虚拟机
// @Accept json
// @Produce json
// @Param task body v1.VmWareReq true "VmWare Request Object"
// @Success 200 {object} _httpUtils.Response{data=int} "code = success? 1 : 0"
// @Router /api/v1/kvm/{name}/destroy [post]
func (c *VmWareCtrl) Destroy(ctx iris.Context) {
	req := v1.VmWareReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.VmWareService.DestroyVm(&req, true)

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to destroy VmWare vm", req.VmId))
	return
}
