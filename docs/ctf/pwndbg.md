# 全局
关闭地址随机化: 
```
echo 0 >/proc/sys/kernel/randomize_va_space
```
 
## Demo
```
from pwn import *

 sh = process('./pwn')
# sh = remote('x.x.x.x',110)

sh.send()
sh.interactive()
```

## brutecrack sample
```
import string
import subprocess
import sys
def main():
    bin="354e846facdf6f3d3205a1465d2fd811"
    flag=""
    length=0
    while 1:
        for i in string.printable:
            p = subprocess.Popen("./%s" % bin,stdin=subprocess.PIPE,stdout=subprocess.PIPE)
            p.stdin.write("1\n")
            p.stdin.write(flag+i+"\n")
            output = p.stdout.readlines()[-1]
            if Todo():
                length+=1
                flag+=i
                print flag
                break
if __name__=="__main__":
    main()
```

## 通信API
```
context.log_level = "debug"
send(payload) 发送payload
sendline(payload) 发送payload，并进行换行（末尾\n）
sendafter(some_string, payload) 接收到 some_string 后, 发送你的 payload
recvn(N) 接受 N(数字) 字符
recvline() 接收一行输出
recvlines(N) 接收 N(数字) 行输出
recvuntil(some_string) 接收到 some_string 为止
```

## 地址API
```
>>> e = ELF('/bin/cat')
>>> print hex(e.address)  # 文件装载的基地址
0x400000
>>> print hex(e.symbols['write']) # 函数地址
0x401680
>>> print hex(e.got['write']) # GOT表的地址
0x60b070
>>> print hex(e.plt['write']) # PLT的地址
0x401680
>>> print hex(e.search('/bin/sh').next())# 字符串/bin/sh的地址
```
