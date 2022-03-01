package vmwareService

import (
	"bytes"
	"encoding/json"
	"fmt"
	v1 "github.com/easysoft/zv/cmd/host/router/v1"
	"log"
)

// GetAllVMs Method return array of Vm and a error variable if occurr some problem
// Return: []Vm and error
func (c *Client) GetAllVMs() ([]Vm, error) {
	var vms []Vm
	responseBody, err := c.httpRequest("api/vms", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Obj:%#v\n", responseBody)
	err = json.NewDecoder(responseBody).Decode(&vms)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Message: I can't read the json structure %s", err)
		return nil, err
	}

	for vm, value := range vms {
		name, desc, err := GetNameDescription(vms[vm].Path)
		if err != nil {
			return nil, err
		}
		vms[vm].Denomination = name
		vms[vm].Description = desc

		responseBody, err := c.httpRequest("api/vms/"+value.IdVM, "GET", bytes.Buffer{})
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(responseBody).Decode(&vms[vm])
		if err != nil {
			return nil, err
		}
		responseBody, err = c.httpRequest("api/vms/"+value.IdVM+"/power", "GET", bytes.Buffer{})
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(responseBody).Decode(&vms[vm])
		if err != nil {
			return nil, err
		}
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Obj:%#v\n", vms)
	return vms, nil
}

// CreateVM method to create a new VM in VmWare Worstation Input:
// tmplId: string with the ID of the origin VM,
//name: string with the denomination of the VM, desc: string with the description of VM
func (c *Client) CreateVM(tmplId string, name string, desc string) (*Vm, error) {
	var vm Vm
	requestBody := new(bytes.Buffer)
	request, err := json.Marshal(map[string]string{
		"name":     name,
		"parentId": tmplId,
	})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request %#v\n", request)
	requestBody.Write(request)
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	response, err := c.httpRequest("api/vms", "POST", *requestBody)
	if err != nil {
		return nil, err
	}
	// Piensa si tiene que ser en este punto en la parte de httpRequest
	// tienes que poner en este punto un control de errores de lo que responde VMW
	// si es diferente de create, ok, o delete que de un error y ponga a nil la vm y salga
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Response Error %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body %#v\n", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		return nil, err
	}

	// The following code we will use in the future when the VmWare fix it the method configparams
	// request, err = json.Marshal(map[string]string{
	// 	"name":  "annotation",
	// 	"value": d,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request %#v\n", request)
	// requestBody.Reset()
	// requestBody.Write(request)
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	// response, err = c.httpRequest("api/vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	// responseBody.Reset()
	// _, err = responseBody.ReadFrom(response)
	// if err != nil {
	// 	log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Response Error %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body %#v\n", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	// if err != nil {
	// 	return nil, err
	// }
	// Falta hacer un PUT para modificar los parametros de la instancia nueva. entre ellos el procesador la memoria y la network
	return &vm, err
}

// ReadVM method return the object Vm with the ID indicate in i.
// Input: i: string with the ID of the VM, Return: pointer at the Vm object
// and error variable with the error if occurr
func (c *Client) ReadVM(i string) (*Vm, error) {
	var vms []Vm
	var vm Vm
	// If you want see the path of the VM it's necessary getting all VMs
	// because the API of VmWare Workstation doesn't permit see this the another way
	response, err := c.httpRequest("api/vms", "GET", bytes.Buffer{})
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	err = json.NewDecoder(response).Decode(&vms)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Message: I can't read the json structure %s", err)
		return nil, err
	}
	for tempvm, value := range vms {
		if value.IdVM == i {
			vm = vms[tempvm]
		}
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:VM %#v\n", vm)
	return &vm, nil
}

