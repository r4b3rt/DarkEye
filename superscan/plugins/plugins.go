package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/zsdevX/DarkEye/common"
	"golang.org/x/time/rate"
	"strings"
	"sync"
	"time"
)

var (
	checkFuncs    = map[int]func(*Plugins){}
	preCheckFuncs = map[int]func(*Plugins){}
	supportPlugin = map[string]string{}
	//GlobalConfig add comment
	GlobalConfig = Config{
		ReverseUrl:      "qvn0kc.ceye.io",
		ReverseCheckUrl: "http://api.ceye.io/v1/records?token=066f3d242991929c823ac85bb60f4313&type=http&filter=",
		RateWait: func(r *rate.Limiter) {
			if r == nil {
				return
			}
			for {
				if r.Allow() {
					break
				} else {
					time.Sleep(time.Millisecond * 10)
				}
			}
		},
	}
)

//PreCheck add comment
func (plg *Plugins) PreCheck() {
	//预处理注意：
	//1、该链上的处理为固定端口，主要为UDP或特殊协议
	//2、此处未做发包限制
	i := 0
	for i < PluginPreCheckNR {
		plg.DescCallback(fmt.Sprintf("Cracking initiated %s", plg.TargetIp))
		preCheckFuncs[i](plg)
		i++
	}
	if len(plg.Cracked) != 0 {
		output(plg)
	}
}

//Check add comment
func (plg *Plugins) Check() {
	GlobalConfig.RateWait(GlobalConfig.Pps) //活跃端口发包限制
	plg.DescCallback(fmt.Sprintf("Crack %s:%s", "*."+plg.TargetIp[len(plg.TargetIp)-3:], plg.TargetPort))
	if !plg.PortOpened &&
		common.IsAlive(plg.TargetIp, plg.TargetPort, plg.TimeOut) != common.Alive {
		return
	}
	plg.PortOpened = true
	plg.Cracked = make([]Account, 0)
	i := 0
	//爆破链
	for i < PluginNR {
		checkFuncs[i](plg)
		//未找到密码
		if plg.TargetProtocol != "" {
			output(plg)
			break
		}
		i++
	}
	if i >= PluginNR {
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
				plg.DescCallback(fmt.Sprintf("Crack %s %s:%s %s/%s",
					pid, "*."+plg.TargetIp[len(plg.TargetIp)-3:], plg.TargetPort, username, pass))
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

//SupportPlugin add comment
func SupportPlugin() {
	for _, v := range supportPlugin {
		color.Green("%v,", v)
	}
}
