# lowSpeedPortScan

一些有强防护的敏感目标使用nmap等上古神器是无法准确扫描出全部的资产端口。本项目属于DarkEye子项目，通过采用`端口校验`技术，解决此类问题，速度会慢点，时间～1天，但是～%100出结果。


项目特色
===
|特性 |状态|
|--------------------------|----------------|
|断点扫描 | 是|
|自动探测防火墙规则频率 | 是|
|限制IP自动回溯端口重新检测 | 是|

支持平台
===
|系统 |状态|
|--------------------------|----------------|
|MacOs | 是|
|Linux | 是|
|Windows | 是|




```golang
macdeMacBook-Pro:lowSpeedPortScan mac$ ./lowSpeedPortScan -h
___.   .__       _________         __   
\_ |__ |__| ____ \_   ___ \_____ _/  |_ 
 | __ \|  |/ ___\/    \  \/\__  \\   __\
 | \_\ \  / /_/  >     \____/ __ \|  |  
 |___  /__\___  / \______  (____  /__|  
     \/  /_____/         \/     \/
1.0.20209102143
https://github.com/zsdevX/DarkEye
大橘Oo0
84500316@qq.com
Examples: 
./lowSpeedPortScan -alive_port 8443 -ip f.u.c.k -port 1-65535 -rate_test -output result.txt
./lowSpeedPortScan -ip f.u.c.k,f.u.c.1-254 -port 1-65535 -rate 200 -output result.txt
----------------

Usage of ./lowSpeedPortScan:
  -alive_port string
        已知开放的端口用来校正扫描 (default "0")
  -ip string
        a.b.c.d（不做扫C，扫C自己想办法或使用nmap --scan-delay 1000ms但是不准 (default "127.0.0.1")
  -min_rate int
        自动计算的速率不能低于min_rate (default 100)
  -output string
        结果保存到该文件 (default "result.txt")
  -port string
        端口格式参考Nmap (default "1-65535")
  -rate int
        端口之间的扫描间隔单位ms，也可用通过-rate_test (default 2000)
  -rate_test
        发包频率
  -thread int
        结果保存到该文件 (default 1)
        
 ```
        





