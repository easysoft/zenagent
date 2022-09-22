package agentConf

import (
	consts "github.com/easysoft/zv/internal/pkg/const"
	_const "github.com/easysoft/zv/pkg/const"
	_commonUtils "github.com/easysoft/zv/pkg/lib/common"
	_fileUtils "github.com/easysoft/zv/pkg/lib/file"
	_httpUtils "github.com/easysoft/zv/pkg/lib/http"
	_i118Utils "github.com/easysoft/zv/pkg/lib/i118"
	"os/user"
	"path/filepath"
)

var (
	Inst = Config{}
)

func Init(app string) {
	_const.IsRelease = _commonUtils.IsRelease()
	if Inst.Language == "" {
		Inst.Language = "zh"
	}
	_i118Utils.InitI118(Inst.Language, app)

	Inst.Server = _httpUtils.AddUrlPostFixIfNeeded(Inst.Server)

	ip, macObj := _commonUtils.GetIp()
	Inst.MacAddress = macObj.String()
	if Inst.NodeIp == "" {
		Inst.NodeIp = ip.String()
	}

	usr, _ := user.Current()
	home := usr.HomeDir
	if Inst.Host != "" {
		home = "/home/" + Inst.User
	}

	Inst.WorkDir = _fileUtils.AddPathSepIfNeeded(filepath.Join(home, consts.AppNameAgent))

	if Inst.RunMode == consts.RunModeHost {
		Inst.DirKvm = _fileUtils.AddPathSepIfNeeded(filepath.Join(home, consts.FolderKvm))
		Inst.DirIso = _fileUtils.AddPathSepIfNeeded(filepath.Join(Inst.DirKvm, consts.FolderIso))
		Inst.DirBaking = _fileUtils.AddPathSepIfNeeded(filepath.Join(Inst.DirKvm, consts.FolderBacking))
		Inst.DirImage = _fileUtils.AddPathSepIfNeeded(filepath.Join(Inst.DirKvm, consts.FolderImage))
		Inst.DirToken = _fileUtils.AddPathSepIfNeeded(filepath.Join(Inst.DirKvm, consts.FolderToken))

		_fileUtils.MkDirIfNeeded(Inst.DirIso)
		_fileUtils.MkDirIfNeeded(Inst.DirBaking)
		_fileUtils.MkDirIfNeeded(Inst.DirImage)
	}
}

type Config struct {
	// fot libvirt testing only
	Host string
	User string

	RunMode    consts.RunMode `yaml:"runMode"`
	Server     string         `yaml:"Server"`
	NodeIp     string         `yaml:"ip"`
	NodePort   int            `yaml:"port"`
	MacAddress string

	Secret   string
	Language string
	NodeName string
	WorkDir  string
	LogDir   string

	DirKvm    string
	DirIso    string
	DirBaking string
	DirImage  string
	DirToken  string

	DB DBConfig
}

type DBConfig struct {
	Prefix string
}