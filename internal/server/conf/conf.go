package serverConf

import (
	"fmt"
	"github.com/easysoft/zv/internal/pkg/lib/common"
	"github.com/easysoft/zv/internal/pkg/lib/file"
	"github.com/easysoft/zv/internal/pkg/lib/log"
	"github.com/easysoft/zv/res/server"
	"github.com/jinzhu/configor"
	logger "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

var (
	Inst = Config{}
)

func Init() {
	exeDir := _fileUtils.GetExeDir()
	configPath := ""
	if _commonUtils.IsRelease() {
		configPath = filepath.Join(exeDir, "server.yml")
		if !_fileUtils.FileExist(configPath) {
			bytes, _ := serverRes.Asset("res/server/server.yml")
			_fileUtils.WriteFile(configPath, string(bytes))
		}
	} else {
		configPath = filepath.Join(exeDir, "cmd", "server", "server.yml")
	}

	_logUtils.Infof("从文件%s加载server配置", configPath)
	if err := configor.Load(&Inst, configPath); err != nil {
		logger.Println(fmt.Sprintf("Config Path:%s ,Error:%s", configPath, err.Error()))
		return
	}

	if Inst.Debug {
		fmt.Println(fmt.Sprintf("配置项：%+v", Inst))
	}
}

func GetRedisUris() []string {
	addrs := make([]string, 0, 0)
	hosts := strings.Split(Inst.Redis.Host, ";")
	ports := strings.Split(Inst.Redis.Port, ";")
	for _, h := range hosts {
		for _, p := range ports {
			addrs = append(addrs, fmt.Sprintf("%s:%s", h, p))
		}
	}
	return addrs
}

type Config struct {
	LogLevel string `yaml:"logLevel" env:"LogLevel" default:"info"`
	Debug    bool   `yaml:"debug" env:"Debug" default:"false"`
	//BinData  bool   `yaml:"binData" default:"false" env:"BinData"`
	Https    bool   `default:"false" env:"Https"`
	CertPath string `default:"" env:"CertPath"`
	CertKey  string `default:"" env:"CertKey"`
	Host     string `default:"0.0.0.0" env:"HostId"`
	Port     int    `default:"8085" env:"Port"`

	Admin struct {
		UserName        string `env:"AdminUserName" default:"admin"`
		Name            string `env:"AdminName" default:"admin"`
		Password        string `env:"AdminPassword" default:"P2ssw0rd"`
		RoleName        string `env:"AdminRoleName" default:"admin"`
		RoleDisplayName string `env:"RoleDisplayName" default:"超级管理员"`
	} `yaml:"admin,flow"`
	DB    DBConfig `yaml:"db,flow"`
	Redis struct {
		Enable bool   `env:"RedisDisable" default:"false"`
		Host   string `env:"RedisHost" default:"localhost"`
		Port   string `env:"RedisPort" default:"6379"`
		Pwd    string `env:"RedisPwd" default:""`
	} `yaml:"redis,flow"`

	Limit struct {
		Disable bool    `env:"LimitDisable" default:"true"`
		Limit   float64 `env:"LimitLimit" default:"1"`
		Burst   int     `env:"LimitBurst" default:"5"`
	}
	Qiniu struct {
		Enable    bool   `env:"QiniuEnable" default:"false"`
		Host      string `env:"QiniuHost" default:""`
		Accesskey string `env:"QiniuAccesskey" default:""`
		Secretkey string `env:"QiniuSecretkey" default:""`
		Bucket    string `env:"QiniuBucket" default:""`
	}
	Options struct {
		UploadMaxSize int64 `env:"uploadMaxSize" default:"100"`
	}
}

type DBConfig struct {
	Prefix   string `yaml:"prefix" env:"DBPrefix" default:"biz_"`
	Name     string `yaml:"name" env:"DBName" default:"ztest"`
	Adapter  string `yaml:"adapter" env:"DBAdapter" default:"sqlite3"`
	Host     string `yaml:"host" env:"DBHost" default:"localhost"`
	Port     string `yaml:"port" env:"DBPort" default:"3306"`
	User     string `yaml:"user" env:"DBUser" default:"root"`
	Password string `yaml:"password" env:"DBPassword" default:"P2ssw0rd"`
}
