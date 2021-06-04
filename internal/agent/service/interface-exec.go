package agentService

import (
	"encoding/json"
	agentConf "github.com/easysoft/zagent/internal/agent/conf"
	commDomain "github.com/easysoft/zagent/internal/comm/domain"
	_const "github.com/easysoft/zagent/internal/pkg/const"
	_fileUtils "github.com/easysoft/zagent/internal/pkg/libs/file"
	_httpUtils "github.com/easysoft/zagent/internal/pkg/libs/http"
	_i118Utils "github.com/easysoft/zagent/internal/pkg/libs/i118"
	_logUtils "github.com/easysoft/zagent/internal/pkg/libs/log"
	"strings"
)

type InterfaceExecService struct {
	CommonService
	InterfaceRequestService *InterfaceRequestService `inject:""`
}

func NewInterfaceExecService() *InterfaceExecService {
	return &InterfaceExecService{}
}

func (s *InterfaceExecService) ExecProcessor(build *commDomain.Build, result *commDomain.TestResult,
	processor commDomain.TestProcessor) {
	tp := processor.Type
	if tp == _const.Simple {
		for _, child := range processor.Children {
			if _, ok := child.(commDomain.TestProcessor); ok {
				s.ExecProcessor(build, result, child.(commDomain.TestProcessor))
			} else {
				s.InterfaceRequestService.Request(build, child.(commDomain.TestInterface))
			}
		}
	}
}

func (s *InterfaceExecService) UploadResult(build commDomain.Build, result commDomain.TestResult) {
	zipFile := build.WorkDir + "testResult.zip"
	err := _fileUtils.ZipFiles(zipFile, build.ProjectDir, strings.Split(build.AutomatedTest.ResultFiles, ","))
	if err != nil {
		_logUtils.Errorf(_i118Utils.Sprintf("fail_to_zip_test_result",
			zipFile, build.ProjectDir, build.AutomatedTest.ResultFiles, err))
	}

	result.Payload = build

	uploadResultUrl := _httpUtils.GenUrl(agentConf.Inst.Server, "build/upload")
	files := []string{zipFile}
	extraParams := map[string]string{}
	json, _ := json.Marshal(result)
	extraParams["result"] = string(json)

	_fileUtils.Upload(uploadResultUrl, files, extraParams)
}
