package commonService

import (
	v1 "github.com/easysoft/zv/cmd/host/router/v1"
	"github.com/easysoft/zv/internal/pkg/domain"
	"github.com/easysoft/zv/internal/server/model"
)

type RpcService struct {
}

func NewRpcService() *RpcService {
	return &RpcService{}
}

func (s RpcService) SeleniumTest(build model.Build) (result _domain.RpcResp) {
	buildTo := model.NewBuildTo(build)
	buildTo.BrowserType = build.BrowserType
	buildTo.BrowserVersion = build.BrowserVersion

	obj := interface{}(buildTo)
	result = s.Request(build.NodeIp, build.NodePort, "job", "Add", &obj)

	return
}

func (s RpcService) AppiumTest(build model.Build) (result _domain.RpcResp) {
	appiumTestTo := model.NewBuildTo(build)
	appiumTestTo.AppiumPort = build.AppiumPort

	obj := interface{}(appiumTestTo)
	result = s.Request(build.NodeIp, build.NodePort, "job", "Add", &obj)

	return
}

func (s RpcService) UnitTest(build model.Build) (result _domain.RpcResp) {
	buildTo := model.NewBuildTo(build)

	obj := interface{}(buildTo)
	result = s.Request(build.NodeIp, build.NodePort, "job", "Add", &obj)

	return
}

func (s RpcService) CreateKvm(hostIp string, hostPort int, req v1.KvmReq) (result _domain.RpcResp) {
	obj := interface{}(req)
	result = s.Request(hostIp, hostPort, "kvm", "Create", &obj)
	return
}
func (s RpcService) DestroyKvm(hostIp string, hostPort int, req v1.KvmReq) (result _domain.RpcResp) {
	obj := interface{}(req)
	result = s.Request(hostIp, hostPort, "kvm", "Destroy", &obj)
	return
}

func (s RpcService) CreateVmWare(hostIp string, hostPort int, req v1.VmWareReq) (result _domain.RpcResp) {
	obj := interface{}(req)
	result = s.Request(hostIp, hostPort, "vmware", "Create", &obj)
	return
}
func (s RpcService) DestroyVmWare(hostIp string, hostPort int, req v1.VmWareReq) (result _domain.RpcResp) {
	obj := interface{}(req)
	result = s.Request(hostIp, hostPort, "vmware", "Destroy", &obj)
	return
}

func (s RpcService) Request(ip string, port int, apiPath string, method string, param *interface{}) (
	rpcResult _domain.RpcResp) {

	//cc := &codec.MsgpackCodec{}
	//
	//data, _ := cc.Encode(param)
	//url := fmt.Sprintf("http://%s:%d/", ip, port)
	//req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	//if err != nil {
	//	msg := fmt.Sprintf("Fail to create request: %s", err.Error())
	//	_logUtils.Errorf(msg)
	//	rpcResult.Fail(msg)
	//	return
	//}
	//
	//// 设置header
	//h := req.Header
	//h.Set(gateway.XServicePath, apiPath)
	//h.Set(gateway.XServiceMethod, method)
	//h.Set(gateway.XMessageID, "10000")
	//h.Set(gateway.XMessageType, "0")
	//h.Set(gateway.XSerializeType, "3")
	//
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	msg := fmt.Sprintf("fail to call, err: %s", err.Error())
	//	_logUtils.Errorf(msg)
	//	rpcResult.Fail(msg)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	msg := fmt.Sprintf("fail to read response, err: %s", err.Error())
	//	_logUtils.Errorf(msg)
	//	rpcResult.Fail(msg)
	//
	//}
	//
	//err = cc.Decode(body, &rpcResult)
	//if err != nil {
	//	msg := fmt.Sprintf("fail to decode reply, err: %s", err.Error())
	//	_logUtils.Errorf(msg)
	//	rpcResult.Fail(msg)
	//}
	//
	//msg := fmt.Sprintf("agent return %d-%s.", rpcResult.Code, rpcResult.Msg)
	//_logUtils.Info(msg)

	return
}
