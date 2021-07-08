package router

import (
	"github.com/easysoft/zagent/cmd/agent/router/handler"
	agentConf "github.com/easysoft/zagent/internal/agent/conf"
	agentConst "github.com/easysoft/zagent/internal/agent/utils/const"
	_i118Utils "github.com/easysoft/zagent/internal/pkg/lib/i118"
	_logUtils "github.com/easysoft/zagent/internal/pkg/lib/log"
	"github.com/smallnest/rpcx/server"
	"strconv"
)

type Router struct {
	ArithCtrl   *handler.ArithCtrl   `inject:""`
	VmCtrl      *handler.VmCtrl      `inject:""`
	JobCtrl     *handler.JobCtrl     `inject:""`
	PerformCtrl *handler.PerformCtrl `inject:""`
}

func NewRouter() *Router {
	router := &Router{}
	return router
}

func (r *Router) App() {
	addr := agentConf.Inst.NodeIp + ":" + strconv.Itoa(agentConf.Inst.NodePort)
	srv := server.NewServer()

	if agentConf.Inst.RunMode == agentConst.Host {
		srv.RegisterName("arith", r.ArithCtrl, "")
		srv.RegisterName("vm", r.VmCtrl, "")

	} else if agentConf.Inst.RunMode == agentConst.Vm {
		srv.RegisterName("job", r.JobCtrl, "")

	} else if agentConf.Inst.RunMode == agentConst.Android || agentConf.Inst.RunMode == agentConst.Ios {
		srv.RegisterName("job", r.JobCtrl, "")

	} else if agentConf.Inst.RunMode == agentConst.Machine {
		srv.RegisterName("perform", r.PerformCtrl, "") // add to agent side queue

	}

	_logUtils.Info(_i118Utils.Sprintf("start_server", addr))
	err := srv.Serve("tcp", addr)
	if err != nil {
		_logUtils.Infof(_i118Utils.Sprintf("fail_to_start_server", addr, err.Error()))
	}
}
