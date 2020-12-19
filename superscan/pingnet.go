//快速发现内网有效的网段
//native的方法有需要root权限，所以暂时用自动自带的ping
package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/beats/libbeat/common/atomic"
	"github.com/fatih/color"
	"github.com/zsdevX/DarkEye/common"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	myOS            = runtime.GOOS
	myCommand       = "ping -c 1 -w 1"
	myCommandOutput = "ttl="
	myShell         = "sh -c "
)

func init() {
	if myOS == "windows" {
		myCommandOutput = "TTL="
		myShell = "CMD /c "
	} else if myOS == "darwin" {
		myCommand = "ping -c 1 -W 1"
		myCommandOutput = ", 0.0%"
	}
}

//将网络分若干个C段，每个C段为最小单位整体扫描
func (s *Scan) PingNet(ipList string) {
	s.pingPrepare()
	ips := strings.Split(ipList, ",")
	for _, ip := range ips {
		base, _, end, err := common.GetIPRange(ip)
		if err != nil {
			color.Red("%v", err)
			return
		}
		ipSeg := net.ParseIP(base).To4()
		for {
			ipSeg[3] = 0
			if s.pingCheck(ipSeg.String()) {
				color.Green("%s is alive", ipSeg.String())
			} else {
				color.Yellow("%s is died", ipSeg.String())
			}
			//判断是否进行下一个C扫
			ipSeg[3] = 0xff
			if common.CompareIP(ipSeg.String(), end) >= 0 {
				break
			}
			//下一个C段
			ipSeg = net.ParseIP(common.GenIP(ipSeg.String(), 1)).To4()
		}
	}
}

//最小的ping单位：1个c段
func (s *Scan) pingCheck(ipSeg string) bool {
	start := 0
	alive := atomic.Bool{}
	ctx, cancel := context.WithCancel(context.TODO())
	wg := sync.WaitGroup{}
	wg.Add(254)
	for {
		start++
		if start >= 255 {
			break
		}
		go func(idx int) {
			defer func() {
				wg.Done()
			}()
			if s.ping(common.GenIP(ipSeg, idx), ctx) {
				alive.Store(true)
				cancel()
			}
		}(start)
	}
	wg.Wait()
	return alive.Load()
}

func (s *Scan) ping(ip string, ctx context.Context) bool {
	cmd := strings.Split(myShell, " ")
	c := exec.CommandContext(ctx, cmd[0], cmd[1], myCommand+" "+ip)
	common.HideCmd(c)
	b, _ := c.Output()
	if b != nil {
		return bytes.Contains(b, []byte(myCommandOutput))
	}
	return false
}

func (s *Scan) pingPrepare() {
	var cmd string
	_, _ = fmt.Fprintf(os.Stderr, "输入探测命令(default: %s):", myCommand)
	n, _ := fmt.Scanln(&cmd)
	if n != 0 {
		myCommand = cmd
	}
	_, _ = fmt.Fprintf(os.Stderr, "输入探测的成功关键字(default: %s)", myCommandOutput)
	n, _ = fmt.Scanln(&cmd)
	if n != 0 {
		myCommandOutput = cmd
	}
	_, _ = fmt.Fprintf(os.Stderr, "输入命令shell环境(default: %s)", myShell)
	n, _ = fmt.Scanln(&cmd)
	if n != 0 {
		myShell = cmd
	}
	color.Yellow("使用命令Shell环境'%s'", myShell)
	color.Yellow("使用探测命令 '%s'检查网络 ", myCommand)
	color.Yellow("使用关键字' %s' 确定网络是否存在", myCommandOutput)
}
