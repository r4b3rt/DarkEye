package scan

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	urlpkg "net/url"
	"strings"
)

type httpDisco struct {
	Server     string
	StatusCode int
	Title      string
	Url        string
}

func (s *discovery) http(ctx context.Context, ip, port string) (interface{}, error) {
	url, err := s.httpIdent(ctx, ip, port)
	if err != nil {
		return nil, err
	}
	disco := &httpDisco{}
	return s.httpFetch(ctx, url, disco, true)
}

func (s *discovery) httpFetch(ctx context.Context, test *urlpkg.URL, disco *httpDisco, redirect bool) (interface{}, error) {
	client := newHttpClient(s.timeout)
	request, err := http.NewRequestWithContext(ctx, "GET", test.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Host = test.Host //very important
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	disco.StatusCode = response.StatusCode
	disco.Url = test.String()

	for k := range response.Header {
		if k == "Server" || k == "X-Powered-By" {
			disco.Server = response.Header.Get(k)
		}
	}

	var body []byte

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(response.Body)
		for {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)

			if err != nil && err != io.EOF {
				return &disco, err
			}

			if n == 0 {
				break
			}
			body = append(body, buf[:n]...)
		}
	default:
		body, _ = ioutil.ReadAll(response.Body)
	}

	if body == nil {
		return disco, nil
	}

	//parse title
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		s.logger.Debug("htmlquery.Parse:", err.Error())
		return disco, nil
	}
	t := htmlquery.FindOne(doc, "//title")

	if t != nil {
		disco.Title = htmlquery.InnerText(t)
	}

	if !redirect {
		return disco, nil
	}

	test, err = s.findRefreshInMeta(test, doc)
	if err == nil {
		return s.httpFetch(ctx, test, disco, true)
	}
	s.logger.Debug("findRefreshInMeta:", err.Error())
	return disco, nil
}

func (s *discovery) findRefreshInMeta(old *urlpkg.URL, doc *html.Node) (*urlpkg.URL, error) {
	//parse refresh
	//<meta http-equiv="refresh" content="0;URL='https://example.com/'">
	//<meta http-equiv="refresh" content="0;/beian/">
	meta := htmlquery.Find(doc, "//meta")

	for _, v := range meta {
		match := false
		for _, a := range v.Attr {
			if !match &&
				(strings.ToLower(a.Key) != "http-equiv" || strings.ToLower(a.Val) != "refresh") {
				break
			}
			match = true
			if strings.ToLower(a.Key) != "content" {
				continue
			}
			val := strings.Split(a.Val, ";")
			if len(val) != 2 {
				return nil, fmt.Errorf("unknown refresh content format")
			}
			val2 := strings.Split(val[1], "=")
			switch len(val2) {
			case 1:
				old.Path = val2[0]
			case 2:
				var err error
				old, err = urlpkg.Parse(val2[1])
				if err != nil {
					return nil, fmt.Errorf("bad content-url format")
				}
			}
			return old, nil
		}
	}
	return nil, fmt.Errorf("not found refresh in meta")
}

//httpIdent return url, host, error
func (s *discovery) httpIdent(ctx context.Context, ip, port string) (*urlpkg.URL, error) {
	test := &urlpkg.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(ip, port),
	}
	ok, err := s._httpIdent(ctx, test)
	if ok {
		return test, nil
	}

	s.logger.Debug("httpIdent.http->https.changed:", err)

	//if error fallthrough ...
	test.Scheme = "https"
	if ok, err = s._httpIdent(ctx, test); ok {
		return test, nil
	}
	return nil, fmt.Errorf("not a http or https")
}

func (s *discovery) _httpIdent(ctx context.Context, test *urlpkg.URL) (ok bool, err error) {
	cli := newHttpClient(s.timeout)
	req, err := http.NewRequestWithContext(ctx, "GET", test.String(), nil)
	if err != nil {
		return
	}
	resp, err := cli.Do(req)
	if err != nil {
		s.logger.Debug("_httpIdent:", err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		if bytes.Contains(body, []byte("The plain HTTP request was sent to HTTPS port")) {
			return ok, fmt.Errorf("maybe access https by http")
		}
	}
	return true, nil
}
