package agentUtils

import (
	"fmt"
	consts "github.com/easysoft/zagent/internal/comm/const"
	_commonUtils "github.com/easysoft/zagent/internal/pkg/lib/common"
	_i118Utils "github.com/easysoft/zagent/internal/pkg/lib/i118"
	_logUtils "github.com/easysoft/zagent/internal/pkg/lib/log"
	agentRes "github.com/easysoft/zagent/res/agent-vm"
	"github.com/fatih/color"
	"io/ioutil"
	"path/filepath"
)

var (
	usageFile = filepath.Join("res", "doc", "usage.txt")
)

func PrintUsage() {
	_logUtils.PrintColor(_i118Utils.Sprintf("usage"), color.FgCyan)

	usage := ReadResData(usageFile)

	app := consts.AppNameAgent
	if _commonUtils.IsWin() {
		app += ".exe"
	}
	usage = fmt.Sprintf(usage, app)
	fmt.Printf("%s\n", usage)
}

func ReadResData(path string) string {
	isRelease := _commonUtils.IsRelease()

	var jsonStr string
	if isRelease {
		data, _ := agentRes.Asset(path)
		jsonStr = string(data)
	} else {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			jsonStr = "fail to read " + path
		} else {
			str := string(buf)
			jsonStr = _commonUtils.RemoveBlankLine(str)
		}
	}

	return jsonStr
}
