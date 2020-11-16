package common

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type HttpRequest struct {
	Method  string
	Url     string
	Body    []byte
	Headers map[string]string

	NoFollowRedirect bool
	TimeOut          time.Duration
}

type HttpResponse struct {
	Status          int32
	ResponseHeaders map[string]string
	Body            []byte
	ContentType     string
}

func (m *HttpRequest) Go() (*HttpResponse, error) {
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
	if m.NoFollowRedirect {
		cli.CheckRedirect = noCheckRedirect
	}
	req, err := http.NewRequest(m.Method, m.Url, bytes.NewReader(m.Body))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	for k, v := range m.Headers {
		req.Header.Set(k, v)
	}
	resp, err := cli.Do(req)
	if err != nil {
		if !strings.Contains(err.Error(), "forbidden redirects") {
			return nil, err
		}
	}
	defer resp.Body.Close()
	response := HttpResponse{
		Status:          int32(resp.StatusCode),
		ResponseHeaders: make(map[string]string),
		ContentType:     resp.Header.Get("Content-Type"),
	}
	for k := range resp.Header {
		if k != "Set-Cookie"{
			response.ResponseHeaders[k] = resp.Header.Get(k)
		}
	}
	for _,ck := range resp.Cookies() {
		response.ResponseHeaders["Set-Cookie"] += ck.String() + ";"
	}

	if resp == nil {
		return nil, nil
	}
	body, err := getRespBody(resp)
	if err != nil {
		return nil, err
	}
	response.Body = body
	return &response, nil
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
	req := HttpRequest{
		Url:     url,
		TimeOut: time.Duration(5),
		Method:  "GET",
		Headers: map[string]string{
			"User-Agent": userAgent,
		},
	}
	response, err := req.Go()
	if err != nil {
		return
	}
	server = response.ResponseHeaders["Server"] + response.ResponseHeaders["X-Powered-By"]
	doc, err := htmlquery.Parse(bytes.NewReader(response.Body))
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

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("forbidden redirects(10)")
	}
	return nil
}

func noCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 0 {
		return errors.New("forbidden redirects")
	}
	return nil
}

func getRespBody(resp *http.Response) ([]byte, error) {
	var body []byte
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gr, _ := gzip.NewReader(resp.Body)
		defer gr.Close()
		for {
			buf := make([]byte, 1024)
			n, err := gr.Read(buf)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if n == 0 {
				break
			}
			body = append(body, buf...)
		}
	} else {
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		body = raw
	}
	return body, nil
}
