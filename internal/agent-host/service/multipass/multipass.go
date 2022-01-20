package multiPassService

import (
	"errors"
	"fmt"
	v1 "github.com/easysoft/zv/cmd/agent-host/router/v1"
	agentConf "github.com/easysoft/zv/internal/agent/conf"
	"github.com/easysoft/zv/internal/comm/domain"
	_fileUtils "github.com/easysoft/zv/internal/pkg/lib/file"
	_logUtils "github.com/easysoft/zv/internal/pkg/lib/log"
	_shellUtils "github.com/easysoft/zv/internal/pkg/lib/shell"
	_stringUtils "github.com/easysoft/zv/internal/pkg/lib/string"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	mpTokenPrefix = "multipass-"

	cmdMpls      = "multipass ls"
	cmdMpInfo    = "multipass info %s"
	cmdMpStart   = "multipass start %s"
	cmdMpStop    = "multipass stop %s"
	cmdMpSuspend = "multipass suspend %s"
	cmdMpReboot  = "multipass restart %s"
	cmdMpDelete  = "multipass delete %s"
	cmdMpPurge   = "multipass purge"
	cmdMpStopAll = "multipass stop --all"
)

type MultiPassService struct {
	syncMap sync.Map
}

func (s *MultiPassService) CreateVm(req *v1.MultiPassReq, removeSameName bool) (dom domain.MultiPass, err error) {
	cmdMpLaunch := "multipass launch "
	name := req.VmUniqueName
	cpus := req.Cpus
	disk := req.Disk / 1000
	mem := req.VmMemory
	filePath := req.ImagePath
	imgFrom := req.ImgFrom

	if name != "" {
		vm := s.GetVmInfo(name)
		if vm.Name != "" {
			msg := fmt.Sprintf("vm %s has existed", name)
			_logUtils.Errorf(msg)
			err = errors.New(msg)
			return
		}
		cmdMpLaunch = cmdMpLaunch + fmt.Sprintf("-n %s ", name)
	}

	if cpus >= 1 {
		cmdMpLaunch = cmdMpLaunch + fmt.Sprintf("-c %v ", cpus)
	}
	if disk >= 0 {
		cmdMpLaunch = cmdMpLaunch + fmt.Sprintf("-d %vG ", disk)
	}
	if mem >= 0 {
		cmdMpLaunch = cmdMpLaunch + fmt.Sprintf("-m %vM ", mem)
	}
	if filePath != "" {
		filePath = filepath.Join(agentConf.Inst.DirImage, filePath)
		if imgFrom == "file" || imgFrom == "url" {
			cmdMpLaunch = cmdMpLaunch + fmt.Sprintf("%s://%s", imgFrom, filePath)
		}
	}
	_, err = _shellUtils.ExeShellWithOutput(cmdMpLaunch)
	dom = s.GetVmInfo(name)
	if dom.State != "Running" && err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	return
}

func (s *MultiPassService) RebootVmByName(name string) (dom domain.MultiPass, err error) {
	vm := s.GetVmInfo(name)
	if vm.Name == "" {
		msg := fmt.Sprintf("vm %s not found", name)
		_logUtils.Errorf(msg)
		err = errors.New(msg)
		return
	}
	_shellUtils.ExeShellWithOutput(fmt.Sprintf(cmdMpReboot, name))
	if s.GetVmInfo(name).State != "Running" && err != nil {
		_logUtils.Errorf("RebootVM error %s", err.Error())
		return
	}
	return
}

func (s *MultiPassService) DestroyVm(name string) (dom domain.MultiPass, err error) {
	if s.GetVmInfo(name).Name == "" {
		msg := fmt.Sprintf("vm %s not found", name)
		_logUtils.Errorf(msg, name)
		err = errors.New(msg)
		return
	}

	_, err = _shellUtils.ExeShellWithOutput(fmt.Sprintf(cmdMpDelete+"&&"+cmdMpPurge, name))
	if s.GetVmInfo(name).Name == "" && err != nil {
		_logUtils.Errorf("DestroyVM error %s", err.Error())
		return
	}
	return
}

