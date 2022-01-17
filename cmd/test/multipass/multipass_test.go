package multipass

import (
	v1 "github.com/easysoft/zv/cmd/agent-host/router/v1"
	multiPassService "github.com/easysoft/zv/internal/agent-host/service/multipass"
	"log"
	"testing"
)

func TestMultiPass(t *testing.T) {
	srv := multiPassService.MultiPassService{}
	req := v1.MultiPassReq{
		VmUniqueName: "test1",
		VmMemory:     4096,
		Disk:         10,
		Cpus:         1,
		FilePath:     "ubuntu.img",
		ImgFrom:      "file",
	}

	vm, _ := srv.CreateVm(&req, false)
	log.Printf("%#v", vm)

	srv.RebootVmByName("test1")

	srv.SuspendVmByName("test1")

	srv.ResumeVmByName("test1")

	srv.DestroyVm("test1")

}
