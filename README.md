# DarkEye
```
本项目旨为从互联网收集目标信息，如实反馈结果。
信息主要来源：SecurityTrails、FoFa、google、网站API接口。
本项目不具备攻击性，漏洞的利用分析主要依靠人。
```
功能介绍
===
### 超级扫描
命令行模式运行
```aidl
* 支持范围扫描（IP、端口）。
* 支持活跃网段检测，即：从B网段中检测有活跃主机的C段。
* 常用协议弱口令爆破。
* 支持获取标题和中间件。
* 绕过防火墙频率限制扫描（限单IP）。
* 支持长亭xray官方poc解析，poc文件可从下列列表白嫖：
    https://github.com/chaitin/xray/tree/master/pocs
    https://github.com/Laura0xiaoshizi/xray_pocs
6、收集结果自动报告输出。
```

#### 帮助
```aidl
./supercan  -h
```

#### 口令+Vul爆破
```aidl
./supercan  -ip 192.168.1.1-192.168.255.255
```
#### 活跃网段检测
```aidl
./supercan  -ip 192.168.1.1-192.168.255.255 -only-check-alive
```

#### Poc测试
```aidl
cd superscan/util
go build poc.go
./poc -test -test-poc ../db_poc/shiro.yml -test-url http://www.baidu.com
```

### 信息搜集
图形界面模式运行
```aidl
* 支持Fofa收集资产信息，无需Key。
* 支持SecurityTrails收集子域信息，并扩展支持提取域名解析的ip、cname、地域、标题。（50个域名/1key，多申请:)）
* 爬取网站（含js、html、xml、json等）贪婪搜索返回内容中任何位置可能存在的接口路径; 敏感路径分级;
* 支持google hack爬取数据，无需翻墙但是需要到https://serpstack.com/申请key（1key/100次/每月, 多申请:)）
* 收集结果自动报告输出    
```

功能截图
===
超级扫描
![avatar](screenshot/superscan.jpg)
主界面
![avatar](screenshot/darkeye.jpg)


支持平台
===
全平台


编译安装
===

### 环境准备
#### QT环境
'信息搜集'的图形界面部分采用qt，需安装qt支持库。
```qt
参考: https://github.com/therecipe/qt/wiki/Installation
```

### 编译方法
## 信息搜集
```golnag
go get github.com/zsdevX/DarkEye
./build mac
./build linux
./build win
编译好后文件都自动发布到dist目录下
```

### 超级扫描
```golang
go get github.com/zsdevX/DarkEye
cd superscan
./build all
```