func (s *MultiPassService) SuspendVmByName(name string) (dom domain.MultiPass, err error) {
	if s.GetVmInfo(name).Name == "" {
		msg := fmt.Sprintf("vm %s not found", name)
		_logUtils.Errorf(msg)
		err = errors.New(msg)
		return
	}

	_, err = _shellUtils.ExeShellWithOutput(fmt.Sprintf(cmdMpSuspend, name))
	if s.GetVmInfo(name).State != "Suspended" && err != nil {
		_logUtils.Errorf("SuspendVM error %s", err.Error())
		return
	}
	return
}

func (s *MultiPassService) ResumeVmByName(name string) (dom domain.MultiPass, err error) {
	if s.GetVmInfo(name).Name == "" {
		msg := fmt.Sprintf("vm %s not found", name)
		_logUtils.Errorf(msg)
		err = errors.New(msg)
		return
	}

	_, err = _shellUtils.ExeShellWithOutput(fmt.Sprintf(cmdMpStart, name))
	if s.GetVmInfo(name).State != "Running" && err != nil {
		_logUtils.Errorf("ResumeVM error %s", err.Error())
		return
	}
	return
}

func (s *MultiPassService) GetToken(port string) (ret v1.VncTokenResp) {
	s.GenWebsockifyTokens()
	obj, ok := s.syncMap.Load(mpTokenPrefix + port)

	if !ok {
		return
	}

	ret = obj.(v1.VncTokenResp)

	return
}

func (s *MultiPassService) GetVmList() (doms []domain.MultiPass, err error) {
	outRets, err := _shellUtils.ExeShellWithOutput(cmdMpls)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}
	dom := domain.MultiPass{}
	var rets []string
	for i := 1; i < len(outRets)-1; i++ {
		rets = strings.Fields(outRets[i])

		dom.Name = rets[0]
		dom.State = rets[1]
		dom.IPv4 = rets[2]
		for i := 3; i < len(rets); i++ {
			dom.Image += rets[i]
		}

		doms = append(doms, dom)
	}
	return
}

func (s *MultiPassService) GetVmInfo(name string) (dom domain.MultiPass) {
	outRets, err := _shellUtils.ExeShellWithOutput(fmt.Sprintf(cmdMpInfo, name))
	if strings.Contains(outRets[0], "info failed:") || err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	vmInfoMap := parseOutput(outRets[:len(outRets)-1])
	dom = domain.MultiPass{
		Name:        vmInfoMap["Name"],
		State:       vmInfoMap["State"],
		IPv4:        vmInfoMap["IPv4"],
		Release:     vmInfoMap["Release"],
		Load:        vmInfoMap["Load"],
		DiskUsage:   vmInfoMap["Diskusage"],
		MemoryUsage: vmInfoMap["Memoryusage"],
		Mounts:      vmInfoMap["Mounts"],
		ImageHash:   vmInfoMap["Imagehash"],
	}

	return
}

func parseOutput(lines []string) (vmInfoMap map[string]string) {
	vmInfoMap = make(map[string]string, len(lines))

	for _, v := range lines {
		rets := strings.Replace(v, " ", "", -1)
		rets = strings.Replace(rets, "\n", "", -1)
		ret := strings.Split(rets, ":")
		vmInfoMap[ret[0]] = ret[1]
	}
	return
}

func (s *MultiPassService) GenWebsockifyTokens() { // create tokenFile
	vms, _ := s.GetVmList()
	port := 5901
	for _, v := range vms {
		if v.State != "Running" {
			continue
		}
		portStr := strconv.Itoa(port)

		// uuid: vmIp:5901
		content := fmt.Sprintf("%s: %s:%s", _stringUtils.NewUuid(), v.IPv4, portStr)

		pth := filepath.Join(agentConf.Inst.DirToken, mpTokenPrefix+portStr+".txt")
		_fileUtils.WriteFile(pth, content)

		arr := strings.Split(content, ":")
		result := v1.VncTokenResp{
			Token: arr[0],
			Ip:    v.IPv4,
			Port:  arr[2],
		}
		s.syncMap.Store(mpTokenPrefix+portStr, result)
	}
}
