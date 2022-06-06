# 安装KVM虚拟机

首先，我们来了解一下KVM虚拟机有关的几个概念和工具。

- kvm：基于内核的虚拟机（引擎）
- qemu：用于模拟虚拟机IO设备
- qemu-img：虚拟机磁盘管理工具
- libvirt：虚拟化服务的API接口
- virsh：基于libvirt实现的命令行工具
- qemu-manager：图形化管理工具

新建KVM虚拟机时，可以指定另一磁盘文件作为 *BackingFile*。BackingFile是一个**只读**的虚拟磁盘基础映像，可以在多个虚拟机间进行共享。基于BackingFile创建和运行虚拟机时，只会在自己的磁盘文件中增量地写入文件，从而提高效率、节省磁盘和维护成本。

**虚拟机快照**保存了虚拟机在某个指定时间点的状态，当我们在自动化测试过程中遭遇问题或错误时，可以利用快照保存、并恢复到执行中的某个时间点。借助BackingFile机制，虚拟机支持形如以下的多层依赖的快照链（<--表示依赖的方向）。

```brainfuck
base image <-- vm01 <-- snap 1 <-- snap 2 <-- vm02(active)
```

可使用以下命令，将处于快照链中的某个虚机，导出形成一个独立的磁盘映像文件，其不再依赖其它映像。

```stylus
qemu-img convert -O qcow2 vm02.qcow2 vm-templ.img
```

假设我们在用户的工作目录中，建立了以下目录。

```crmsh
kvm          根目录
   iso       存放光盘镜像
   base      存放基础磁盘镜像BackingFile
   share     存放共享磁盘镜像，用户存储测试工具、驱动等
   image     存放测试机的磁盘镜像 
   xml       存放导出的虚拟机XML配置文件
```

##### 1.  下面用一个例子，给大家介绍下快速创建测试虚拟机的方法。

- 按照 [上一篇文章](1-host.md)中的步骤，创建一个Win10虚拟机；
- 在虚拟机中，安装好Chrome浏览器和测试相关软件；
- 使用以下命令，新建一个用于保存工具的共享磁盘；

   ```apache
   qemu-img create -f qcow2 -o cluster_size=2M kvm/share/tools.qcow2 10G
   ```

- 挂载共享磁盘到虚拟机,复制工具和文件到该共享磁盘；
- 移除该虚拟机，注意在确认对话框中，请选择“不删除相关磁盘文件”；
- 移动虚机主磁盘文件到基础镜像目录，如kvm/base/windows/win10。
- 执行以下命令，以上述基础镜像作为BackingFile，创建新的虚拟机磁盘；

   ```
   qemu-img create -f qcow2 -o cluster_size=2M,backing_file=kvm/base/windows/win10/pro-x64-zh_cn.qcow2 kvm/image/test-win10-pro-x64-zh_cn-01.qcow2 50G
   ```

- 打开qemu-manager，新建测试虚拟机，挂在新创建的虚拟机磁盘和共享磁盘。

##### 2.  除了使用图形界面的qemu-manager外，这里也提供一种命令行的方法，可用于代码中。

- 导出虚拟机XML配置文件；

   ```
   virsh dumpxml test-win10-pro-x64-zh_cn > kvm/xml/test-win10-pro-x64-zh_cn.xml
   ```

- 修改XML配置文件中的以下字段：

  - name
  - uuid
  - vcpu
  - memory和currentMemory
  - mac address
  - 第1块disk的source file

- 在第1块disk的Element中，加入以下BackingFile有关的内容：

   ```
   <backingStore type="file" index="2">
      <format type="qcow2"/>
      <source file="/home/aaron/kvm/base/windows/win10/pro-x64-zh_cn.qcow2"/> 
   <backingStore/>
   ```

- 如需要用页面VNC访问虚拟机桌面，找到XML的graphics元素，修改成以下内容；

   ```abnf
   <graphics type="vnc" port="-1" autoport="yes" listen="0.0.0.0" passwd="P2ssw0rd">
     <listen type="address" address="0.0.0.0"/>
   </graphics>
   ```

- 使用以下命令定义虚拟机；

   ```
   virsh define kvm /xml/test-win10-pro-x64-zh_cn.xml
   ```

- 使用以下命令启动虚拟机；

   ```powershell
   virsh start test-win10-pro-x64-zh_cn
   ```

- 使用以下命令获取虚拟机的VNC端口编号，在VNC软件中使用”5900+该数字“的端口，访问虚拟机远程桌面；

   ```stata
   virsh vncdisplay test-win10-pro-x64-zh_cn
   ```

3. ##### 最后，我们介绍推荐的第三种方法，在代码中使用libvirt*接口来管理虚拟机。禅道基于GO语言的实现，请见[这里](https://github.com/easysoft/zagent/blob/main/internal/agent-host/service/kvm/libvirt.go) 。

#### **虚拟机有关常用命令：**

```js
# 查看虚拟机信息
qemu-img info --backing-chain kvm/image/test-win10-pro-x64-zh_cn-01.qcow2

# 修改虚拟机磁盘大小
qemu-img resize x64-pro-zh_cn.qcow2 +10G

# 查看虚拟机里列表
virsh list --all

# 查看虚拟机VNC端口
virsh vncdisplay win10-test

# 导出虚拟机XML配置文件
virsh dumpxml win10-test > win10-test.xml

# 创建虚拟机磁盘镜像
qemu-img create -f qcow2 -o cluster_size=2M,backing_file=base.qcow2 win10-test.qcow2 40G

# 转换虚拟机镜像
qemu-img convert -O qcow2 vm02.qcow2 vm-templ.img

# 定义、取消定义，启动、停止虚拟机
virsh define win10-test.xml
virsh start win10-test
virsh destroy win10-test
virsh undefine win10-test
```

