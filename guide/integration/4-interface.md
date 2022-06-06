## 宿主机注册接口实现

#### 请求服务器，管理测试任务

- [POST 创建测试任务](https://zagent-server.zentao.net/docs#/paths/~1api~1v1~1client~1task~1create/post)
- [GET 列出测试任务](https://zagent-server.zentao.net/docs#/paths/~1api~1v1~1client~1task~1list/get)
- [GET 获取测试任务](https://zagent-server.zentao.net/docs#/paths/~1api~1v1~1client~1task~1{id}/get)
- [PUT 更新测试任务](https://zagent-server.zentao.net/docs#/paths/~1api~1v1~1client~1task~1{id}/put)
- [DEL 删除测试任务](https://zagent-server.zentao.net/docs#/paths/~1api~1v1~1client~1task~1{id}/delete)

#### 请求服务器，注册资源

- [POST 注册宿主机](https://zagent-server.zentao.net/docs#/paths/~1api~1v1~1client~1host~1register/post)
- [POST 注册虚拟机](https://zagent-server.zentao.net/docs#/paths/~1api~1v1~1client~1vm~1register/post)

---

注意：以下接口，仅在您自行实现服务端，用以代替ZAgent的Server时，需要考虑。

#### 调用宿主机代理，管理虚拟机

- [POST 创建KVM虚拟机](https://zagent-host.zentao.net/docs#/paths/~1api~1v1~1kvm~1create/post)
- [POST 重启KVM虚拟机](https://zagent-host.zentao.net/docs#/paths/~1api~1v1~1kvm~1{name}~1reboot/post)
- [POST 暂停KVM虚拟机](https://zagent-host.zentao.net/docs#/paths/~1api~1v1~1kvm~1{name}~1suspend/post)
- [POST 恢复KVM虚拟机](https://zagent-host.zentao.net/docs#/paths/~1api~1v1~1kvm~1{name}~1resume/post)

#### 调用虚拟机代理，执行测试任务

- [POST 创建任务](https://zagent-vm.zentao.net/docs#/paths/~1api~1v1~1vmware~1create/post)

