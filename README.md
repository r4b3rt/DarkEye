<h1 align="center">👏👏👏 欢迎使用 DarkEye 👏👏👏</h1>

![Go Report Card](https://img.shields.io/github/release-date/b1gcat/DarkEye) [![Go Report Card](https://goreportcard.com/badge/github.com/b1gcat/DarkEye)](https://goreportcard.com/report/github.com/b1gcat/DarkEye)


> 赶过年前架构调整

## Demo

## 🚀快速使用

### 1、主机发现
支持多种`loader: tcp、ping、http、nb`
```bash
./dist/superscan_darwin_amd64 -action disco-host -ip 192.168.1.1-254
```

### 2、网段发现

#### 指定探测协议

支持两种`loader: tcp、ping`

```bash
 ./dist/superscan_darwin_amd64 -action disco-net -loader ping -ip 192.168.1-254 
```

### 3、协议爆破
可查看帮助选取loader，默认为所有协议插件
```bash
./dist/superscan_darwin_amd64 -action risk -loader ssh -p 22  -ip 192.168.1.253 		
```

修改爆破字典

```bash
./dist/superscan_darwin_amd64 -action risk -ip 192.168.1.253 -p 22 -user varbing -pass pass.txt
```

### 4、IP域名碰撞

```bash
./dist/superscan_darwin_amd64 -action ip-host -ip 192.168.1.1-2 -p 80 -host host.txt
```

## ⚡️技巧

1. 查看帮助：`./dist/superscan_darwin_amd64 -h`。
2. 当IP数量多时，可以使用`-t 256`增加IP并发。
3. 当端口数量多时，可以使用`-tt 100`增加端口并发。
3. `-ip`参数可接：掩码：`a.b.c.d/24`、范围：`a.b.c.1-254`、子网范围 :`a.b.1-254`、IP:`a.b.c.d`


## 🛠 编译安装

```bash
git clone https://github.com/b1gcat/DarkEye.git
./build all

Tips:编译好后文件都自动发布到dist目录下
```

# 404StarLink 2.0 - Galaxy

![](https://github.com/knownsec/404StarLink-Project/raw/master/logo.png)

DarkEye 是 404Team [星链计划2.0](https://github.com/knownsec/404StarLink2.0-Galaxy)中的一环，如果对DarkEye
有任何疑问又或是想要找小伙伴交流，可以参考星链计划的加群方式。

- [https://github.com/knownsec/404StarLink2.0-Galaxy#community](https://github.com/knownsec/404StarLink2.0-Galaxy#community)



