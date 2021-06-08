package plugins

import (
	"context"
	"fmt"
	"github.com/b1gcat/DarkEye/common"
	"net"
)

func webCheck(s *Service) {
	timeOutSec := 1 + Config.TimeOut/1000
	ctx, _ := context.WithCancel(Config.ParentCtx)

	targets := make([]string, 0)
	target := net.JoinHostPort(s.parent.TargetIp, s.parent.TargetPort)
	targets = append(targets, target)
	if Config.WebSiteDomainList != nil {
		for _, t := range Config.WebSiteDomainList {
			targets = append(targets, t)
		}
	}

	httpUrl := ""
	httpTitle := ""
	httpCode := ""
	httpServer := ""

	defer func() {
		if s.parent.Hit {
			s.parent.Result.Output.Set("http_url", httpUrl)
			s.parent.Result.Output.Set("http_server", httpServer)
			s.parent.Result.Output.Set("http_code", httpCode)
			s.parent.Result.Output.Set("http_title", httpTitle)
		}
	}()
	for k, site := range targets {
		for _, proto := range []string{"http", "https"} {
			server, title, code := common.GetHttpTitle(ctx, proto, target, site, timeOutSec)
			//网不通或非http协议
			if code == 0 {
				continue
			}
			s.parent.Hit = true

			s.parent.Result.ServiceName = "http[s]"
			httpUrl += fmt.Sprint(k, ": ", proto, "://", site, "\r\n")
			httpServer += fmt.Sprint(k, ": ", server, "\r\n")
			httpTitle += fmt.Sprint(k, ": ", title, "\r\n")
			httpCode += fmt.Sprint(k, ": ", fmt.Sprint(code), "\r\n")
		}
	}
}

func init() {
	services["web"] = Service{
		name:  "web",
		check: webCheck,
	}
}
