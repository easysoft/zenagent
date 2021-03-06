package consts

type HostCapability string

const (
	PlatformVm     HostCapability = "vm"
	PlatformDocker HostCapability = "docker"

	PlatformNative HostCapability = "native"
	PlatformCloud  HostCapability = "cloud"
)

func (e HostCapability) ToString() string {
	return string(e)
}

type HostVendor string

const (
	HostVendorVirtualBox  HostVendor = "virtualbox"
	HostVendorVmWare      HostVendor = "vmware"
	HostVendorHuaweiCloud HostVendor = "huaweicloud"
	HostVendorAliyun      HostVendor = "aliyun"

	HostVendorPve       HostVendor = "pve"
	HostVendorPortainer HostVendor = "portainer"
)

func (e HostVendor) ToString() string {
	return string(e)
}

type HostStatus string

const (
	HostOnline  HostStatus = "online"
	HostOffline HostStatus = "offline"
	HostBusy    HostStatus = "busy" // report by agent on host
)

func (e HostStatus) ToString() string {
	return string(e)
}

type VmStatus string

const (
	VmCreated    VmStatus = "created"        // set first time created
	VmLaunch     VmStatus = "launch"         // set after success to call vm creating remotely
	VmFailCreate VmStatus = "vm_fail_create" // set after fail to call vm creating remotely

	VmRunning VmStatus = "running" // report by agent on host
	VmShutOff VmStatus = "shutoff" // report by agent on host

	VmBusy  VmStatus = "busy"  // report by agent in vm
	VmReady VmStatus = "ready" // report by agent in vm

	VmUnknown     VmStatus = "unknown"      // report by agent on host, not running, destroy and shutoff
	VmDestroy     VmStatus = "destroy"      // final status
	VmDestroyFail VmStatus = "destroy_fail" // set after fail to call vm destroy remotely
)

func (e VmStatus) ToString() string {
	return string(e)
}

type BuildProgress string

const (
	// start group
	ProgressCreated BuildProgress = "created"

	// res group
	ProgressResPending     BuildProgress = "res_pending"
	ProgressResLaunched    BuildProgress = "res_launched"
	ProgressResReady       BuildProgress = "res_ready"
	ProgressResFailed      BuildProgress = "res_failed"
	ProgressResDestroy     BuildProgress = "res_destroy"
	ProgressResFailDestroy BuildProgress = "res_fail_destroy"

	// exec group
	ProgressRunning BuildProgress = "running"
	ProgressRunFail BuildProgress = "run_fail"

	// end group
	ProgressCompleted BuildProgress = "completed"
	ProgressTimeout   BuildProgress = "timeout"
	ProgressTerminal  BuildProgress = "terminal"
	ProgressCancel    BuildProgress = "cancel"
)

func (e BuildProgress) ToString() string {
	return string(e)
}

type BuildStatus string

const (
	StatusCreated BuildStatus = "created"
	StatusPass    BuildStatus = "pass"
	StatusFail    BuildStatus = "fail"
)

func (e BuildStatus) ToString() string {
	return string(e)
}

type BuildType string

const (
	ZtfTest       BuildType = "ztf"
	SeleniumTest  BuildType = "selenium"
	AppiumTest    BuildType = "appium"
	UnitTest      BuildType = "unittest"
	InterfaceTest BuildType = "interface"
)

func (e BuildType) ToString() string {
	return string(e)
}

type DeviceStatus string

const (
	DeviceOff     DeviceStatus = "off"
	DeviceOn      DeviceStatus = "on"
	DeviceActive  DeviceStatus = "active"
	DeviceBusy    DeviceStatus = "busy"
	DeviceUnknown DeviceStatus = "unknown"
)

func (e DeviceStatus) ToString() string {
	return string(e)
}

type ServiceStatus string

const (
	ServiceOffline ServiceStatus = "offline"
	ServiceOnline  ServiceStatus = "online"
	ServiceReady   ServiceStatus = "ready"
	ServiceBusy    ServiceStatus = "busy"
)

func (e ServiceStatus) ToString() string {
	return string(e)
}

type OsDevice string

const (
	Android OsDevice = "android"
	Ios     OsDevice = "ios"
	Harmony OsDevice = "harmony"
)

func (e OsDevice) ToString() string {
	return string(e)
}

type OsCategory string

