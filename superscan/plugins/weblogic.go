package plugins

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"strings"
	"time"
)

func init() {
	supportPlugin["webLogic"] = "webLogic"
}

func checkWebLogic(plg *Plugins) {
	response, _ := webLogicTest(plg, "test", "test")
	if response == nil {
		return
	}
	if loc, ok := response.ResponseHeaders["Location"]; ok {
		//判断webLogic指纹
		if !strings.Contains(loc, "LoginForm.jsp") {
			return
		}
		//ADMINCONSOLESESSION=O9BKOxQLzl7ZoJBlf4IgIicF0g0WGpfNMrUaSWWIA2G5gdlL6yvH!-1249850057;
		if ck, ok := response.ResponseHeaders["Set-Cookie"]; ok {
			cks := strings.Split(ck, ";")
			if len(cks) >= 1 {
				plg.tmp.cookie = cks[0]
			}
		}
	} else {
		return
	}

	crack("webLogic", plg, dic.DIC_USERNAME_WEBLOGIC, dic.DIC_PASSWORD_WEBLOGIC, webLogicConn)
}

func webLogicConn(plg *Plugins, user string, pass string) (ok int) {
	ok = OKNext

	response, err := webLogicTest(plg, user, pass)
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

func webLogicTest(plg *Plugins, user, pass string) (*common.HttpResponse, error) {
	url := fmt.Sprintf("http://%s:%s/console/j_security_check", plg.TargetIp, plg.TargetPort)
	if plg.tmp.tls {
		url = fmt.Sprintf("https://%s:%s/console/j_security_check", plg.TargetIp, plg.TargetPort)
	}
	req := common.HttpRequest{
		Url:     url,
		TimeOut: time.Duration(1 + plg.TimeOut/1000),
		Method:  "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent":   common.UserAgents[0],
			"Cookie":       plg.tmp.cookie,
		},
		NoFollowRedirect: true,
		Body: []byte(fmt.Sprintf("j_username=%s&j_password=%s&j_character_encoding=UTF-8",
			user, pass)),
	}
	return req.Go()
}
