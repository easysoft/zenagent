package hostAgentKvmService

import "C"
import (
	"encoding/json"
	"fmt"
	v1 "github.com/easysoft/zv/cmd/host/router/v1"
	agentConf "github.com/easysoft/zv/internal/agent/conf"
	"github.com/easysoft/zv/internal/comm/const"
	"github.com/easysoft/zv/internal/comm/domain"
	_fileUtils "github.com/easysoft/zv/internal/pkg/lib/file"
	_logUtils "github.com/easysoft/zv/internal/pkg/lib/log"
	_stringUtils "github.com/easysoft/zv/internal/pkg/lib/string"
	"github.com/libvirt/libvirt-go"
	"github.com/libvirt/libvirt-go-xml"
	"path/filepath"
	"strings"
)

const (
	LibvirtConnStrLocal  = "qemu:///system"
	LibvirtConnStrRemote = "qemu+ssh://%s:22/system?socket=/var/run/libvirt/libvirt-sock"
)

var ()

type LibvirtService struct {
	LibvirtConn *libvirt.Connect
	QemuService *QemuService `inject:""`
}

func NewLibvirtService() *LibvirtService {
	connStr := ""
	if agentConf.Inst.Host == "" {
		connStr = LibvirtConnStrLocal
	} else {
		connStr = fmt.Sprintf(LibvirtConnStrRemote, agentConf.Inst.Host)
	}

	s := &LibvirtService{}
	s.Connect(connStr)

	return s
}

func (s *LibvirtService) ListTmpl() (ret []v1.KvmRespTempl, err error) {
	if s.LibvirtConn == nil {
		return
	}

	domains, err := s.LibvirtConn.ListAllDomains(0)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	for _, domain := range domains {
		name, _ := domain.GetName()
		if strings.Index(name, "tmpl-") > -1 {
			newXml := ""
			newXml, err = domain.GetXMLDesc(0)
			if err != nil {
				continue
			}

			domainCfg := &libvirtxml.Domain{}
			err = domainCfg.Unmarshal(newXml)

			tmpl := v1.KvmRespTempl{
				Name: domainCfg.Name,
				Type: domainCfg.Type,
				UUID: domainCfg.UUID,

				CpuCoreNum:  domainCfg.VCPU.Value,
				MemoryValue: domainCfg.Memory.Value,
				MemoryUnit:  domainCfg.Memory.Unit,

				OsArch:     domainCfg.OS.Type.Arch,
				MacAddress: domainCfg.Devices.Interfaces[0].MAC.Address,
			}

			if len(domainCfg.Devices.Interfaces) > 0 {
				tmpl.VncPost = domainCfg.Devices.Graphics[0].VNC.Port
			}

			mainDiskIndex := s.QemuService.GetMainDiskIndex(domainCfg)
			tmpl.DiskFile = domainCfg.Devices.Disks[mainDiskIndex].Source.File.File

			backingStore := domainCfg.Devices.Disks[mainDiskIndex].BackingStore
			if backingStore != nil && backingStore.Source.File != nil {
				tmpl.BackingFile = backingStore.Source.File.File
				tmpl.BackingFormat = backingStore.Format.Type
			}

			ret = append(ret, tmpl)
		}
	}

	return
}

func (s *LibvirtService) CreateVm(req *v1.KvmReq, removeSameName bool) (dom *libvirt.Domain,
	vmVncPort int, vmRawPath, vmBackingPath string, err error) {

	reqMsg, err := json.Marshal(req)
	_logUtils.Infof("%s", reqMsg)

	vmMacAddress := req.VmMacAddress
	vmUniqueName := req.VmUniqueName
	vmBackingPath = filepath.Join(agentConf.Inst.DirKvm, req.VmBacking)
	vmTemplateName := req.VmTemplate
	vmCpu := req.VmCpu
	vmMemorySize := req.VmMemorySize
	vmDiskSize := req.VmDiskSize

	if removeSameName {
		s.DestroyVmByName(vmUniqueName, true)
	}

	// gen tmpl xml definition
	tmplXml := s.GetVmDef(vmTemplateName)
	tmplDomCfg := &libvirtxml.Domain{}
	err = tmplDomCfg.Unmarshal(tmplXml)
	if err != nil {
		return
	}

	// get BackingPath
	mainDiskIndex := s.QemuService.GetMainDiskIndex(tmplDomCfg)
	if vmBackingPath == "" {
		// use tmpl vm's image as backing path
		vmBackingPath = tmplDomCfg.Devices.Disks[mainDiskIndex].Source.File.File
	}

	vmXml := ""
	vmXml, vmRawPath, _ = s.QemuService.GenVmDef(tmplXml, vmMacAddress, vmUniqueName, vmBackingPath, vmCpu, vmMemorySize)
	if err != nil {
		_logUtils.Errorf("err gen vm xml, err %s", err.Error())
		return
	}

	err = s.QemuService.createDiskFile(vmBackingPath, vmRawPath, vmDiskSize)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	dom, err = s.LibvirtConn.DomainCreateXML(vmXml, libvirt.DOMAIN_NONE)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	// get new vm info
	newXml := ""
	newXml, err = dom.GetXMLDesc(0)
	if err != nil {
		return
	}

	newDomCfg := &libvirtxml.Domain{}
	err = newDomCfg.Unmarshal(newXml)

	vmMacAddress = newDomCfg.Devices.Interfaces[0].MAC.Address
	vmVncPort = newDomCfg.Devices.Graphics[0].VNC.Port

	return
}

