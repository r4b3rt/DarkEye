package plugins

import (
	"context"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strings"
	"time"
)

func checkWebLogic(s *Service) {
	response, _ := webLogicTest(s, "test", "test")
	if response == nil {
		return
	}
	if loc, ok := response.ResponseHeaders["Location"]; ok {
		//判断webLogic指纹
		if !strings.Contains(loc, "LoginForm.jsp") {
			return
		}
		//随意填写一个过期cookie
		s.vars["cookie"] = "ADMINCONSOLESESSION=O9BKOxQLzl7ZoJBlf4IgIicF0g0WGpfNMrUaSWWIA2G5gdlL6yvH!-1249850057"
		if ck, ok := response.ResponseHeaders["Set-Cookie"]; ok {
			cks := strings.Split(ck, ";")
			if len(cks) >= 1 {
				s.vars["cookie"] = cks[0]
			}
		}
	} else {
		return
	}
	s.name = "webLogic"
	s.parent.Result.ServiceName = s.name
	s.connect = webLogicConn
	s.crack()
}

func webLogicConn(parent context.Context, s *Service, user, pass string) (ok int) {
	ok = OKNext
	response, err := webLogicTest(s, user, pass)
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			//连接限制
			ok = OKWait
			return
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			ok = OKTimeOut
			return
		}
		if response == nil {
			//异常?
			ok = OKStop
			return
		}
	}
	if loc, ok := response.ResponseHeaders["Location"]; ok {
		//判断webLogic指纹
		if !strings.Contains(loc, "LoginForm.jsp") {
			return OKDone
		}
	} else {
		return OKDone
	}

	return ok
}

func webLogicTest(s *Service, user, pass string) (*common.HttpResponse, error) {
	url := fmt.Sprintf("%s/console/j_security_check", s.parent.Result.Web.Url)
	ctx, _ := context.WithCancel(Config.ParentCtx)
	cookie, _ := s.vars["cookie"]
	req := common.HttpRequest{
		Ctx:     ctx,
		Url:     url,
		TimeOut: time.Duration(1 + Config.TimeOut/1000),
		Method:  "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent":   common.UserAgents[0],
			"Cookie":       cookie,
		},
		NoFollowRedirect: true,
		Body: []byte(fmt.Sprintf("j_username=%s&j_password=%s&j_character_encoding=UTF-8",
			user, pass)),
	}
	return req.Go()
}
