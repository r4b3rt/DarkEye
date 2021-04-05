package plugins

import (
	"context"
	"encoding/json"
	"github.com/fatih/color"
	"github.com/zsdevX/DarkEye/common"
	"strings"
	"sync"
)

//PreCheck add comment
func (plg *Plugins) PreCheck() {
	if ShouldStop() {
		return
	}
	//预处理注意：
	//1、该链上的处理为固定端口，主要为UDP或特殊协议
	//2、此处未做发包限制
	i := 0
	for _, v := range preCheckFuncs {
		v.doit(plg, &v)
		i++
	}
	if len(plg.Cracked) != 0 {
		output(plg)
	}
}

//Check add comment
func (plg *Plugins) Check() {
	if ShouldStop() {
		return
	}
	GlobalConfig.RateWait(GlobalConfig.Pps) //活跃端口发包限制
	if !plg.PortOpened &&
		common.IsAlive(plg.TargetIp, plg.TargetPort, plg.TimeOut) != common.Alive {
		return
	}
	plg.PortOpened = true
	plg.Cracked = make([]Account, 0)
	i := 0
	//爆破链
	for _, v := range checkFuncs {
		if plg.available(v.name, v.port) {
			v.doit(plg, &v)
			//未找到密码
			if plg.TargetProtocol != "" {
				output(plg)
				break
			}
		}
		i++
	}
	if i >= len(checkFuncs) {
		color.Yellow("\n%s %s:%s %v\n", "[√]",
			plg.TargetIp, plg.TargetPort, "Opened")
	}
	return
}

func crack(pid string, plg *Plugins, dictUser, dictPass []string, callback func(*Plugins, string, string) int) {
	//如果用户指定字典强制切换
	if GlobalConfig.UserList != nil {
		dictUser = GlobalConfig.UserList
	}
	if GlobalConfig.PassList != nil {
		dictPass = GlobalConfig.PassList
	}
	wg := sync.WaitGroup{}
	wg.Add(len(dictUser))
	limiter := make(chan int, GlobalConfig.Thread)
	ctx, cancel := context.WithCancel(context.TODO())
	for _, user := range dictUser {
		limiter <- 1
		go func(username string) {
			defer func() {
				<-limiter
				wg.Done()
			}()
			if ShouldStop() {
				return
			}
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
				//限速
				GlobalConfig.RateWait(GlobalConfig.Pps)
				ok := callback(plg, username, pass)
				switch ok {
				case OKNoAuth:
					fallthrough
				case OKDone:
					//密码正确一次退出
					plg.Lock()
					if pass == "" || ok == OKNoAuth {
						pass = "空"
					}
					plg.Cracked = append(plg.Cracked, Account{Username: username, Password: pass})
					plg.TargetProtocol = pid
					plg.highLight = true
					plg.Unlock()

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
	output, _ := json.Marshal(plg.Cracked)
	//var out bytes.Buffer
	//_ = json.Indent(&out, output, "", "\t")
	if plg.highLight {
		color.Green("\n[√] %s %s:%s %v\n",
			plg.TargetProtocol, plg.TargetIp, plg.TargetPort, string(output))
	} else {
		color.Yellow("\n[√] %s %s:%s %v\n",
			plg.TargetProtocol, plg.TargetIp, plg.TargetPort, string(output))
	}
}

func (plg *Plugins) available(name, port string) bool {
	//强制指纹识别的协议
	if port == "" {
		return true
	}
	if plg.TargetPort != port {
		if GlobalConfig.UsingPlugin == "" {
			return false
		}
	}
	return true
}

//SupportPlugin add comment
func SupportPlugin() {
	list := ""
	defer func() {
		color.Green("%v", list)
	}()

	if GlobalConfig.UsingPlugin != "" {
		list += GlobalConfig.UsingPlugin
		return
	}
	for k := range checkFuncs {
		list += k + " "
	}

	for k := range preCheckFuncs {
		list += k + " "
	}
	strings.TrimSpace(list)
}

func ShouldStop() bool {
	select {
	case <-GlobalConfig.Ctx.Done():
		GlobalConfig.Stop.Store(true)
		return true
	default:
		if GlobalConfig.Stop.Load() {
			return true
		}
	}
	return false
}