func (c *Client) UpdateVM(id, name, desc string, processors, memory uint) (*Vm, error) {
	var buffer bytes.Buffer

	request, err := json.Marshal(map[string]int{
		"processors": int(processors),
		"memory":     int(memory),
	})
	if err != nil {
		return nil, err
	}
	buffer.Write(request)

	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: Request Body %#v\n", buffer.String())
	_, err = c.httpRequest("api/vms/"+id, "PUT", buffer)
	if err != nil {
		return nil, err
	}

	vm, err := GetVM(c, id)
	if err != nil {
		return nil, err
	}

	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: VM before %#v\n", vm)

	if name != "" || desc != "" {
		err = SetNameDescription(vm.Path, name, desc)
		if err != nil {
			return nil, err
		}

		vm, err = GetVM(c, id)
		if err != nil {
			return nil, err
		}

		log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: VM after %#v\n", vm)
	}

	return vm, err
}

// RegisterVM method to register a new VM in VmWare Worstation GUI:
// n: string with the VM NAME, p: string with the path of the VM
func (c *Client) RegisterVM(n string, p string) (*Vm, error) {
	var vm Vm
	requestBody := new(bytes.Buffer)
	request, err := json.Marshal(map[string]string{
		"name": n,
		"path": p,
	})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:Request %#v\n", request)
	requestBody.Write(request)
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:Request Body %#v\n", requestBody.String())
	response, err := c.httpRequest("api/vms/registration", "POST", *requestBody)
	if err != nil {
		return nil, err
	}
	// Piensa si tiene que ser en este punto en la parte de httpRequest
	// tienes que poner en este punto un control de errores de lo que responde VMW
	// si es diferente de create, ok, o delete que de un error y ponga a nil la vm y salga
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:response raw %#v\n", response)
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: RegisterVM Obj:Response Error %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:Response Body %#v\n", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		return nil, err
	}

	// Falta hacer un PUT para modificar los parametros de la instancia nueva. entre ellos el procesador la memoria y la network
	return &vm, err
}

func (c *Client) GetVmNic(id string) (nic *Nic, err error) {
	response, err := c.httpRequest(fmt.Sprintf("api/vms/%s/nic", id), "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: GetNicGetVmNic Obj:%#v\n", err)
		return
	}

	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: GetVmNic Obj:%#v, %#v\n", err, responseBody.String())
		return
	}

	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetVmNic Obj:%#v\n", responseBody)

	resp := NicResp{}
	err = json.NewDecoder(responseBody).Decode(&resp)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: GetVmNic Message: I can't read the json structure %s", err)
		return
	}

	if len(resp.Nics) > 0 {
		nic = &resp.Nics[0]
	}

	return
}

func (c *Client) SetRes(id string, processors, memory int) (err error) {
	req := v1.VmWareParam{
		Processors: processors,
		Memory:     memory,
	}

	request, err := json.Marshal(req)
	requestBody := new(bytes.Buffer)
	requestBody.Write(request)

	response, err := c.httpRequest("api/vms/"+id, "PUT", *requestBody)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: SetPower Obj:%#v\n", err)
		return err
	}

	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: SetPower Obj:%#v, %#v\n", err, responseBody.String())
		return err
	}

	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: SetPower Obj:%#v\n", responseBody)
	return nil

	return
}

func (c *Client) DestroyVM(id string) error {
	response, err := c.httpRequest("api/vms/"+id, "DELETE", bytes.Buffer{})
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: DestroyVM Obj:%#v\n", err)
		return err
	}
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: DestroyVM Obj:%#v, %#v\n", err, responseBody.String())
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: DestroyVM Obj:%#v\n", responseBody)
	return nil
}

func (c *Client) PowerOn(id string) error {
	return c.SetPower(id, On)
}
func (c *Client) PowerOff(id string) error {
	return c.SetPower(id, Off)
}
func (c *Client) ShutDown(id string) error {
	return c.SetPower(id, Shutdown)
}
func (c *Client) SetPower(id string, status VmWareStatus) error {
	requestBody := new(bytes.Buffer)
	requestBody.Write([]byte(status))

	response, err := c.httpRequest("api/vms/"+id+"/power", "PUT", *requestBody)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: SetPower Obj:%#v\n", err)
		return err
	}

	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: SetPower Obj:%#v, %#v\n", err, responseBody.String())
		return err
	}

	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: SetPower Obj:%#v\n", responseBody)
	return nil
}
