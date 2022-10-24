package handler

import (
	v1 "github.com/easysoft/zv/cmd/server/router/v1"
	consts "github.com/easysoft/zv/internal/pkg/const"
	"github.com/easysoft/zv/internal/pkg/domain"
	serverService "github.com/easysoft/zv/internal/server/service"
	_httpUtils "github.com/easysoft/zv/pkg/lib/http"
	_logUtils "github.com/easysoft/zv/pkg/lib/log"
	"github.com/kataras/iris/v12"
	"time"
)

type HostCtrl struct {
	BaseCtrl

	AssertService *serverService.AssertService `inject:""`
}

func NewHostCtrl() *HostCtrl {
	return &HostCtrl{}
}

func (c *HostCtrl) Register(ctx iris.Context) {
	req := v1.HostRegisterReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		_, _ = ctx.JSON(_httpUtils.RespData(consts.ResultFail, err.Error(), nil))
		return
	}

	_logUtils.Infof("remove: %v", ctx.Request().RemoteAddr)

	//success := c.AssertService.RegisterHost(req)
	//if !success {
	//	ctx.StopWithJSON(http.StatusInternalServerError, "register fail")
	//	return
	//}

	data := domain.RegisterResp{
		Token:           "123",
		ExpiredTimeUnix: time.Now().Unix() + 24*3600,
	}
	_, _ = ctx.JSON(data)

	return
}
