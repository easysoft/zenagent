package main

import (
	"flag"
	hostInit "github.com/easysoft/zv/cmd/host/init"
	"github.com/easysoft/zv/internal/agent/conf"
	"github.com/easysoft/zv/internal/agent/utils/common"
	"github.com/easysoft/zv/internal/agent/utils/const"
	consts "github.com/easysoft/zv/internal/comm/const"
	"github.com/easysoft/zv/internal/pkg/lib/log"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"syscall"
)

var (
	help    bool
	flagSet *flag.FlagSet
	runMode string
)

// @title ZAgent宿主机API文档
// @version 1.0
// @contact.name API Support
// @contact.url https://github.com/easysoft/zv/issues
// @contact.email 462626@qq.com
func main() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
		os.Exit(0)
	}()

	flagSet = flag.NewFlagSet(consts.AppNameAgentHost, flag.ContinueOnError)

	flagSet.StringVar(&runMode, "t", agentConst.Host.ToString(), "")
	flagSet.StringVar(&agentConf.Inst.Server, "s", "http://127.0.0.1:8085", "")
	flagSet.StringVar(&agentConf.Inst.NodeName, "n", "", "")
	flagSet.StringVar(&agentConf.Inst.NodeIp, "i", "127.0.0.1", "")
	flagSet.IntVar(&agentConf.Inst.NodePort, "p", 8086, "")
	flagSet.StringVar(&agentConf.Inst.Language, "l", "zh", "")

	flagSet.BoolVar(&help, "h", false, "")

	action := ""
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	switch action {
	case "help", "-h":
		agentUtils.PrintUsage(consts.AppNameAgentHost)

	default:
		start()
	}
}

func start() {
	_logUtils.Init(consts.AppNameAgentHost)

	if err := flagSet.Parse(os.Args[1:]); err == nil {
		agentConf.Inst.RunMode = agentConst.RunMode(runMode)
		hostInit.Init()
	}
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
