package plugins

import (
	"context"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strings"
	"time"
)

//PreCheck add comment
func (plg *Plugins) PreCheck() {
	//预处理注意：
	//1、该链上的处理为固定端口，主要为UDP或特殊协议
	//2、此处未做发包限制
	for _, srv := range preServices {
		if !plg.available(srv.name, srv.port, true) {
			continue
		}
		s := new(Service)
		*s = srv
		s.parent = plg
		s.check(s)
	}
}

//Check add comment
func (plg *Plugins) Check() {
	Config.rateLimiter() //活跃端口发包限制
	if common.IsAlive(Config.ParentCtx, plg.TargetIp, plg.TargetPort, Config.TimeOut) != common.Alive {
		return
	}
	plg.Result.PortOpened = true
	for _, srv := range services {
		if !plg.available(srv.name, srv.port, false) {
			continue
		}
		s := new(Service)
		*s = srv
		s.parent = plg
		s.check(s)
		if s.parent.Hit {
			break
		}
	}
	return
}

func (s *Service) crack() {
	dictUser := s.user
	dictPass := s.pass
	//如果用户指定字典强制切换
	if Config.UserList != nil {
		dictUser = Config.UserList
	}
	if Config.PassList != nil {
		dictPass = Config.PassList
	}
	task := common.NewTask(s.thread, Config.ParentCtx)
	defer task.Wait("crack")
	for _, user := range dictUser {
		if !task.Job() {
			break
		}
		go func(username string) {
			defer task.UnJob()
			if stop := s.job(task.Ctx, username, dictPass); stop {
				task.Die()
				return
			}
		}(user)
	}
}

func (s *Service) job(parent context.Context, user string, dictPass []string) (stop bool) {
	if user == "空" {
		user = ""
	}
	for _, pass := range dictPass {
		if pass == "空" {
			pass = ""
		}
		pass = strings.Replace(pass, "%user%", user, -1)
		//限速
		Config.rateLimiter()
		ok := s.connect(parent, s, user, pass)
		if ok == OKStop || ok == OKTerm {
			//非协议退出
			return true
		}
		s.parent.Hit = true
		s.parent.Result.ServiceName = s.name
		switch ok {
		case OKNoAuth:
			fallthrough
		case OKDone:
			//密码正确一次退出
			if pass == "" || ok == OKNoAuth {
				pass = "空"
			}
			s.parent.Result.Cracked = Account{Username: user, Password: pass}
			return true
		case OKWait:
			//太快了服务器限制
			common.Log(s.name, fmt.Sprintf("爆破受限,影响主机:%s:%s",
				s.parent.TargetIp, s.parent.TargetPort), common.ALERT)
			return true
		case OKTimeOut:
			common.Log(s.name, fmt.Sprintf("爆破超时,影响主机:%s:%s",
				s.parent.TargetIp, s.parent.TargetPort), common.ALERT)
			return true
		case OKForbidden:
			common.Log(s.name, fmt.Sprintf("服务器配置受限。影响主机:%s:%s",
				s.parent.TargetIp, s.parent.TargetPort), common.ALERT)
			return true
		default:
			stop = false
			//密码错误.OKNext
		}
	}
	return false
}

func (c *config) rateLimiter() {
	if c.PPS == nil {
		return
	}
	for {
		if c.PPS.Allow() {
			break
		} else {
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func (plg *Plugins) available(srvName, srvPort string, preCheck bool) bool {
	if Config.SelectPlugin == "" {
		//未指定时，采用程序内置逻辑判断
		//1. 插件未指定端口强制扫描
		if srvPort == "" {
			return true
		}
		if preCheck {
			return true
		}
		//2. 插件指定端口扫描
		if plg.TargetPort == srvPort {
			return true
		}
		return false
	}
	//指定插件时, 按照用户意图爆破所有端口
	services := strings.Split(Config.SelectPlugin, ",")
	for _, s := range services {
		if srvName == s {
			return true
		}
	}
	return false
}

//SupportPlugin add comment
func SupportPlugin() {
	list := ""
	defer func() {
		common.Log("Plugins:", list, common.ALERT)
	}()
	if Config.SelectPlugin != "" {
		list = Config.SelectPlugin
		return
	}
	for k := range preServices {
		list += k + " "
	}
	for k := range services {
		list += k + " "
	}
	strings.TrimSpace(list)
}
