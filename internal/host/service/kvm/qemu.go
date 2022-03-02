package hostAgentKvmService

import "C"
import (
	"errors"
	"fmt"
	agentConf "github.com/easysoft/zv/internal/agent/conf"
	"github.com/easysoft/zv/internal/comm/const"
	"github.com/easysoft/zv/internal/comm/domain"
	_logUtils "github.com/easysoft/zv/internal/pkg/lib/log"
	_shellUtils "github.com/easysoft/zv/internal/pkg/lib/shell"
	_sshUtils "github.com/easysoft/zv/internal/pkg/lib/ssh"
	_stringUtils "github.com/easysoft/zv/internal/pkg/lib/string"
	"github.com/libvirt/libvirt-go"
	"github.com/libvirt/libvirt-go-xml"
	"golang.org/x/crypto/ssh"
	"path/filepath"
	"strings"
)

type QemuService struct {
}

func NewQemuService() *QemuService {
	return &QemuService{}
}

func (s *QemuService) GenVmDef(tmplXml, macAddress, vmName, backingPath string, vmCpu, vmMemory uint) (
	vmXml string, rawPath string, err error) {

	domCfg := &libvirtxml.Domain{}
	err = domCfg.Unmarshal(tmplXml)
	if err != nil {
		return
	}

	return s.GenVmDefFromCfg(domCfg, macAddress, vmName, backingPath, vmCpu, vmMemory)
}

func (s *QemuService) GenVmDefFromCfg(domCfg *libvirtxml.Domain, macAddress, vmName, backingPath string, vmCpu, vmMemory uint) (
	vmXml string, rawPath string, err error) {

	rawPath = filepath.Join(agentConf.Inst.DirImage, vmName+".qcow2")

	mainDiskIndex := s.GetMainDiskIndex(domCfg)

	domCfg.Name = vmName
	domCfg.UUID = _stringUtils.NewUuidWithSep()
	domCfg.Devices.Disks[mainDiskIndex].Source.File = &libvirtxml.DomainDiskSourceFile{
		File: rawPath,
	}

	if backingPath != "" {
		domCfg.Devices.Disks[mainDiskIndex].BackingStore = &libvirtxml.DomainDiskBackingStore{
			Index: 0,
			Format: &libvirtxml.DomainDiskFormat{
				Type: "qcow2",
			},
			Source: &libvirtxml.DomainDiskSource{
				File: &libvirtxml.DomainDiskSourceFile{
					File: backingPath,
				},
			},
		}
	}

	//<graphics type="vnc" port="-1" autoport="yes" listen="0.0.0.0" passwd="pass">
	//<listen type="address" address="0.0.0.0"/>
	//</graphics>
	domCfg.Devices.Graphics = []libvirtxml.DomainGraphic{
		{
			VNC: &libvirtxml.DomainGraphicVNC{
				AutoPort: "yes",
				Port:     -1,
				Listen:   "0.0.0.0",
				Passwd:   "pass",
				Listeners: []libvirtxml.DomainGraphicListener{
					{
						Address: &libvirtxml.DomainGraphicListenerAddress{
							Address: "0.0.0.0",
						},
					},
				},
			},
		},
	}

	domCfg.Devices.Controllers = s.removeUnnecessaryPciCtrl(domCfg)

	if vmCpu != 0 {
		domCfg.VCPU.Value = vmCpu
	}

	if vmMemory != 0 {
		domCfg.Memory = &libvirtxml.DomainMemory{
			Unit:  "M",
			Value: vmMemory,
		}
		domCfg.CurrentMemory = &libvirtxml.DomainCurrentMemory{
			Unit:  "M",
			Value: vmMemory,
		}
	}

	domCfg.Devices.Interfaces[0].MAC.Address = macAddress
	//domCfg.OS.Type.Machine = s.GenMachine()

	vmXml, _ = domCfg.Marshal()

	return
}

func (s *QemuService) GenVmDefTest(src, vmName, rawPath, backingPath string, vmCpu, vmMemory uint) (
	xml, macAddress string, err error) {

	domCfg := &libvirtxml.Domain{}
	err = domCfg.Unmarshal(src)
	if err != nil {
		return
	}

	mainDiskIndex := s.GetMainDiskIndex(domCfg)

	domCfg.Name = vmName
	domCfg.UUID = _stringUtils.NewUuidWithSep()
	domCfg.Devices.Disks[mainDiskIndex].Source.File = &libvirtxml.DomainDiskSourceFile{
		File: rawPath,
	}
	domCfg.Devices.Disks[mainDiskIndex].BackingStore = &libvirtxml.DomainDiskBackingStore{
		Index: 0,
		Format: &libvirtxml.DomainDiskFormat{
			Type: "qcow2",
		},
		Source: &libvirtxml.DomainDiskSource{
			File: &libvirtxml.DomainDiskSourceFile{
				File: backingPath,
			},
		},
	}

	//<graphics type="vnc" port="-1" autoport="yes" listen="0.0.0.0" passwd="pass">
	//<listen type="address" address="0.0.0.0"/>
	//</graphics>
	domCfg.Devices.Graphics = []libvirtxml.DomainGraphic{
		{
			VNC: &libvirtxml.DomainGraphicVNC{
				AutoPort: "yes",
				Port:     -1,
				Listen:   "0.0.0.0",
				Passwd:   "pass",
				Listeners: []libvirtxml.DomainGraphicListener{
					{
						Address: &libvirtxml.DomainGraphicListenerAddress{
							Address: "0.0.0.0",
						},
					},
				},
			},
		},
	}

	domCfg.Devices.Controllers = s.removeUnnecessaryPciCtrl(domCfg)

	if vmCpu != 0 {
		domCfg.VCPU.Value = vmCpu
	}

	if vmMemory != 0 {
		domCfg.Memory = &libvirtxml.DomainMemory{
			Unit:     "M",
			Value:    vmMemory,
			DumpCore: "yes",
		}
		domCfg.CurrentMemory = &libvirtxml.DomainCurrentMemory{
			Unit:  "M",
			Value: vmMemory,
		}
	}

	macAddress = s.GenMacAddress()
	domCfg.Devices.Interfaces[0].MAC.Address = macAddress

	//domCfg.OS.Type.Machine = s.GenMachine()

	xml, _ = domCfg.Marshal()

	return
}

