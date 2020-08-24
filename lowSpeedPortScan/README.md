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
macbook:lowSpeedPortScan mac$ ./lowSpeedPortScan -h
Usage of ./lowSpeedPortScan:
  -alive_port string
        已知开放的端口用来校正扫描 (default "80")
  -examples
        显示使用示例
  -ip string
        a.b.c.d（不做扫C，扫C自己想办法或使用nmap --scan-delay 1000ms, 不准 (default "127.0.0.1")
  -min_speed int
        自动计算的速率不能低于min_speed (default 100)
  -output string
        结果保存到该文件 (default "result.txt")
  -port string
        端口格式参考Nmap (default "1-65535")
  -speed int
        端口之间的扫描间隔单位ms，也可用通过-test_speed自动计算 (default 2000)
  -speed_test
        检测防火墙限制频率
 ```
        





