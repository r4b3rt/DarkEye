# Debug

## Gdb
```
打印指令
查看内存指令x：
x /nuf 0x123456 //常用，x指令的格式是：x空格/nfu，nfu代表三个参数
n代表显示几个单元（而不是显示几个字节，后面的u表示一个单元多少个字节），放在’/'后面
u代表一个单元几个字节，b(一个字节)，h(俩字节)，w(四字节)，g(八字节)
f代表显示数据的格式，f和u的顺序可以互换，也可以只有一个或者不带n，用的时候很灵活
x 按十六进制格式显示变量。
d 按十进制格式显示变量。
u 按十六进制格式显示无符号整型。
o 按八进制格式显示变量。
t 按二进制格式显示变量。
a 按十六进制格式显示变量。
c 按字符格式显示变量。
f 按浮点数格式显示变量。
s 按字符串显示。
b 按字符显示。
i 显示汇编指令。

x /10gx 0x123456 //常用，从0x123456开始每个单元八个字节，十六进制显示是个单元的数据
x /10xd $rdi //从rdi指向的地址向后打印10个单元，每个单元4字节的十进制数
x /10i 0x123456 //常用，从0x123456处向后显示十条汇编指令
```

### Debug attach tips
```
attach进程sigsegment错误不退出
catch fork
catch vfork
set follow-fork-mode child

dump memory ./result.txt 0x7fffffffdaa0 0x7fffffffdae0
```
#### Dump 
```
#!/bin/bash

grep rw-p /proc/$1/maps \
| sed -n 's/^\([0-9a-f]*\)-\([0-9a-f]*\) .*$/\1 \2/p' \
| while read start stop; do \
    gdb --batch --pid $1 -ex \
        "dump memory $1-$start-$stop.dump 0x$start 0x$stop"; \
done
```

#其它
```
Tip1:
x64中的前六个参数依次保存在RDI, RSI, RDX, RCX, R8和 R9中，如果还有更多的参数的话才会保存在栈上
```