const (
	Windows OsCategory = "windows"
	Linux   OsCategory = "linux"
	Mac     OsCategory = "mac"
)

func (e OsCategory) ToString() string {
	return string(e)
}

type OsType string

const (
	Win10 OsType = "win10"
	Win7  OsType = "win7"
	WinXP OsType = "winxp"

	Ubuntu OsType = "ubuntu"
	CentOS OsType = "centos"
	Debian OsType = "debian"
)

func (e OsType) ToString() string {
	return string(e)
}

type OsLang string

const (
	EN_US OsLang = "en_us"
	ZH_CN OsLang = "zh_cn"
	ZH_TW OsLang = "zh_tw"
)

func (e OsLang) ToString() string {
	return string(e)
}

type BrowserType string

const (
	Chrome  BrowserType = "chrome"
	Firefox BrowserType = "firefox"
	Edge    BrowserType = "edge"
	IE      BrowserType = "ie"
)

func (e BrowserType) ToString() string {
	return string(e)
}

type AuthType string

const (
	None        AuthType = ""
	BasicAuth   AuthType = "basicAuth"
	BearerToken AuthType = "bearerToken"
	OAuth2      AuthType = "oauth2"
)

func (e AuthType) ToString() string {
	return string(e)
}

type OAuth2TypeGrantType string

const (
	AuthCode                         OAuth2TypeGrantType = "authCode"
	AuthCodeWithPKCE                 OAuth2TypeGrantType = "authCodeWithPKCE" // pkce: Proof Key for Code Exchange
	Implicit                         OAuth2TypeGrantType = "implicit"
	ResourceOwnerPasswordCredentials OAuth2TypeGrantType = "resourceOwnerPasswordCredentials"
	ClientCredentials                OAuth2TypeGrantType = "clientCredentials"
)

type OAuth2ClientAuthType string

const (
	AsBasicAuthHeader OAuth2ClientAuthType = "asBasicAuthHeader"
	CredentInBody     OAuth2ClientAuthType = "credentInBody"
)

func (e OAuth2ClientAuthType) ToString() string {
	return string(e)
}

type CodeChallengeMethod string

const (
	SHA256 CodeChallengeMethod = "sha256"
	Plain  CodeChallengeMethod = "plain"
)

func (e CodeChallengeMethod) ToString() string {
	return string(e)
}

type TestType string

const (
	Auto      TestType = "auto"
	Interface TestType = "interface"
	Case      TestType = "case"
	Scenario  TestType = "scenario"
)

func (e TestType) ToString() string {
	return string(e)
}

type PreviewType string

const (
	Json PreviewType = "json"
	Html PreviewType = "html"
	Xml  PreviewType = "xml"
	Text PreviewType = "text"
)

func (e PreviewType) ToString() string {
	return string(e)
}

type ProcessorType string

const (
	Simple    ProcessorType = "simple"
	DataLoop  ProcessorType = "data_loop"
	Extractor ProcessorType = "extractor"
)

func (e ProcessorType) ToString() string {
	return string(e)
}

type ErrorAction string

const (
	ActionContinue        ErrorAction = "continue"
	ActionStartNextThread ErrorAction = "start_next_thread"
	ActionLoop            ErrorAction = "loop"
	ActionStopThread      ErrorAction = "stop_thread"
	ActionStopTest        ErrorAction = "stop_test"
	ActionStopTestNow     ErrorAction = "stop_test_now"
)

func (e ErrorAction) ToString() string {
	return string(e)
}

type DataSource string

const (
	ZenData DataSource = "zendata"
	CSV     DataSource = "csv"
	Excel   DataSource = "excel"
)

func (e DataSource) ToString() string {
	return string(e)
}

type ExtractorType string

const (
	Value       ExtractorType = "value"
	XPath       ExtractorType = "xpath"
	JSONPath    ExtractorType = "json_path"
	CssSelector ExtractorType = "sss_selector"
)

func (e ExtractorType) ToString() string {
	return string(e)
}

type ExtractorSource string

const (
	Body   ExtractorSource = "body"
	Header ExtractorSource = "header"
)

func (e ExtractorSource) ToString() string {
	return string(e)
}

type EntityType string

const (
	Task  EntityType = "task"
	Queue EntityType = "queue"
	Build EntityType = "build"
	Vm    EntityType = "vm"
)

func (e EntityType) ToString() string {
	return string(e)
}
