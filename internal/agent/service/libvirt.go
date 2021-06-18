package agentService

import "C"
import (
	"fmt"
	agentConf "github.com/easysoft/zagent/internal/agent/conf"
	commConst "github.com/easysoft/zagent/internal/comm/const"
	commDomain "github.com/easysoft/zagent/internal/comm/domain"
	_logUtils "github.com/easysoft/zagent/internal/pkg/libs/log"
	"github.com/libvirt/libvirt-go"
	"github.com/libvirt/libvirt-go-xml"
	"path/filepath"
)

const (
	ConnStrLocal  = "qemu:///system"
	ConnStrRemote = "qemu+ssh://192.168.0.56:22/system?keyfile=~/.ssh/id_rsa"
)

var (
	Conn *libvirt.Connect
)

type LibvirtService struct {
	QemuService *QemuService `inject:""`
}

func NewLibvirtService() *LibvirtService {
	return &LibvirtService{}
}

func (s *LibvirtService) CreateVm(vm *commDomain.Vm) (
	dom *libvirt.Domain, macAddress string, vncPort int, err error) {
	s.setVmProps(vm)

	srcXml := s.GetVmDef(vm.Src)

	basePath := ""
	if vm.Base != "" {
		basePath = filepath.Join(agentConf.Inst.DirBase, vm.Base)
	}
	basePath += ".qcow2"

	vmXml := ""
	rawPath := filepath.Join(agentConf.Inst.DirImage, vm.Name+".qcow2")
	vmXml, vm.MacAddress, _ = s.QemuService.GenVmDef(srcXml, vm.Name, rawPath, basePath, 0)

	if err != nil || vm.DiskSize == 0 {
		_logUtils.Errorf("wrong vm disk size %d, err %s", vm.DiskSize, err.Error())
		return
	}

	s.QemuService.createDiskFile(basePath, vm.Name, vm.DiskSize)

	dom, err = Conn.DomainCreateXML(vmXml, 0)

	if err == nil {
		newXml := ""
		newXml, err = dom.GetXMLDesc(0)
		if err != nil {
			return
		}

		newDomCfg := &libvirtxml.Domain{}
		err = newDomCfg.Unmarshal(newXml)

		macAddress = newDomCfg.Devices.Interfaces[0].MAC.Address
		vncPort = newDomCfg.Devices.Graphics[0].VNC.Port
	}

	return
}

func (s *LibvirtService) GetVm(name string) (dom *libvirt.Domain) {
	s.Connect(ConnStrLocal)
	defer func() {
		if res, _ := Conn.Close(); res != 0 {
			_logUtils.Errorf("close() == %d, expected 0", res)
		}
	}()

	dom, err := Conn.LookupDomainByName(name)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	return
}

func (s *LibvirtService) StartVm(dom *libvirt.Domain) (err error) {
	err = dom.Create()
	return
}
func (s *LibvirtService) DestroyVm(dom *libvirt.Domain) (err error) {
	err = dom.Destroy()
	return
}
func (s *LibvirtService) UndefineVm(dom *libvirt.Domain) (err error) {
	err = dom.Undefine()
	return
}

func (s *LibvirtService) GetVmDef(name string) (xml string) {
	dom := s.GetVm(name)
	if dom == nil {
		return
	}

	xml, _ = dom.GetXMLDesc(0)

	return
}

func (s *LibvirtService) Connect(str string) {
	if Conn != nil {
		active, err := Conn.IsAlive()
		if err != nil {
			_logUtils.Errorf(err.Error())
		}
		if active {
			return
		}
	}

	var err error
	Conn, err = libvirt.NewConnect(str)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	active, err := Conn.IsAlive()
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}
	if !active {
		_logUtils.Errorf("not active")
	}
	return
}

func (s *LibvirtService) setVmProps(vm *commDomain.Vm) {
	osCategory := commConst.Windows
	osType := commConst.Win10
	osVersion := "x64-pro"
	osLang := commConst.ZH_CN

	vm.Base = fmt.Sprintf("%s/%s/%s-%s", osCategory.ToString(), osType.ToString(),
		osVersion, osLang.ToString())
}
