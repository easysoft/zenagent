package main

import (
	"flag"
	"fmt"
	ini "github.com/easysoft/zv/cmd/server/init"
	"github.com/easysoft/zv/internal/comm/const"
	"github.com/easysoft/zv/internal/pkg/lib/log"
	"os"
)

var (
	version      = "master"
	printVersion = flag.Bool("v", false, "打印版本号")
	printRouter  = flag.Bool("r", false, "打印路由列表")
)

// @title ZAgent服务器API文档
// @version 1.0
// @contact.name API Support
// @contact.url https://github.com/easysoft/zv/issues
// @contact.email 462626@qq.com
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] [command]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  -c <path>\n")
		fmt.Fprintf(os.Stderr, "    设置配置文件路径\n")
		fmt.Fprintf(os.Stderr, "  -v <true or false> 默认为: false\n")
		fmt.Fprintf(os.Stderr, "    打印版本号\n")
		fmt.Fprintf(os.Stderr, "  -s <true or false> 默认为: false\n")
		fmt.Fprintf(os.Stderr, "    填充基础数据\n")
		fmt.Fprintf(os.Stderr, "  -p <true or false> 默认为: true\n")
		fmt.Fprintf(os.Stderr, "    同步权限\n")
		fmt.Fprintf(os.Stderr, "  -r <true or false> 默认为: false\n")
		fmt.Fprintf(os.Stderr, "    打印路由列表\n")
		fmt.Fprintf(os.Stderr, "\n")
		//flag.PrintDefaults()
	}
	flag.Parse()

	_logUtils.Init(consts.AppNameServer)

	ini.Init(version, printVersion, printRouter)
}