func (s *QemuService) GenMacAddress() string {
	cmd := "dd if=/dev/urandom count=1 2>/dev/null | md5sum | sed 's/^\\(..\\)\\(..\\)\\(..\\)\\(..\\).*$/\\1:\\2:\\3:\\4/'"
	output, err := _shellUtils.ExeShell(cmd)
	if err != nil {
		_logUtils.Errorf("fail to exec cmd %s, err %s.", cmd, err.Error())
		return ""
	}

	mac := "fa:92:" + _stringUtils.TrimAll(output)
	return mac
}

func (s *QemuService) GenMachine() string {
	cmd := "qemu-system-x86_64 -M ? | grep default | awk '{print $1}'"
	output, err := _shellUtils.ExeShell(cmd)
	if err != nil {
		_logUtils.Errorf("fail to exec cmd %s, err %s.", cmd, err.Error())
		return ""
	}

	machine := _stringUtils.TrimAll(output)
	return machine
}

func (s *QemuService) GetBaseImagePath(vm domain.Vm) (path string) {
	dir := filepath.Join(agentConf.Inst.DirBaking, vm.OsCategory.ToString(), vm.OsType.ToString())
	name := fmt.Sprintf("%s-%s", vm.OsVersion, vm.OsLang.ToString())

	path = filepath.Join(dir, name)

	return
}

func (s *QemuService) createDiskFile(backingPath, diskPath string, diskSize uint) (err error) {
	var cmd string
	if backingPath == "" {
		cmd = fmt.Sprintf("qemu-img create -f qcow2 %s %dG",
			diskPath, diskSize/1000)
	} else {
		cmd = fmt.Sprintf("qemu-img create -f qcow2 -o cluster_size=2M,backing_file=%s %s %dG",
			backingPath, diskPath, diskSize/1000)
	}

	removeCmd := fmt.Sprintf("rm -rf %s", diskPath)
	_shellUtils.ExeShellInDir(removeCmd, agentConf.Inst.DirKvm)

	if agentConf.Inst.Host == "" { // local
		_, err = _shellUtils.ExeShellInDir(cmd, agentConf.Inst.DirKvm)
		if err != nil {
			msg := fmt.Sprintf("fail to create disk, cmd %s, err %s.", cmd, err.Error())
			_logUtils.Errorf(msg)
			err = errors.New(msg)
			return
		}

	} else { // remote
		var conn *ssh.Client
		conn, err = _sshUtils.Connect(agentConf.Inst.Host, agentConf.Inst.User)
		if err != nil {
			_logUtils.Errorf(err.Error())
			return
		}
		defer conn.Close()

		var session *ssh.Session
		session, err = conn.NewSession()
		if err != nil {
			_logUtils.Errorf(err.Error())
			return
		}
		defer session.Close()

		var cmdInfo []byte
		cmdInfo, err = session.CombinedOutput(cmd)
		if err != nil {
			_logUtils.Errorf(err.Error())
			return
		} else {
			_logUtils.Infof(string(cmdInfo))
		}
	}

	return
}

func (s *QemuService) setVmProps(vm *domain.Vm) {
	osCategory := consts.Windows
	osType := consts.Win10
	osVersion := "x64-pro"
	osLang := consts.ZH_CN

	vm.Backing = fmt.Sprintf("%s/%s/%s-%s", osCategory.ToString(), osType.ToString(),
		osVersion, osLang.ToString())
}

func (s *QemuService) GetMainDiskIndex(domCfg *libvirtxml.Domain) (ret int) {
	for index, item := range domCfg.Devices.Disks {
		if item.Device == "disk" && strings.Index(item.Source.File.File, "share") < 0 {
			ret = index
			return
		}
	}

	return
}

func (s *QemuService) getFirstPciCtrlIndex(domCfg *libvirtxml.Domain) (ret int) {
	for index, item := range domCfg.Devices.Controllers {
		if item.Type == "pci" && *item.Index == 0 {
			ret = index
			return
		}
	}

	return
}

func (s *QemuService) removeUnnecessaryPciCtrl(domCfg *libvirtxml.Domain) (ret []libvirtxml.DomainController) {
	for _, item := range domCfg.Devices.Controllers {
		if item.Type == "pci" && *item.Index != 0 {
			continue
		}
		//item.Model = "pci-root"
		ret = append(ret, item)
	}

	return
}

func (s *QemuService) GetDisk(dom *libvirt.Domain) (path string, err error) {
	var xml string

	xml, err = dom.GetXMLDesc(0)
	if err != nil {
		return
	}

	domCfg := &libvirtxml.Domain{}
	err = domCfg.Unmarshal(xml)
	if err != nil {
		return
	}

	mainDiskIndex := s.GetMainDiskIndex(domCfg)
	path = domCfg.Devices.Disks[mainDiskIndex].Source.File.File

	return
}
