# 部署宿主机的Multipass环境
这里使用Ubuntu20.04桌面版来作为宿主机，搭建虚拟化环境。其他系统的安装方式见
- mac https://multipass.run/docs/installing-on-macos
- windows https://multipass.run/docs/installing-on-windows
### 1.安装multipass
```
snap install multipass
```
### 2. 创建虚拟机
这里创建以名为ubuntu-1命名的虚拟机实例作为基础实例
```
multipass launch -n ubuntu-1
```