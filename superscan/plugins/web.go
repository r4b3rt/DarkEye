package plugins

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strings"
)

func webCheck(plg *Plugins, f *funcDesc) {
	timeOutSec := plg.TimeOut / 1000
	if timeOutSec == 0 {
		timeOutSec = 1
	}
	plg.Web.Server, plg.Web.Title, plg.Web.Code = common.GetHttpTitle("http", plg.TargetIp+":"+plg.TargetPort, timeOutSec)
	plg.Web.Url = fmt.Sprintf("http://%s:%s", plg.TargetIp, plg.TargetPort)
	//部分http访问https有title
	if strings.Contains(plg.Web.Title, "The plain HTTP request was sent to HTTPS port") {
		plg.Web.Title = ""
	}
	if plg.Web.Server == "" && plg.Web.Title == "" {
		plg.Web.Server, plg.Web.Title, plg.Web.Code = common.GetHttpTitle("https", plg.TargetIp+":"+plg.TargetPort, timeOutSec)
		plg.Web.Tls = true
		plg.Web.Url = fmt.Sprintf("https://%s:%s", plg.TargetIp, plg.TargetPort)
	}
	if plg.Web.Server != "" || plg.Web.Title != "" {
		if plg.Web.Tls {
			plg.TargetProtocol = "https"
		} else {
			plg.TargetProtocol = "http"
		}
		webCrackByFinger(plg, f)
	}
}

func webCrackByFinger(plg *Plugins, f *funcDesc) {
	if strings.Contains(plg.Web.Title, "Apache Tomcat") {
		tomcatCheck(plg)
		return
	}
	//Other
	checkWebLogic(plg, f)
}
