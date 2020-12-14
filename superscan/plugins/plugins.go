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
	checkFuncs    = map[int]func(*Plugins){}
	supportPlugin = map[string]string{}
	userList      []string
	passList      []string
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
			output(plg)
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
	plg.TargetPort = "137"
	nbCheck(plg)
	plg.TargetPort = "445"
	ms17010Check(plg)
	if plg.PortOpened || len(plg.Cracked) != 0 {
		output(plg)
	}

}

func SetDicByFile(userFile, passFile string) {
	if userFile != "" {
		userList = common.GenArraryFromFile(userFile)
	}
	if passFile != "" {
		passList = common.GenArraryFromFile(passFile)
	}
	if userList != nil {
		color.Green("使用用户字典 %s", userFile)
	}
	if passList != nil {
		color.Green("使用密码字典 %s", passFile)
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
					plg.Lock()
					if pass == "" || ok == OKNoauth {
						pass = "空"
					}
					plg.Cracked = append(plg.Cracked, Account{Username: username, Password: pass})
					plg.Unlock()
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

func output(plg *Plugins) {
	if plg.highLight {
		color.Green("\n%s %s:%s %v\n",
			plg.TargetProtocol, plg.TargetIp, plg.TargetPort, plg.Cracked)
	} else {
		color.Yellow("\n%s %s:%s %v\n",
			plg.TargetProtocol, plg.TargetIp, plg.TargetPort, plg.Cracked)
	}
}

func SupportPlugin() {
	for _, v := range supportPlugin {
		color.Green("%v,", v)
	}
	color.Green("to be continue.\n,")
}
