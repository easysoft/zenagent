package hostRouter

import (
	"github.com/easysoft/zv/cmd/agent/router/handler"
	hostHandler "github.com/easysoft/zv/cmd/host/router/handler"
	_httpUtils "github.com/easysoft/zv/internal/pkg/lib/http"
	serverConf "github.com/easysoft/zv/internal/server/conf"
	"github.com/kataras/iris/v12"
)

type Router struct {
	api *iris.Application

	ArithCtrl *handler.ArithCtrl   `inject:""`
	JobCtrl   *hostHandler.JobCtrl `inject:""`

	KvmCtrl    *hostHandler.KvmCtrl    `inject:""`
	VmWareCtrl *hostHandler.VmWareCtrl `inject:""`
	VncCtrl    *hostHandler.VncCtrl    `inject:""`
}

func NewRouter(app *iris.Application) *Router {
	router := &Router{api: app}
	return router
}

func (r *Router) App() {
	iris.LimitRequestBodySize(serverConf.Inst.Options.UploadMaxSize)
	r.api.UseRouter(_httpUtils.CrsAuth())

	app := r.api.Party("/api").AllowMethods(iris.MethodOptions)
	{
		v1 := app.Party("/v1")
		{
			v1.PartyFunc("/job", func(client iris.Party) {
				client.Post("/add", r.JobCtrl.Add).Name = "创建任务"
			})
			v1.PartyFunc("/kvm", func(client iris.Party) {
				client.Post("/listTempl", r.KvmCtrl.ListTempl).Name = "克隆虚机"
				client.Post("/create", r.KvmCtrl.Create).Name = "创建虚机"
				client.Post("/clone", r.KvmCtrl.Clone).Name = "克隆虚机"
				client.Post("/{name:string}/destroy", r.KvmCtrl.Destroy).Name = "摧毁虚机"
				client.Post("/{name:string}/reboot", r.KvmCtrl.Reboot).Name = "重启虚机"
				client.Post("/{name:string}/suspend", r.KvmCtrl.Suspend).Name = "暂停虚机"
				client.Post("/{name:string}/resume", r.KvmCtrl.Resume).Name = "恢复虚机"
			})
			v1.PartyFunc("/vmware", func(client iris.Party) {
				client.Post("/create", r.VmWareCtrl.Create).Name = "创建虚机"
				client.Post("/destroy", r.VmWareCtrl.Destroy).Name = "摧毁虚机"
			})
			v1.PartyFunc("/vnc", func(client iris.Party) {
				client.Get("/getToken", r.VncCtrl.GetToken).Name = "获取VNC的Token"
			})
		}
	}
}
