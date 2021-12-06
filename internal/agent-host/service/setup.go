package hostAgentService

import (
	"fmt"
	v1 "github.com/easysoft/zagent/cmd/agent-host/router/v1"
	agentConf "github.com/easysoft/zagent/internal/agent/conf"
	_fileUtils "github.com/easysoft/zagent/internal/pkg/lib/file"
	_stringUtils "github.com/easysoft/zagent/internal/pkg/lib/string"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type SetupService struct {
	syncMap sync.Map
}

func NewSetupService() *SetupService {
	srv := SetupService{}
	srv.GenWebsockifyTokens()

	return &srv
}

func (s *SetupService) GetToken(port string) (ret v1.VncTokenResp) {
	obj, ok := s.syncMap.Load(port)

	if !ok {
		return
	}

	ret = obj.(v1.VncTokenResp)

	return
}

func (s *SetupService) GenWebsockifyTokens() {
	port := 5901
	for port <= 5910 {
		portStr := strconv.Itoa(port)

		// uuid: 192.168.1.215:5901
		content := fmt.Sprintf("%s: %s:%s", _stringUtils.NewUuid(), agentConf.Inst.NodeIp, portStr)

		pth := filepath.Join(agentConf.Inst.DirToken, portStr+".txt")
		_fileUtils.WriteFile(pth, content)

		arr := strings.Split(content, ":")
		result := v1.VncTokenResp{
			Token: arr[0],
			Ip:    arr[1],
			Port:  arr[2],
		}
		s.syncMap.Store(portStr, result)

		port++
	}
}