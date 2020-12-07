package plugins

import (
	"github.com/zsdevX/DarkEye/common"
	"strings"
)

func init() {
	checkFuncs[WEBSrv] = webCheck
}

func webCheck(plg *Plugins) {
	timeOutSec := plg.TimeOut / 1000
	if timeOutSec == 0 {
		timeOutSec = 1
	}
	cracked := Account{}
	plg.RateWait(plg.RateLimiter)
	cracked.Server, cracked.Title, cracked.Code = common.GetHttpTitle("http", plg.TargetIp+":"+plg.TargetPort, timeOutSec)
	//部分http访问https有title
	if strings.Contains(cracked.Title, "The plain HTTP request was sent to HTTPS port") {
		cracked.Title = ""
	}
	if cracked.Server == "" && cracked.Title == "" {
		cracked.Server, cracked.Title, cracked.Code = common.GetHttpTitle("https", plg.TargetIp+":"+plg.TargetPort, timeOutSec)
	}
	if cracked.Server != "" || cracked.Title != "" {
		plg.TargetProtocol = "web"
		plg.Lock()
		plg.Cracked = append(plg.Cracked, cracked)
		plg.Unlock()
	}
}
