# 在线信息收集
## 域名

### 域名所有者信息
 ```
https://whois.aliyun.com/ //阿里
https://www.whois365.com/cn/ //全球查
http://ipwhois.cnnic.net.cn/ //中国互联网信息中心
http://whois.domaintools.com/ //Whois Lookup, Domain Availability & IP Search - DomainTools
http://whois.chinaz.com/ //站长
https://x.threatbook.cn/ //微步
http://whois.aizhan.com //爱站
kali命令: whois www.xxx.com
http://tianyancha.com //天眼查
http://www.beianbeian.com/ //ICP备案查询网
http://beian.miit.gov.cn //工信部
http://www.gsxt.gov.cn/index.html //国家企业信用信息公示系统
```
### 子域名

#### 第三方dns数据库
```
curl --request GET --url "https://api.securitytrails.com/v1/history/${host}/dns/a" -H "apikey: ${you-api-key}" --header 'accept: application/json' 

https://viewdns.info/
https://securitytrails.com/
https://dnsdb.io/zh-cn/
https://dnsdumpster.com/
https://searchdns.netcraft.com/
https://x.threatbook.cn/en
Virustotal//输入”domain:target.com”,返回子域名列表和一些辅助信息。
```

#### 在线爆破
```
http://tool.chinaz.com/subdomain //站长工具
https://phpinfo.me/domain //在线子域名爆破工具
http://z.zcjun.com/ //在线子域名爆破-子成君提供
https://dns.bufferover.run/dns?q=,
```

#### 基于证书查域名
```
https://censys.io
https://crt.sh/
https://developers.facebook.com/tools/ct/
https://spyse.com/search/certificate
https://google.com/transparencyreport/https/ct/
```

### 在线端口爆破
```
http://tool.chinaz.com/port
http://www.shungg.cn/sm/
http://coolaf.com/tool/port
https://tool.lu/portscan/
http://tool.cc/port/
http://duankou.wlphp.com/
https://www.astrill.com/zh/port-scan
```

### 敏感信息
```
site:218.87.21.*

inurl:editor/db/ 
inurl:eWebEditor/db/ 
inurl:bbs/data/ 
inurl:databackup/ 
inurl:blog/data/ 
inurl:bokedata 
inurl:bbs/database/ 
inurl:conn.asp 
inc/conn.asp
Server.mapPath(“.mdb”)
allinurl:bbs data
filetype:mdb inurl:database
filetype:inc conn
inurl:data filetype:mdb
intitle:"index of" data
intitle:"index of" etc
intitle:"Index of" .sh_history
intitle:"Index of" .bash_history
intitle:"index of" passwd
intitle:"index of" people.lst
intitle:"index of" pwd.db
intitle:"index of" etc/shadow
intitle:"index of" spwd
intitle:"index of" master.passwd
intitle:"index of" htpasswd
inurl:service.pwd
```

### 搜索引擎
```aidl
FOFA
zoomEye
```
