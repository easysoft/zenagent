package hostHandler

import (
	v1 "github.com/easysoft/zv/cmd/host/router/v1"
	hostKvmService "github.com/easysoft/zv/internal/host/service/kvm"
	_const "github.com/easysoft/zv/internal/pkg/const"
	_httpUtils "github.com/easysoft/zv/internal/pkg/lib/http"
	"github.com/kataras/iris/v12"
	"strconv"
)

type KvmCtrl struct {
	VmService      *hostKvmService.VmService      `inject:""`
	LibvirtService *hostKvmService.LibvirtService `inject:""`
}

func NewKvmCtrl() *KvmCtrl {
	return &KvmCtrl{}
}

// ListTempl
// @summary 获取KVM虚拟机模板信息
// @Produce json
// @Success 200 {object} _httpUtils.Response{data=[]v1.KvmRespTempl} "code = success? 1 : 0"
// @Router /api/v1/kvm/listTempl [get]
func (c *KvmCtrl) ListTempl(ctx iris.Context) {
	domainCfgs, err := c.LibvirtService.ListTmpl()

	if err != nil {
		ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, "fail to list vm tmpl", err))
		return
	}

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to list vm tmpl", domainCfgs))

	return
}

// Create
// @summary 创建KVM虚拟机
// @Accept json
// @Produce json
// @Param kvmReq body v1.KvmReq true "Kvm Request Object"
// @Success 200 {object} _httpUtils.Response{data=v1.KvmResp} "code = success? 1 : 0"
// @Router /api/v1/kvm/create [post]
func (c *KvmCtrl) Create(ctx iris.Context) {
	req := v1.KvmReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, err.Error(), nil))
		return
	}

	dom, vmVncPort, vmRawPath, vmBackingPath, err := c.LibvirtService.CreateVm(&req, true)

	if err != nil {
		ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, "fail to create vm", err))
		return
	}

	vmName, _ := dom.GetName()
	vm := v1.KvmResp{
		Name:        vmName,
		MacAddress:  req.VmMacAddress,
		VncPort:     strconv.Itoa(vmVncPort),
		ImagePath:   vmRawPath,
		BackingPath: vmBackingPath,
	}

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to create vm", vm))

	return
}

// Clone
// @summary 克隆KVM虚拟机
// @Accept json
// @Produce json
// @Param kvmReqClone body v1.KvmReqClone true "Kvm Request Object"
// @Success 200 {object} _httpUtils.Response{data=v1.KvmResp} "code = success? 1 : 0"
// @Router /api/v1/kvm/create [post]
func (c *KvmCtrl) Clone(ctx iris.Context) {
	req := v1.KvmReqClone{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, err.Error(), nil))
		return
	}

	if req.VmSrc == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, "request vmSrc field can not be empty.", nil))
		return
	}

	dom, vmVncPort, vmRawPath, vmBackingPath, err := c.LibvirtService.CloneVm(&req, true)

	if err != nil {
		ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, "fail to clone vm", err))
		return
	}

	vmName, _ := dom.GetName()
	vm := v1.KvmResp{
		Name:        vmName,
		MacAddress:  req.VmMacAddress,
		VncPort:     strconv.Itoa(vmVncPort),
		ImagePath:   vmRawPath,
		BackingPath: vmBackingPath,
	}

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to create vm", vm))

	return
}

// Destroy
// @summary 摧毁KVM虚拟机
// @Accept json
// @Produce json
// @Param name path string true "Kvm Name"
// @Success 200 {object} _httpUtils.Response{} "code = success? 1 : 0"
// @Router /api/v1/kvm/{name}/destroy [post]
func (c *KvmCtrl) Destroy(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	if name == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(_const.ResultFail, "vm name is empty", nil))
		return
	}

	c.LibvirtService.DestroyVmByName(name, true)

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to destroy vm", name))
	return
}

// Reboot
// @summary 重启KVM虚拟机
// @Accept json
// @Produce json
// @Param name path string true "Kvm Name"
// @Success 200 {object} _httpUtils.Response{} "code = success? 1 : 0"
// @Router /api/v1/kvm/{name}/reboot [post]
func (c *KvmCtrl) Reboot(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	if name == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(_const.ResultFail, "vm name is empty", nil))
		return
	}

	c.LibvirtService.RebootVmByName(name)

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to reboot vm", name))
	return
}

// Suspend
// @summary 暂停KVM虚拟机
// @Accept json
// @Produce json
// @Param name path string true "Kvm Name"
// @Success 200 {object} _httpUtils.Response{} "code = success? 1 : 0"
// @Router /api/v1/kvm/{name}/suspend [post]
func (c *KvmCtrl) Suspend(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	if name == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(_const.ResultFail, "vm name is empty", nil))
		return
	}

	c.LibvirtService.SuspendVmByName(name)

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to suspend vm", name))
	return
}

// Resume
// @summary 恢复KVM虚拟机
// @Accept json
// @Produce json
// @Param name path string true "Kvm Name"
// @Success 200 {object} _httpUtils.Response{} "code = success? 1 : 0"
// @Router /api/v1/kvm/{name}/resume [post]
func (c *KvmCtrl) Resume(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	if name == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(_const.ResultFail, "vm name is empty", nil))
		return
	}

	c.LibvirtService.ResumeVmByName(name)

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to resume vm", name))
	return
}

//func (c *KvmCtrl) Boot(ctx iris.Context) {
//	req :=v1.KvmReq{}
//	if err := ctx.ReadJSON(&req); err != nil {
//		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, err.Error(), nil))
//		return
//	}
//
//	c.LibvirtService.BootVmByName(req.VmUniqueName)
//
//	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to boot vm", req.VmUniqueName))
//	return
//}
//func (c *KvmCtrl) Shutdown(ctx iris.Context) {
//	req :=v1.KvmReq{}
//	if err := ctx.ReadJSON(&req); err != nil {
//		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, err.Error(), nil))
//		return
//	}
//
//	c.LibvirtService.ShutdownVmByName(req.VmUniqueName)
//
//	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success to shutdown vm", req.VmUniqueName))
//	return
//}
