package plugins

import (
	"encoding/base64"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"net/http"
	"time"
)

func basicAuthCheck(plg *Plugins, user, pass []string) {
	crack(plg.TargetProtocol, plg, user, pass, basicAuthConn)
}

func basicAuthConn(plg *Plugins, user string, pass string) (ok int) {
	ok = OKNext

	url := fmt.Sprintf("http://%s:%s%s", plg.TargetIp, plg.TargetPort, plg.tmp.urlPath)
	if plg.tmp.tls {
		url = fmt.Sprintf("https://%s:%s%s", plg.TargetIp, plg.TargetPort, plg.tmp.urlPath)
	}
	authKey := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	req := common.HttpRequest{
		Url:              url,
		TimeOut:          time.Duration(1 + plg.TimeOut/1000),
		Method:           "GET",
		NoFollowRedirect: true,
		Headers: map[string]string{
			"Authorization": "Basic " + authKey,
			"User-Agent":    common.UserAgents[0],
		},
	}
	response, err := req.Go()
	if err != nil {
		ok = OKStop
		return
	}
	if response.Status == http.StatusOK {
		ok = OKDone
	} else if response.Status != http.StatusUnauthorized {
		ok = OKStop
	}
	return
}