func (s *LibvirtService) CloneVm(req *v1.KvmReqClone, removeSameName bool) (dom *libvirt.Domain,
	vmVncPort int, vmRawPath, vmBackingPath string, err error) {

	reqMsg, err := json.Marshal(req)
	_logUtils.Infof("%s", reqMsg)

	// get src vm config
	tmplXml := s.GetVmDef(req.VmSrc)
	srcDomCfg := &libvirtxml.Domain{}
	err = srcDomCfg.Unmarshal(tmplXml)
	if err != nil {
		return
	}

	vmMacAddress := req.VmMacAddress
	vmUniqueName := req.VmUniqueName
	if removeSameName {
		s.DestroyVmByName(vmUniqueName, true)
	}

	vmCpu := req.VmCpu
	vmMemorySize := req.VmMemorySize
	vmDiskSize := req.VmDiskSize

	// update empty values from tmpl
	if vmCpu == 0 {
		vmCpu = srcDomCfg.VCPU.Value
	}
	if vmMemorySize == 0 {
		vmCpu = srcDomCfg.Memory.Value
	}

	// get BackingPath
	mainDiskIndex := s.QemuService.GetMainDiskIndex(srcDomCfg)
	if vmBackingPath == "" {
		// try to use the src vm's backing file
		backingStore := srcDomCfg.Devices.Disks[mainDiskIndex].BackingStore
		if backingStore != nil && backingStore.Source.File != nil {
			vmBackingPath = backingStore.Source.File.File
		}

		// use src vm image as if src is a template without backing path
		if vmBackingPath == "" {
			vmBackingPath = srcDomCfg.Devices.Disks[mainDiskIndex].Source.File.File
		}
	}

	// gen xml definition
	vmXml := ""
	vmXml, vmRawPath, _ = s.QemuService.GenVmDef(tmplXml, vmMacAddress, vmUniqueName, vmBackingPath, vmCpu, vmMemorySize)
	if err != nil {
		_logUtils.Errorf("err gen vm xml, err %s", err.Error())
		return
	}

	err = s.QemuService.createDiskFile(vmBackingPath, vmRawPath, vmDiskSize)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	dom, err = s.LibvirtConn.DomainCreateXML(vmXml, libvirt.DOMAIN_NONE)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	newXml := ""
	newXml, err = dom.GetXMLDesc(0)
	if err != nil {
		return
	}

	newDomCfg := &libvirtxml.Domain{}
	err = newDomCfg.Unmarshal(newXml)

	vmMacAddress = newDomCfg.Devices.Interfaces[0].MAC.Address
	vmVncPort = newDomCfg.Devices.Graphics[0].VNC.Port

	return
}

func (s *LibvirtService) ListVm() (domains []libvirt.Domain) {
	if s.LibvirtConn == nil {
		return
	}

	domains, err := s.LibvirtConn.ListAllDomains(0)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	return
}

func (s *LibvirtService) GetVm(name string) (dom *libvirt.Domain, err error) {
	dom, err = s.LibvirtConn.LookupDomainByName(name)
	if err != nil {
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
func (s *LibvirtService) DestroyVmByName(name string, removeDiskImage bool) {
	dom, err := s.GetVm(name)
	if err != nil {
		return
	}

	dickPath, err := s.QemuService.GetDisk(dom)

	err = dom.Destroy()

	if removeDiskImage && dickPath != "" {
		_fileUtils.RmDir(dickPath)
	}

	return
}

func (s *LibvirtService) BootVmByName(name string) {
	dom, err := s.GetVm(name)
	if err != nil {
		return
	}

	err = dom.Reboot(libvirt.DOMAIN_REBOOT_DEFAULT)
	return
}
func (s *LibvirtService) RebootVmByName(name string) {
	dom, err := s.GetVm(name)
	if err != nil {
		return
	}

	err = dom.Reboot(libvirt.DOMAIN_REBOOT_DEFAULT)
	return
}
func (s *LibvirtService) ShutdownVmByName(name string) {
	dom, err := s.GetVm(name)
	if err != nil {
		return
	}

	err = dom.ShutdownFlags(libvirt.DOMAIN_SHUTDOWN_PARAVIRT)
	return
}

func (s *LibvirtService) SuspendVmByName(name string) {
	dom, err := s.GetVm(name)
	if err != nil {
		return
	}

	err = dom.Suspend()
	return
}
func (s *LibvirtService) ResumeVmByName(name string) {
	dom, err := s.GetVm(name)
	if err != nil {
		return
	}

	err = dom.Resume()
	return
}

func (s *LibvirtService) UndefineVm(dom *libvirt.Domain) (err error) {
	err = dom.Undefine()

	return
}

func (s *LibvirtService) GetVmDef(name string) (xml string) {
	dom, err := s.GetVm(name)
	if err != nil {
		return
	}

	xml, _ = dom.GetXMLDesc(0)

	return
}

func (s *LibvirtService) Connect(str string) {
	if s.LibvirtConn != nil {
		active, err := s.LibvirtConn.IsAlive()
		if err != nil {
			_logUtils.Errorf(err.Error())
		}
		if active {
			return
		}
	}

	var err error
	s.LibvirtConn, err = libvirt.NewConnect(str)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	active, err := s.LibvirtConn.IsAlive()
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}
	if !active {
		_logUtils.Errorf("not active")
	}
	return
}

func (s *LibvirtService) setVmProps(vm *domain.Vm) {
	osCategory := consts.Windows
	osType := consts.Win10
	osVersion := "x64-pro"
	osLang := consts.ZH_CN

	vm.Backing = fmt.Sprintf("%s/%s/%s-%s", osCategory.ToString(), osType.ToString(),
		osVersion, osLang.ToString())

	vm.Tmpl = fmt.Sprintf("tmpl-%s-%s-%s",
		osType.ToString(), osVersion, osLang.ToString())
	vm.Name = fmt.Sprintf("test-%s-%s-%s-%s",
		osType.ToString(), osVersion, osLang.ToString(), _stringUtils.NewUuid())
}
