按照[这篇文章](../deploy/4-zagent.md)，完成ZAgent服务端部署后，在根目录下会生成一个名为*server.db*的SQLite数据库文件。请使用SQLite客户端，打开并完成以下操作。

未来我们会提供Web网页来初始化这些业务数据，您也可以在对接的管理系统中自行实现。

- 将前面步骤安装完成的KVM宿主机，填入*biz_host*表中，内容参照如下记录。

| id   | name               | ip           | port | status | priority | capabilities | vendor |
| ---- | ------------------ | ------------ | ---- | ------ | -------- | ------------ | ------ |
| 1    | host_ 192.168.0.66 | 192.168.0.66 | 8086 | online | 100      | vm           | native |

