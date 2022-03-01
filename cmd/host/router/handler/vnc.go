package hostHandler

import (
	hostAgentService "github.com/easysoft/zv/internal/host/service"
	_httpUtils "github.com/easysoft/zv/internal/pkg/lib/http"
	"github.com/kataras/iris/v12"
)

type VncCtrl struct {
	SetupService *hostAgentService.SetupService `inject:""`
}

func NewVncCtrl() *VncCtrl {
	return &VncCtrl{}
}

// GetToken
// @summary 根据VNC Port获取Token
// @Accept json
// @Produce json
// @Param port query string true "Vnc Port"
// @Success 200 {object} _httpUtils.Response{iris.Map} "code = success? 1 : 0"
// @Router /api/v1/vnc/getToken [get]
func (c *VncCtrl) GetToken(ctx iris.Context) {
	port := ctx.URLParam("port")

	if port == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, "no port param", nil))
		return
	}

	ret := c.SetupService.GetToken(port)
	if ret.Token == "" {
		_, _ = ctx.JSON(_httpUtils.ApiRes(iris.StatusInternalServerError, "token not found", nil))
		return
	}

	ctx.JSON(_httpUtils.ApiRes(iris.StatusOK, "success", ret))

	return
}
