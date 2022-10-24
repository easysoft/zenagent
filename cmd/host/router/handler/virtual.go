package hostHandler

import (
	v1 "github.com/easysoft/zv/cmd/host/router/v1"
	kvmService "github.com/easysoft/zv/internal/host/service/kvm"
	virtualService "github.com/easysoft/zv/internal/host/service/virtual"
	agentConf "github.com/easysoft/zv/internal/pkg/conf"
	consts "github.com/easysoft/zv/internal/pkg/const"
	"github.com/easysoft/zv/internal/pkg/domain"
	_httpUtils "github.com/easysoft/zv/pkg/lib/http"
	"github.com/kataras/iris/v12"
	"sync"
)

var (
	vmMacMap = sync.Map{}
)

type VirtualCtrl struct {
	NatService   *virtualService.NatService   `inject:""`
	NoVncService *virtualService.NoVncService `inject:""`
	KvmService   *kvmService.KvmService       `inject:""`
}

func NewVirtualCtrl() *VirtualCtrl {
	return &VirtualCtrl{}
}

func (c *VirtualCtrl) NotifyHost(ctx iris.Context) {
	req := domain.VmNotifyReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
		return
	}

	data := domain.VmNotifyResp{
		Secret: agentConf.Inst.Secret,
	}

	// get vm ip
	vmIp, err := c.KvmService.GetVmIpByMac(req.MacAddress)
	if err != nil {
		ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
		return
	}

	// map vm agent port to host
	//vmAgentPortMapped, err := natHelper.GetValidPort()
	//if err != nil {
	//	ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
	//	return
	//}
	//err = natHelper.ForwardPort(vmIp, consts.AgentServicePost, vmAgentPortMapped, consts.Http)
	//if err != nil {
	//	ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
	//	return
	//}
	//
	//data.AgentPortOnHost = vmAgentPortMapped
	data.Ip = vmIp

	ctx.JSON(_httpUtils.RespData(consts.ResultPass, "success to refresh secret", data))
	return
}

// @summary 新增虚拟机到宿主机端口的映射
// @Accept json
// @Produce json
// @Param VmPortMapReq body v1.KvmReq true "Kvm Request Object"
// @Success 200 {object} _domain.Response{data=v1.VmPortMapResp} "code = success? 1 : 0"
// @Router /api/v1/virtual/addVmPortMap [post]
func (c *VirtualCtrl) AddVmPortMap(ctx iris.Context) {
	req := v1.VmPortMapReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
		return
	}

	resp, err := c.NatService.AddVmPortMap(req)
	if err != nil {
		ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
		return
	}

	ctx.JSON(_httpUtils.RespData(consts.ResultPass, "success", resp))
}

// @summary 删除虚拟机到宿主机的端口映射
// @Accept json
// @Produce json
// @Param VmPortMapReq body v1.KvmReq true "Kvm Request Object"
// @Success 200 {object} _domain.Response "code = success? 1 : 0"
// @Router /api/v1/virtual/removeVmPortMap [post]
func (c *VirtualCtrl) RemoveVmPortMap(ctx iris.Context) {
	req := v1.VmPortMapReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
		return
	}

	err := c.NatService.RemoveVmPortMap(req)
	if err != nil {
		ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
		return
	}

	ctx.JSON(_httpUtils.RespData(consts.ResultPass, "success", nil))
}

// @summary 根据VNC Port获取Token
// @Accept json
// @Produce json
// @Param port query string true "VNC Port"
// @Success 200 {object} _domain.Response{data=v1.VncTokenResp} "code = success? 1 : 0"
// @Router /api/v1/virtual/getToken [get]
func (c *VirtualCtrl) GetToken(ctx iris.Context) {
	port := ctx.URLParam("port")

	if port == "" {
		_, _ = ctx.JSON(_httpUtils.RespData(consts.ResultFail, "no port param", nil))
		return
	}

	ret, _ := c.NoVncService.GetToken(port)
	if ret.Token == "" {
		_, _ = ctx.JSON(_httpUtils.RespData(consts.ResultFail, "token not found", nil))
		return
	}

	ctx.JSON(_httpUtils.RespData(consts.ResultPass, "success", ret))

	return
}
