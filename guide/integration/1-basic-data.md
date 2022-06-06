## ZAgent基础数据准备

在需要对接的第三方系统中，新建包含以下记录的**数据表**；您也可以在代码中使用**枚举**，来同ZAgent的基础数据定义保持一致。

这些数据将被用于双方接口的访问中，请保持其*编码*值始终一致。

### 系统大类：

| id   | 名称    | 编码    |
| ---- | ------- | ------- |
| 1    | Windows | windows |
| 2    | Linux   | linux   |
| 3    | Mac     | mac     |

### 系统分类：

| id   | 名称   | 编码   |
| ---- | ------ | ------ |
| 1    | Win10  | win10  |
| 2    | Win7   | win7   |
| 3    | WinXP  | winxp  |
| 4    | Ubuntu | ubuntu |
| 5    | CentOS | centos |
| 6    | Debian | debian |

### 系统语言：

| id   | 名称  | 编码  |
| ---- | ----- | ----- |
| 1    | EN_US | en_us |
| 2    | ZH_CN | zh_cn |
| 3    | ZH_TW | zh_tw |

### 浏览器表：

| id   | 名称    | 编码    |
| ---- | ------- | ------- |
| 1    | Chrome  | chrome  |
| 2    | Firefox | firefox |
| 3    | Edge    | edge    |
| 4    | IE      | ie      |