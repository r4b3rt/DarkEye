# DarkEye

DarkEye项目旨为收集目标信息。DarkEye从互联网（SecurityTrails、fofa、github等）收集情报或目标公开的接口信息，仅做汇总并不做分析，本项目不具备攻击性，漏洞的利用主要依靠人或其它工具来支撑。

TODO LIST
===
```$xslt
* 敏感接口检测增加批量导入
* 敏感接口增加等级方便查看
* 敏感接口增加保存内容为csv
* 增加poc检测框架，爬页面时自动匹配
```

支持平台
===
|系统 |状态|
|--------------------------|----------------|
|MacOs | 支持|
|Linux | 支持|
|Windows | 支持|


支持功能
===
|功能 |描述|
|--------------------------|----------------|
|收集C段资产 | 输入IP后通过**FoFa**自动收集资产信息（**免key**），同时**判断收集资产有效性**，若有其它好途径请留言会增加;收集的结果会自动保存csv格式|
|收集子域 | 通过**SecurityTrails**收集子域名，并扩展支持提取**域名解析的ip、cname、地域、标题**; 需要使用key建议官网申请2-3个（50个域名/1key）;收集的结果自动保存为csv格式|
|敏感接口 | 爬取网站（含js、html、xml、json等）贪婪搜索返回内容中任何位置可能存在的接口路径|
|端口扫描 | 支持扫IP，IP范围扫描、**支持获取标题和中间件**、支持端口范围和指定端口扫描（默认为常用端口）、**绕过防火墙频率限制扫描（仅支持单IP）**;收集结果自动保存为csv格式|

![avatar](screenshot/portscan.png)

![avatar](screenshot/darkeye.png)


功能使用
===
|功能 |描述|
|--------------------------|----------------|
|收集C段资产| UI操作方式，直接运行**DarkEye**即可|
|收集子域| UI操作方式，直接运行**DarkEye**即可|
|敏感接口| UI操作方式，直接运行**DarkEye**即可|
|端口扫描 | 命令行运行： **portscan -h**可查看帮助|


安装
===

##### QT环境

```qt
参考: https://github.com/therecipe/qt
```

##### Build Windows/macOS/Linux

```golnag
go get github.com/zsdevX/DarkEye
./build mac
./build linux
./build win
编译好后文件都自动发布到dist目录下

```

