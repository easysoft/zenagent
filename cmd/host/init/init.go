package hostInit

import (
	"fmt"
	hostCron "github.com/easysoft/zv/cmd/host/cron"
	hostRouter "github.com/easysoft/zv/cmd/host/router"
	agentConf "github.com/easysoft/zv/internal/agent/conf"
	consts "github.com/easysoft/zv/internal/comm/const"
	hostAgentService "github.com/easysoft/zv/internal/host/service"
	hostKvmService "github.com/easysoft/zv/internal/host/service/kvm"
	"github.com/easysoft/zv/internal/pkg/db"
	_commonUtils "github.com/easysoft/zv/internal/pkg/lib/common"
	serverConf "github.com/easysoft/zv/internal/server/conf"
	"github.com/facebookgo/inject"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	_ "github.com/easysoft/zv/res/host/docs"
)

func Init() {
	agentConf.Init(consts.AppNameAgentHost)

	irisServer := NewServer(nil)
	irisServer.App.Logger().SetLevel(serverConf.Inst.LogLevel)

	router := hostRouter.NewRouter(irisServer.App)
	injectObj(router)

	router.App()

	iris.RegisterOnInterrupt(func() {
		defer _db.GetInst().Close()
	})

	err := irisServer.Serve()
	if err != nil {
		panic(err)
	}
}

func injectObj(router *hostRouter.Router) {
	// inject
	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	err := g.Provide(
		// setup
		&inject.Object{Value: hostAgentService.NewSetupService()},

		// cron
		&inject.Object{Value: hostCron.NewAgentCron()},

		// service
		&inject.Object{Value: hostKvmService.NewLibvirtService()},
		&inject.Object{Value: hostKvmService.NewVmService()},

		// controller
		&inject.Object{Value: router},
	)

	if err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}

	err = g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}
}

type Server struct {
	App       *iris.Application
	AssetFile http.FileSystem
}

func NewServer(assetFile http.FileSystem) *Server {
	app := iris.Default()
	return &Server{
		App:       app,
		AssetFile: assetFile,
	}
}

func (s *Server) Serve() error {
	if _commonUtils.IsPortInUse(agentConf.Inst.NodePort) {
		panic(fmt.Sprintf("端口 %d 已被使用", agentConf.Inst.NodePort))
	}

	err := s.App.Run(
		iris.Addr(fmt.Sprintf("%s:%d", agentConf.Inst.NodeIp, agentConf.Inst.NodePort)),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(time.RFC3339),
	)

	if err != nil {
		return err
	}

	return nil
}
