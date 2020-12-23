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
	"time"
)

var (
	myOS            = runtime.GOOS
	myCommand       = "ping -c 1 -w 1"
	myCommandOutput = "ttl="
	myShell         = "sh -c "
	mPrivileged     = false
)

func init() {
	if myOS == "windows" {
		myCommandOutput = "TTL="
		myShell = "CMD /c "
		myCommand = "ping -n 1 -w 1"
	} else if myOS == "darwin" {
		myCommand = "ping -c 1 -W 1"
		myCommandOutput = ", 0.0%"
	}
}

//PingNet 将网络分若干个C段，每个C段为最小单位整体扫描
func (s *Scan) PingNet(ipList string) {
	s.pingPrepare()
	ips := strings.Split(ipList, ",")
	for _, ip := range ips {
		if !strings.Contains(ip, "-") {
			//如果没有范围，则探测该段下所有ip活跃情况
			s.pingHost(ip)
			continue
		}
		//如果有'-'，则只探测网段存活
		base, _, end, err := common.GetIPRange(ip)
		if err != nil {
			color.Red("%v", err)
			return
		}
		ipSeg := net.ParseIP(base).To4()
		for {
			ipSeg[3] = 0
			if s.pingCheck(ipSeg.String(), false) {
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

func (s *Scan) pingHost(ip string) {
	base, _, _, err := common.GetIPRange(ip)
	if err != nil {
		color.Red("%v", err)
		return
	}
	ipSeg := net.ParseIP(base).To4()
	ipSeg[3] = 0
	s.pingCheck(ipSeg.String(), true)

}

//最小的ping单位：1个c段
func (s *Scan) pingCheck(ipSeg string, perHost bool) bool {
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
			select {
			case <-ctx.Done():
				return
			default:
			}
			tip := common.GenIP(ipSeg, idx)
			var ok bool
			if mPrivileged {
				ok = s.pingWithPrivileged(ctx, tip) == nil
			} else {
				ok = s.ping(ctx, tip)
			}
			if ok {
				if perHost {
					color.Green("%s is alive", tip)
				} else {
					alive.Store(true)
					cancel()
				}
			}
		}(start)
	}
	wg.Wait()
	return alive.Load()
}

func (s *Scan) pingWithPrivileged(ctx context.Context, ip string) error {
	data := []byte{8, 0, 247, 255, 0, 0, 0, 0}
	d := net.Dialer{Timeout: time.Duration(s.TimeOut) * time.Millisecond}
	conn, err := d.DialContext(ctx, "ip4:icmp", ip)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.Write(data); err != nil {
		return err
	}
	_ = conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(s.TimeOut)))
	recv := make([]byte, 1024)
	_, err = conn.Read(recv)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scan) ping(ctx context.Context, ip string) bool {
	cmd := strings.Split(myShell, " ")
	c := exec.CommandContext(ctx, cmd[0], cmd[1], myCommand+" "+ip)
	//common.HideCmd(c)
	b, _ := c.Output()
	if b != nil {
		return bytes.Contains(b, []byte(myCommandOutput))
	}
	return false
}

func (s *Scan) pingPrepare() {
	if s.pingWithPrivileged(context.Background(), "127.0.0.1") == nil {
		mPrivileged = true
		return
	}
	color.Yellow("当前为非管理权限模式，需要使用原生的命令（例如：ping）检测。请设置命令参数：")
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
	color.Yellow("\n使用命令Shell环境'%s'", myShell)
	color.Yellow("使用探测命令 '%s'检查网络 ", myCommand)
	color.Yellow("使用关键字' %s' 确定网络是否存在", myCommandOutput)
}
