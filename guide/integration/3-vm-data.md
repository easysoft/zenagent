按照[这篇文章](../deploy/3-zagent.md)，完成ZAgent服务端部署后，在根目录下，会生成一个名为server.db的SQLite数据库文件。请使用SQLite客户端，打开并完成以下操作。

未来我们会提供Web网页来初始化这些业务数据，您也可以在对接的管理系统中自行实现。

- 使用virt-manager，将安装好的虚拟机命名为 *win10-x64-pro-x64-zh_cn*；
- 移动虚拟机磁盘从目录*kvm/image*，到*kvm/base*；
- 在*biz_tmpl*表中，填入以下内容：

| id   | name                    | os_category | os_type | os_lang | os_version | host_id |
| ---- | ----------------------- | ----------- | ------- | ------- | ---------- | ------- |
| 1    | win10-x64-pro-x64-zh_cn | windows     | win10   | zh_cn   | pro-x64    | 1       |

- 使用virt-manager，查看该虚拟机的磁盘文件路径并记录为path，在*biz_baking*表中填写以下内容：

| id   | name                | path                                   | os_category | os_type | os_lang | os_version | suggest_cpu_count | suggest_memory_size | suggest_disk_size |
| ---- | ------------------- | -------------------------------------- | ----------- | ------- | ------- | ---------- | ----------------- | ------------------- | ----------------- |
| 1    | win10-pro-x64-zh_cn | base/windows/win10/pro-x64-zh_cn.qcow2 | windows     | win10   | zh_cn   | pro-x64    | 2                 | 4000                | 50000             |

- 在*biz_backing_browser_r*表中设置*biz_baking*记录到*biz_browser*的关联关系，用以定义”该基础磁盘镜像文件中，安装了哪些浏览器“；
- 在*biz_host_backing_r*表中填入*biz_host*记录到*biz_baking*的关联关系 ，用以定义”该宿主机上，包涵了哪些基础磁盘镜像“。