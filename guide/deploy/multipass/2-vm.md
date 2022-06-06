# 准备multipass虚拟机基础镜像
### 进入虚拟机
```
multipass shell ubunut-1
```
### 1. ubuntu换源
```
sudo cp -ra /etc/apt/sources.list /etc/apt/sources.list.bak
sudo vim /etc/apt/sources.list
```
##### 用以下内容替换
```
deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse
    deb-src http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse

deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse

deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse

deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse

deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
```
##### 更新缓存和升级
```
sudo apt-get update
sudo apt-get upgrade
```

### 2. 安装桌面、协议
```
sudo apt install ubuntu-desktop xrdp
```
#### 设置用户帐号密码，便于后续的桌面访问
```
sudo passwd ubuntu
```
### 3. 安装配置 vncserver
##### 安装vnc服务
```
sudo apt-get install tightvncserver
```

##### 启动vnc服务（注意：不启动的话，下面所需的配置文件并不会生成）
```
vncserver
```
##### 备份vncserver 原配置文件
```
sudo cp  ~/.vnc/xstartup ~/.vnc/xstartup.bak
```

##### 编辑配置文件
```
sudo vim ~/.vnc/xstartup
```

##### 用以下内容替换
```
#!/bin/sh                                                                       

unset SESSION_MANAGER
unset DBUS_SESSION_BUS_ADDRESS
export XKL_XMODMAP_DISABLE=1
export XDG_CURRENT_DESKTOP="GNOME-Flashback:GNOME"
export XDG_MENU_PREFIX="gnome-flashback-"
[ -x /etc/vnc/xstartup ] && exec /etc/vnc/xstartup
[ -r $HOME/.Xresources ] && xrdb $HOME/.Xresources
xsetroot -solid grey
vncconfig -iconic &
#gnome-terminal &    
#nautilus &   
gnome-session --session=gnome-flashback-metacity --disable-acceleration-check &
```
##### 重新启动 vncserver
```
vncserver
```
##### 灰屛解决
```
sudo chmod +x ~/.vnc/xstartup
sudo apt-get install gnome-panel gnome-settings-daemon metacity nautilus gnome-terminal
```

##### 设置vncserver开机启动
```
sudo vim  /etc/init.d/tightvncserver
```
##### 将以下脚本中export USERR=‘ubunut'的 ubuntu 换成你自己的用户名,复制到tightvncserver
```
#!/bin/sh
### BEGIN INIT INFO
# Provides:          tightvncserver
# Required-Start:    $local_fs
# Required-Stop:     $local_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Start/stop tightvncserver
### END INIT INFO

# More details see:
# http://www.penguintutor.com/linux/tightvnc

### Customize this entry
# Set the USER variable to the name of the user to start tightvncserver under
export USER='ubuntu'
### End customization required

eval cd ~$USER

case "$1" in
  start)
    # 启动命令行。此处自定义分辨率、控制台号码或其它参数。
    su $USER -c '/usr/bin/tightvncserver -depth 16 -geometry 1928x1080 -dpi 100 :1'
    echo "Starting TightVNC server for $USER "
    ;;
  stop)
    # 终止命令行。此处控制台号码与启动一致。
    su $USER -c '/usr/bin/tightvncserver -kill :1'
    echo "Tightvncserver stopped"
    ;;
  *)
    echo "Usage: /etc/init.d/tightvncserver {start|stop}"
    exit 1
    ;;
esac
exit 0

```
##### 赋予权限是脚本生效
```
cd /etc/init.d
sudo chmod +x tightvncserver
sudo update-rc.d tightvncserver defaults
```
### 4. 安装好测试工具和文件

### 5. 复制镜像
镜像位置在/var/snap/multipass/common/data/multipassd/vault/instances/ubuntu-1
复制到 img目录下