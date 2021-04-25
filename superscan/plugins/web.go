package plugins

import (
	"context"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strings"
)

func webCheck(s *Service) {
	timeOutSec := 1 + Config.TimeOut/1000
	ctx, _ := context.WithCancel(Config.ParentCtx)
	s.parent.Result.Web.Server, s.parent.Result.Web.Title, s.parent.Result.Web.Code =
		common.GetHttpTitle(ctx, "http", s.parent.TargetIp+":"+s.parent.TargetPort, timeOutSec)
	s.parent.Result.Web.Url = fmt.Sprintf("http://%s:%s", s.parent.TargetIp, s.parent.TargetPort)
	//部分http访问https有title
	if strings.Contains(s.parent.Result.Web.Title, "The plain HTTP request was sent to HTTPS port") {
		s.parent.Result.Web.Title = ""
	}
	//http失败后尝试https访问
	if s.parent.Result.Web.Server == "" && s.parent.Result.Web.Title == "" {
		s.parent.Result.Web.Server, s.parent.Result.Web.Title, s.parent.Result.Web.Code =
			common.GetHttpTitle(ctx, "https", s.parent.TargetIp+":"+s.parent.TargetPort, timeOutSec)
		s.parent.Result.Web.Tls = true
		s.parent.Result.Web.Url = fmt.Sprintf("https://%s:%s", s.parent.TargetIp, s.parent.TargetPort)
	}
	//记录
	if s.parent.Result.Web.Server != "" || s.parent.Result.Web.Title != "" {
		if s.parent.Result.Web.Tls {
			s.parent.Result.ServiceName = "https"
		} else {
			s.parent.Result.ServiceName = "http"
		}
		s.parent.Hit = true
		//尝试web相关爆破
		webCrack(s)
	}
}

func webCrack(s *Service) {
	//尝试找爆破的web
	switch WhatWeb(s.parent.Result.Web.Title) {
	case Tomcat:
		tomcatCheck(s)
	case WebLogic:
		checkWebLogic(s)
	}
}

func init() {
	services["web"] = Service{
		name:  "web",
		check: webCheck,
	}
}

const (
	UnknownWeb = iota
	WebLogic
	Tomcat
)

func WhatWeb(finger string) int {
	if strings.Contains(finger, "Apache Tomcat") {
		return Tomcat
	}
	if strings.Contains(finger, "WebLogic") {
		return WebLogic
	}
	return UnknownWeb
}
