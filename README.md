# DarkEye

DarkEye项目旨为收集目标信息，如实反馈结果，无主观判断。DarkEye从互联网（SecurityTrails、FoFa、google）收集情报或目标公开的接口信息，仅做汇总并不做分析，本项目不具备攻击性，漏洞的利用主要依靠人或其它工具来支撑。

功能介绍
===
### 超级扫描
```aidl
1、支持范围扫描（IP、端口）
2、常用协议弱口令爆破
3、支持获取标题和中间件
4、绕过防火墙频率限制扫描（限单IP）
5、收集结果自动保存为csv报告
```

### 信息搜集
```aidl
1、支持Fofa收集资产信息，无需Key。
2、支持SecurityTrails收集子域信息，并扩展支持提取域名解析的ip、cname、地域、标题。（50个域名/1key，多申请:)）
3、爬取网站（含js、html、xml、json等）贪婪搜索返回内容中任何位置可能存在的接口路径; 敏感路径分级;
4、支持google hack爬取数据，无需翻墙但是需要到https://serpstack.com/申请key（1key/100次/每月, 多申请:)）
5、支持长亭xray官方poc解析，poc文件可从下列列表白嫖：
    https://github.com/chaitin/xray/tree/master/pocs
    https://github.com/Laura0xiaoshizi/xray_pocs
```

功能截图
===
超级扫描
![avatar](screenshot/portscan.jpg)
主界面
![avatar](screenshot/darkeye.jpg)
POC扫描
![avatar](screenshot/pocscan.jpg)


支持平台
===
|系统 |状态|
|--------------------------|----------------|
|MacOs | 支持|
|Linux | 支持|
|Windows | 支持|


安装
===

##### QT环境

```qt
参考: https://github.com/therecipe/qt/wiki/Installation
```

##### Build Windows/macOS/Linux

```golnag
go get github.com/zsdevX/DarkEye
./build mac
./build linux
./build win
编译好后文件都自动发布到dist目录下

```

