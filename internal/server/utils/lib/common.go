package serverUitls

import (
	"fmt"
	"github.com/easysoft/zv/internal/comm/const"
	"strings"
)

func GenVmHostName(queueId uint,
	osPlatform consts.OsCategory, osName consts.OsType, osLang consts.OsLang) (ret string) {
	ret = fmt.Sprintf("queue%d-%s-%s-%s", queueId, osPlatform, osName, osLang)

	return
}

func IsVm(platform string) (result bool) {
	result = strings.Index(platform, consts.PlatformVm.ToString()) > -1
	return
}
func IsDocker(platform string) (result bool) {
	result = strings.Index(platform, consts.PlatformDocker.ToString()) > -1
	return
}
func IsNative(platform string) (result bool) {
	result = strings.Index(platform, consts.PlatformNative.ToString()) > -1
	return
}
func IsHuaweiCloud(platform string) (result bool) {
	result = strings.Index(platform, consts.HostVendorHuaweiCloud.ToString()) > -1
	return
}
func IsAliyun(platform string) (result bool) {
	result = strings.Index(platform, consts.HostVendorAliyun.ToString()) > -1
	return
}
func IsVirtualBox(platform string) (result bool) {
	result = strings.Index(platform, consts.HostVendorVirtualBox.ToString()) > -1
	return
}
func IsVmWare(platform string) (result bool) {
	result = strings.Index(platform, consts.HostVendorVmWare.ToString()) > -1
	return
}

func IsCloud(platform string) (result bool) {
	result = strings.Index(platform, consts.PlatformCloud.ToString()) > -1
	return
}
