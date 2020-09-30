package common

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Http struct {
	Method         string
	Url            string
	Data           []byte
	Cookie         string
	Origin         string
	ContentType    string
	Referer        string
	Agent          string
	H              string
	TimeOut        time.Duration
	ResponseServer string
}

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("stopped after 64 redirects")
	}
	return nil
}

func (m *Http) Http() ([]byte, error) {
	ctx, cancel := context.WithCancel(context.TODO()) // or parant context
	_ = time.AfterFunc(m.TimeOut*time.Second, func() {
		cancel()
	})
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	jar, _ := cookiejar.New(nil)
	cli := http.Client{
		Transport:     tr,
		CheckRedirect: defaultCheckRedirect,
		Jar:           jar,
	}
	req, _ := http.NewRequest(m.Method, m.Url, bytes.NewReader(m.Data))
	req = req.WithContext(ctx)
	if m.Method == "POST" {
		req.Header.Add("Content-Type", m.ContentType)
	}
	if m.Origin != "" {
		req.Header.Add("Origin", m.Origin)
	}
	if m.Referer == "" {
		v, _ := url.Parse(m.Url)
		m.Referer = v.Scheme + "://" + v.Host
	}
	req.Header.Add("Referer", m.Referer)
	req.Header.Add("User-Agent", m.Agent)
	req.Header.Add("Accept-Encoding", "xzip")
	if m.Cookie != "" {
		req.Header.Add("Cookie", m.Cookie)
	}
	if m.H != "" {
		h := strings.Split(m.H, "=")
		req.Header.Add(h[0], h[1])
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(fmt.Sprintf("Bad status %d", resp.StatusCode))
	}
	m.ResponseServer = resp.Header.Get("Server ")
	m.ResponseServer += resp.Header.Get("X-Powered-By")
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func IsAlive(ip, port string, timeout int) bool {
	ctx, cancel := context.WithCancel(context.TODO()) // or parant context
	_ = time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
		cancel()
	})
	d := net.Dialer{}
	c, err := d.DialContext(ctx, "tcp", ip+":"+port)
	if err != nil {
		return false
	}
	defer c.Close()
	return true
}

func GetHttpTitle(proto, domain string) (server, title string) {
	url := fmt.Sprintf(proto+"://%s", domain)
	userAgent := UserAgents[0]
	req := Http{
		Url:     url,
		TimeOut: time.Duration(5),
		Method:  "GET",
		Referer: url,
		Agent:   userAgent,
	}
	body, err := req.Http()
	if err != nil {
		return
	}
	server = req.ResponseServer
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		return
	}
	t := htmlquery.Find(doc, "//title")
	if len(t) != 0 {
		title = htmlquery.InnerText(t[0])
	}
	if !ISUtf8([]byte(title)) {
		if message, err := simplifiedchinese.GBK.NewDecoder().String(title); err == nil {
			 title = message
		}
	}
	title = TrimUseless(title)

	return
}
