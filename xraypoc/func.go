package xraypoc

import (
	"bytes"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/xraypoc/celtypes"
	"net/url"
	"strings"
	"time"
)

func myReverseCheck(reverse *xraypoc_celtypes.Reverse, timeout int64) bool {
	if reverse.ReverseCheckUrl == "" {
		return false
	}
	time.Sleep(time.Second * time.Duration(timeout))
	filter := strings.Split(reverse.Domain, ".")[0]
	if len(filter) >= 20 { //filter长度限制
		filter = string([]byte(filter)[:20])
	}
	apiUrl := fmt.Sprintf("%s%s", reverse.ReverseCheckUrl, filter)
	req := common.HttpRequest{
		Method:  "GET",
		Url:     apiUrl,
		TimeOut: 10,
	}
	response, err := req.Go()
	if err != nil {
		return false
	}
	if !bytes.Contains(response.Body, []byte(`"data": []`)) {
		return true
	}
	return false
}

//UrlConvertString add comment
func UrlConvertString(url *xraypoc_celtypes.UrlType) (myUrl string) {
	if url.Scheme == "" {
		return
	}
	myUrl = url.Scheme + ":" + "//" + url.Host
	if url.Path != "" && url.Path[0] != '/' {
		myUrl += "/"
	}
	myUrl += url.Path
	if url.Query != "" {
		myUrl += "?" + url.Query
	}
	if url.Fragment != "" {
		myUrl += "#" + url.Fragment
	}
	return
}

//StringConvertUrl add comment
func StringConvertUrl(myUrl string) (*xraypoc_celtypes.UrlType, error) {
	u, err := url.Parse(myUrl)
	if err != nil {
		return nil, err
	}
	return &xraypoc_celtypes.UrlType{
		Scheme:   u.Scheme,
		Domain:   u.Hostname(),
		Host:     u.Host,
		Port:     u.Port(),
		Path:     u.EscapedPath(),
		Query:    u.RawQuery,
		Fragment: u.Fragment,
	}, nil
}
