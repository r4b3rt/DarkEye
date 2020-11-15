package xraypoc

import (
	"github.com/zsdevX/DarkEye/hack/poc/xraypoc/celtypes"
	"net/url"
)

func myReverseCheck(reverse *xraypoc_celtypes.Reverse, timeout int64) bool {
	//todo
	return true
}

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
