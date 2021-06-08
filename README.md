<h1 align="center">👏👏👏 欢迎使用 DarkEye 👏👏👏</h1>

![Go Report Card](https://img.shields.io/github/release-date/b1gcat/DarkEye) [![Go Report Card](https://goreportcard.com/badge/github.com/b1gcat/DarkEye)](https://goreportcard.com/report/github.com/b1gcat/DarkEye)


> Whatever!

## Demo

## 🚀快速使用 

### 弱口令、指纹发现

```bash
df -ip 192.168.1.1-192.168.1.2
扫描任务完成macdeMacBook-Pro:DarkEye mac$ ./dist/df_darwin_amd64 -ip 45.88.13.188 -website-domain-list www.hackdoor.org -port-list 80
INFO[0000] 已加载1个IP,共计1个端口,启动每IP扫描端口线程数128,同时可同时检测IP数量32 
INFO[0000] Plugins::netbios snmp postgres redis smb web memcached mssql mysql ftp mongodb ssh  

Cracking...              100% [==================================================================================================================================================] (1/1, 13 it/min) 
```

### 网段发现

```bash
fiware-wilma:DarkEye mac$ ./dist/df_darwin_amd64 -ip 192.168.1.0-192.168.255.0  -only-alive-network
当前为非管理权限模式，请切换管理员运行。
如果不具备管理员权限需要设置原生的命令（例如：ping）检测。请设置命令参数：
输入探测命令(default: ping -c 1 -W 1):
输入探测的成功关键字(default: , 0.0%)
输入命令shell环境(default: sh -c )

使用命令Shell环境'sh -c '
使用探测命令 'ping -c 1 -W 1'检查网络 
使用关键字' , 0.0%' 确定网络是否存在
192.168.1.0 is alive
192.168.2.0 is died
192.168.3.0 is died
192.168.4.0 is died
192.168.5.0 is died

```

### 主机发现

```bash
fiware-wilma:DarkEye mac$ ./dist/df_darwin_amd64 -ip 192.168.1.0-192.168.1.255 --alive-host-check
当前为非管理权限模式，请切换管理员运行。
如果不具备管理员权限需要设置原生的命令（例如：ping）检测。请设置命令参数：
输入探测命令(default: ping -c 1 -W 1):
输入探测的成功关键字(default: , 0.0%)
输入命令shell环境(default: sh -c )

使用命令Shell环境'sh -c '
使用探测命令 'ping -c 1 -W 1'检查网络 
使用关键字' , 0.0%' 确定网络是否存在
192.168.1.1 is alive
192.168.1.3 is alive
192.168.1.0 is died

```

### 主机碰撞

```bash
./dist/df_darwin_amd64 -ip 192.168.1.0-192.168.1.255 -website-domain website.txt
```

## 支持平台

```
Windows、Linux、MacOs、Arm、Mips[el]、FreeBsd ...
```


## 🛠 编译安装

```bash
git clone https://github.com/b1gcat/DarkEye.git
cd DarkEye
./build all

Tips:编译好后文件都自动发布到dist目录下
```

# 404StarLink 2.0 - Galaxy

![](https://github.com/knownsec/404StarLink-Project/raw/master/logo.png)

DarkEye 是 404Team [星链计划2.0](https://github.com/knownsec/404StarLink2.0-Galaxy)中的一环，如果对DarkEye 有任何疑问又或是想要找小伙伴交流，可以参考星链计划的加群方式。

- [https://github.com/knownsec/404StarLink2.0-Galaxy#community](https://github.com/knownsec/404StarLink2.0-Galaxy#community)


