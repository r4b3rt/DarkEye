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
	plg.RateWait(plg.RateLimiter)

	plg.Web.Server, plg.Web.Title, plg.Web.Code = common.GetHttpTitle("http", plg.TargetIp+":"+plg.TargetPort, timeOutSec)
	//部分http访问https有title
	if strings.Contains(plg.Web.Title, "The plain HTTP request was sent to HTTPS port") {
		plg.Web.Title = ""
	}
	if plg.Web.Server == "" && plg.Web.Title == "" {
		plg.Web.Server, plg.Web.Title, plg.Web.Code = common.GetHttpTitle("https", plg.TargetIp+":"+plg.TargetPort, timeOutSec)
	}
	if plg.Web.Server != "" || plg.Web.Title != "" {
		plg.TargetProtocol = "[web]"
	}
}
