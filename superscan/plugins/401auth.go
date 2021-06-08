package plugins

import (
	"context"
	"encoding/base64"
	"github.com/zsdevX/DarkEye/common"
	"net/http"
	"strings"
	"time"
)

func _401AuthCheck(s *Service) {
	s.crack()
}

func _401AuthConn(parent context.Context, s *Service, user, pass string) (ok int) {
	ok = OKNext
	authKey := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	ctx, _ := context.WithCancel(parent)
	vul, _ := s.parent.Result.Output.Get("http_vul_page")
	req := common.HttpRequest{
		Url:              vul.(string),
		TimeOut:          time.Duration(1 + Config.TimeOut/1000),
		Method:           "GET",
		NoFollowRedirect: true,
		Headers: map[string]string{
			"Authorization": "Basic " + authKey,
			"User-Agent":    common.UserAgents[0],
		},
		Ctx: ctx,
	}
	response, err := req.Go()
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
	if response.Status == http.StatusOK {
		ok = OKDone
	} else if response.Status != http.StatusUnauthorized {
		ok = OKStop
	}
	return
}
