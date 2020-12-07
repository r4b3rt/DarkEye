package plugins

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/zsdevX/DarkEye/common"
	"strings"
	"sync"
)

var (
	checkFuncs = map[int]func(*Plugins){}
	userList   []string
	passList   []string
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
		color.Yellow("\n%s %s:%s %v\n",
			plg.TargetProtocol, plg.TargetIp, plg.TargetPort, plg.NetBios)
	}
}

func SetDicByFile(userfile, passfile string) {
	if userfile != "" {
		userList = common.GenArraryFromFile(userfile)
	}
	if passfile != "" {
		passList = common.GenArraryFromFile(passfile)
	}
	if userList != nil {
		color.Green("使用用户字典 %s", userfile)
	}
	if passList != nil {
		color.Green("使用密码字典 %s", passfile)
	}
	return
}

func crack(pid string, plg *Plugins, dictUser, dictPass []string, callback func(*Plugins, string, string) int) {
	//如果用户指定字典强制切换
	if userList != nil {
		dictUser = userList
	}
	if passList != nil {
		dictPass = passList
	}
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
				if pass == "空" {
					pass = ""
				}
				pass = strings.Replace(pass, "%user%", username, -1)
				plg.DescCallback(fmt.Sprintf("Cracking %s %s:%s %s/%s",
					pid, plg.TargetIp, plg.TargetPort, username, pass))
				ok := callback(plg, username, pass)
				switch ok {
				case OKNoauth:
					fallthrough
				case OKDone:
					//密码正确一次退出
					plg.locker.Lock()
					if pass == "" || ok == OKNoauth {
						pass = "空"
					}
					plg.Cracked = append(plg.Cracked, Account{Username: username, Password: pass})
					plg.locker.Unlock()
					plg.TargetProtocol = pid
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
				case OKForbidden:
					plg.TargetProtocol = pid
					color.Red("\n%s服务器配置受限。影响主机:%s:%s",
						pid, plg.TargetIp, plg.TargetPort)
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
