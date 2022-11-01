package consts

import (
	"fmt"
	_const "github.com/easysoft/zagent/pkg/const"
	"path/filepath"
	"sync"
	"time"
)

var ()

var (
	ConfigFile     = filepath.Join("conf", AppNameAgent+".yaml")
	LogDir         = fmt.Sprintf("log%s", _const.PthSep)
	ControlActions = []string{"start", "stop", "restart", "install", "uninstall"}

	DownloadDir   = ""
	NovncDir      = ""
	WebsockifyDir = ""

	AuthSecret  = ""
	AuthToken   = ""
	ExpiredDate = time.Now()

	Verbose          = false
	ExistVncPortMap  = sync.Map{}
	ExistHttpPortMap = sync.Map{}
	ExistSshPortMap  = sync.Map{}
)
