package plugins

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	checkFuncs = map[int]func(*Plugins){}
)

func (plg *Plugins) Check() {
	plg.RateWait(plg.RateLimiter) //活跃端口发包限制
	if common.IsAlive(plg.TargetIp, plg.TargetPort, plg.TimeOut) != common.Alive {
		return
	}
	plg.PortOpened = true
	plg.Cracked = make([]Account, 0)
	i := 0
	//预处理
	for i < PluginNR {
		checkFuncs[i](plg)
		plg.DescCallback("Cracking ...")
		//未找到密码
		if plg.TargetProtocol != "" {
			if plg.highLight {
				color.Green("\n%s %s:%s %v\n",
					plg.TargetProtocol, plg.TargetIp, plg.TargetPort, plg.Cracked)
			} else {
				color.Yellow("\n%s %s:%s %v\n",
					plg.TargetProtocol, plg.TargetIp, plg.TargetPort, plg.Cracked)
			}
			break
		}
		i++
	}
	if i >= PluginNR {
		color.Yellow("\n%s %s:%s %v\n", "[-]",
			plg.TargetIp, plg.TargetPort, "Opened")
	}
	return
}

func (plg *Plugins) PreCheck() {
	//预处理
	//137端口机器检查
	nbCheck(plg)
	if plg.PortOpened {
		color.Yellow("\n%s %s:%s %v\n", plg.TargetProtocol, plg.TargetIp, 137, plg.NetBios)
	}
}

func loadDic(name string) []string {
	filename := filepath.Join("dic", fmt.Sprintf("dic_%s", name))
	file, err := os.Open(filename)
	if err != nil {
		color.Red("未发现字典文件:" + filename)
		return nil
	}
	defer file.Close()
	result := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		one := scanner.Text()
		if strings.HasPrefix(one, "#") {
			continue
		}
		one = strings.TrimSpace(one)
		one = strings.Trim(one, "\r\n")
		if one == "空" { //超级口令的""特殊表示
			result = append(result, "")
		} else {
			result = append(result, one)
		}
	}
	return result
}

func crack(pid string, plg *Plugins, dictUser, dictPass []string, callback func(*Plugins, string, string) int) {
	wg := sync.WaitGroup{}
	wg.Add(len(dictUser))
	limiter := make(chan int, plg.Worker)
	ctx, cancel := context.WithCancel(context.TODO())
	for _, user := range dictUser {
		limiter <- 1
		go func(username string) {
			defer func() {
				<-limiter
				wg.Done()
			}()
			for _, pass := range dictPass {
				select {
				case <-ctx.Done():
					return
				default:
				}
				pass = strings.Replace(pass, "%user%", username, -1)
				plg.DescCallback(fmt.Sprintf("Cracking %s %s:%s %s/%s",
					pid, plg.TargetIp, plg.TargetPort, username, pass))
				ok := callback(plg, username, pass)
				switch ok {
				case OKDone:
					//密码正确一次退出
					plg.locker.Lock()
					plg.Cracked = append(plg.Cracked, Account{Username: username, Password: pass})
					plg.locker.Unlock()
					plg.highLight = true
					cancel()
					return
				case OKWait:
					//太快了服务器限制
					color.Red("\n%s爆破受限，建议降低参数'plugin-worker'数值.影响主机:%s:%s",
						pid, plg.TargetIp, plg.TargetPort)
					cancel()
					return
				case OKTimeOut:
					color.Red("\n%s爆破超时，建议提高参数'timeout'数值.影响主机:%s:%s",
						pid, plg.TargetIp, plg.TargetPort)
					cancel()
					return
				case OKStop:
					//非协议退出
					cancel()
					return
				default:
					//密码错误.OKNext
					plg.TargetProtocol = pid
				}
			}
		}(user)
	}
	wg.Wait()
}
