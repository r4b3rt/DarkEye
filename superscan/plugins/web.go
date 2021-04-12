package plugins

import (
	"context"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strings"
)

func webCheck(s *Service) {
	timeOutSec := Config.TimeOut / 1000
	if timeOutSec == 0 {
		timeOutSec = 1
	}
	ctx, _ := context.WithCancel(Config.ParentCtx)
	s.parent.Result.Web.Server, s.parent.Result.Web.Title, s.parent.Result.Web.Code =
		common.GetHttpTitle(ctx, "http", s.parent.TargetIp+":"+s.parent.TargetPort, timeOutSec)
	s.parent.Result.Web.Url = fmt.Sprintf("http://%s:%s", s.parent.TargetIp, s.parent.TargetPort)
	//部分http访问https有title
	if strings.Contains(s.parent.Result.Web.Title, "The plain HTTP request was sent to HTTPS port") {
		s.parent.Result.Web.Title = ""
	}
	if s.parent.Result.Web.Server == "" && s.parent.Result.Web.Title == "" {
		s.parent.Result.Web.Server, s.parent.Result.Web.Title, s.parent.Result.Web.Code =
			common.GetHttpTitle(ctx, "https", s.parent.TargetIp+":"+s.parent.TargetPort, timeOutSec)
		s.parent.Result.Web.Tls = true
		s.parent.Result.Web.Url = fmt.Sprintf("https://%s:%s", s.parent.TargetIp, s.parent.TargetPort)
	}
	if s.parent.Result.Web.Server != "" || s.parent.Result.Web.Title != "" {
		if s.parent.Result.Web.Tls {
			s.parent.Result.ServiceName = "https"
		} else {
			s.parent.Result.ServiceName = "http"
		}
		s.parent.Hit = true
		webCrackByFinger(s)
	}
}

func webCrackByFinger(s *Service) {
	if strings.Contains(s.parent.Result.Web.Title, "Apache Tomcat") {
		tomcatCheck(s)
		return
	}
	//Other
	checkWebLogic(s)
}
