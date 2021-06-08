package common

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	//Die add comment
	Die = 1
	//Alive add comment
	Alive = 2
	//TimeOut add comment
	TimeOut = 3
)

//HttpRequest add comment
type HttpRequest struct {
	Method  string
	Url     string
	Body    []byte
	Headers map[string]string
	Proxy   string
	Ctx     context.Context

	NoFollowRedirect bool
	TimeOut          time.Duration //秒
}

//HttpResponse add comment
type HttpResponse struct {
	Status          int32
	ResponseHeaders map[string]string
	Body            []byte
	ContentType     string
}

//Go add comment
func (m *HttpRequest) Go() (*HttpResponse, error) {
	var proxy func(r *http.Request) (i *url.URL, e error)
	if m.Proxy != "" {
		proxy = func(request *http.Request) (i *url.URL, e error) {
			return url.Parse(m.Proxy)
		}
	}
	tr := &http.Transport{
		Proxy:           proxy,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   m.TimeOut * time.Second,
			KeepAlive: m.TimeOut * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: time.Second * m.TimeOut,
		TLSHandshakeTimeout:   time.Second * m.TimeOut,
	}
	jar, _ := cookiejar.New(nil)
	cli := http.Client{
		Transport:     tr,
		CheckRedirect: defaultCheckRedirect,
		Jar:           jar,
		//Timeout:       m.TimeOut * time.Second,
	}
	if m.NoFollowRedirect {
		cli.CheckRedirect = noCheckRedirect
	}
	if m.Ctx == nil {
		m.Ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(m.Ctx, m.Method, m.Url, bytes.NewReader(m.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range m.Headers {
		if strings.ToLower(k) == "host" {
			req.Host = v
			continue
		}
		req.Header[k] = []string{v}
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
		if k != "Set-Cookie" {
			response.ResponseHeaders[k] = resp.Header.Get(k)
		}
	}
	for _, ck := range resp.Cookies() {
		response.ResponseHeaders["Set-Cookie"] += ck.String() + ";"
	}
	body, err := getRespBody(resp)
	if err != nil {
		return &response, err
	}
	response.Body = body
	return &response, nil
}

//WhatWeb add comment
func WhatWeb(ctx context.Context, proto, target, host string, timeOutSec int) (server, title string, code int32, finger string) {
	req := HttpRequest{
		Url:     fmt.Sprintf(proto+"://%s", target),
		TimeOut: time.Duration(timeOutSec),
		Method:  "GET",
		Headers: map[string]string{
			"User-Agent": UserAgents[0],
			"Host":       host,
		},
		Ctx: ctx,
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
	code = response.Status
	t := htmlquery.Find(doc, "//title")
	if len(t) != 0 {
		title = htmlquery.InnerText(t[0])
	}
	if !ISUtf8([]byte(title)) {
		if message, err := simplifiedchinese.GBK.NewDecoder().String(title); err == nil {
			title = message
		}
	}
	title = TrimLRS.ReplaceAllString(title, "")
	//获取finger
	//todo：fast match？
	headerStr, err := json.Marshal(&response.ResponseHeaders)
	if err != nil {
		Log("whatWeb.Marshal", err.Error(), FAULT)
		return
	}
	finger = getFinger(headerStr, response.Body)
	return
}

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("forbidden redirects(10)")
	}
	return nil
}

func noCheckRedirect(_ *http.Request, via []*http.Request) error {
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

func getFinger(header, body []byte) (finger string) {
	match := true
	for _, f := range httpFinger.Fingerprint {
		match = true
		switch f.Method {
		case "keyword":
			buf := header
			if f.Location == "body" {
				buf = body
			}
			for _, k := range f.Keyword {
				if !bytes.Contains(buf, []byte(k)) {
					match = false
					break
				}
			}
			if match {
				finger = f.Cms
				return
			}
		}
	}
	return
}
