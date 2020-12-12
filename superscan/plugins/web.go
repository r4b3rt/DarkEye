package plugins

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"strings"
)

func init() {
	checkFuncs[WEBSrv] = webCheck
	supportPlugin["tomcat"] = "tomcat"
}

func webCheck(plg *Plugins) {
	timeOutSec := plg.TimeOut / 1000
	if timeOutSec == 0 {
		timeOutSec = 1
	}
	cracked := Account{}
	plg.RateWait(plg.RateLimiter)
	cracked.Server, cracked.Title, cracked.Code = common.GetHttpTitle("http", plg.TargetIp+":"+plg.TargetPort, timeOutSec)
	cracked.Url = fmt.Sprintf("http://%s:%s", plg.TargetIp, plg.TargetPort)
	//部分http访问https有title
	if strings.Contains(cracked.Title, "The plain HTTP request was sent to HTTPS port") {
		cracked.Title = ""
	}
	if cracked.Server == "" && cracked.Title == "" {
		cracked.Server, cracked.Title, cracked.Code = common.GetHttpTitle("https", plg.TargetIp+":"+plg.TargetPort, timeOutSec)
		cracked.Tls = true
		cracked.Url = fmt.Sprintf("https://%s:%s", plg.TargetIp, plg.TargetPort)
	}
	if cracked.Server != "" || cracked.Title != "" {
		plg.TargetProtocol = "web"
		webTodo(plg, &cracked)
		plg.Lock()
		plg.Cracked = append(plg.Cracked, cracked)
		plg.Unlock()
	}
}

func webTodo(plg *Plugins, ck *Account) {
	if strings.Contains(ck.Title, "Apache Tomcat") {
		//爆破manager
		plg.tmp.tls = ck.Tls
		plg.tmp.urlPath = "/manager/html"
		plg.TargetProtocol = "tomcat"
		basicAuthCheck(plg, dic.DIC_USERNAME_TOMCAT, dic.DIC_PASSWORD_TOMCAT)
		return
	}
	//Other
	checkWebLogic(plg)
}
