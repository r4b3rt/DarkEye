package scan

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	urlpkg "net/url"
	"strings"
)

type cms struct {
	Favicon string
	Cms     string
}
type httpDisco struct {
	Host        string
	Server      string
	StatusCode  int
	Title       string
	Url         string
	RedirectUrl []string
	Cms         []cms
}

func (s *discovery) http(ctx context.Context, ip, port string) (interface{}, error) {
	hosts := make([]string, 0)
	hosts = append(hosts, ip)
	if s.host != nil {
		hosts = append(hosts, s.host...)
	}
	//ip host check
	r := make([]httpDisco, 0)
	for _, h := range hosts {
		host := net.JoinHostPort(h, port)
		url, err := s.httpIdent(ctx, host, ip, port)
		if err != nil {
			s.logger.Debug("http.httpIdent:", err.Error())
			continue
		}
		disco := httpDisco{
			Host:        h,
			RedirectUrl: make([]string, 0),
			Cms:         make([]cms, 0),
		}
		x, _ := s.httpFetch(ctx, host, url, &disco, true)
		if x != nil {
			r = append(r, disco)
		}
	}
	switch len(r) {
	default:
		return r, nil
	case 0:
		return nil, nil
	}
}

func (s *discovery) httpFetch(ctx context.Context, host string, test *urlpkg.URL, disco *httpDisco, redirect bool) (interface{}, error) {
	client := newHttpClient(s.timeout, disco)
	request, err := http.NewRequestWithContext(ctx, "GET", test.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Host = host //very important
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
	} else {
		defaultFavicon := test.String() + "/favicon.ico"
		s.findFavicon(ctx, host, test, defaultFavicon, disco)
		if favicon := getFavicon(body); favicon != nil {
			for _, fav := range favicon {
				s.findFavicon(ctx, host, test, fav[1], disco)
			}
		}
	}

	//parse title
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		s.logger.Debug("html.query.Parse:", err.Error())
		return disco, nil
	}
	t := htmlquery.FindOne(doc, "//title")

	if t != nil {
		disco.Title = htmlquery.InnerText(t)
		if !isUtf8([]byte(disco.Title)) {
			if message, err := simplifiedchinese.GBK.NewDecoder().String(disco.Title); err == nil {
				disco.Title = message
			}
		}
	}

	if !redirect {
		return disco, nil
	}

	test, err = s.findRefreshInMeta(test, doc)
	if err == nil {
		return s.httpFetch(ctx, host, test, disco, true)
	}
	s.logger.Debug("findRefreshInMeta:", err.Error())
	return disco, nil
}

func (s *discovery) findFavicon(ctx context.Context, host string, test *urlpkg.URL, favicon string, disco *httpDisco) {
	cli := newHttpClient(s.timeout, disco)
	cli.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	url := test.String()
	if strings.HasPrefix(favicon, "http") {
		url = favicon
	} else {
		url = test.String() + favicon
	}
	s.logger.Debug("getting favicon:", url)
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		s.logger.Debug("findFavicon.NewRequestWithContext:", err.Error())
		return
	}
	request.Host = host
	request.Header.Add("Accept-Encoding", "xZip") //
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	response, err := cli.Do(request)
	if err != nil {
		s.logger.Debug("findFavicon.do:", err.Error())
		return
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			s.logger.Debug("findFavicon.ReadAll:", err.Error())
			return
		}
		fh := mmh3Hash32(standBase64(body))
		finger, _ := faviconHash[fh]
		disco.Cms = append(disco.Cms, cms{fh, finger})
	}
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
func (s *discovery) httpIdent(ctx context.Context, host, ip, port string) (*urlpkg.URL, error) {
	test := &urlpkg.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(ip, port),
	}
	ok, err := s._httpIdent(ctx, host, test)
	if ok {
		return test, nil
	}

	s.logger.Debug("httpIdent.http->https.changed:", err)

	//if error fallthrough ...
	test.Scheme = "https"
	if ok, err = s._httpIdent(ctx, host, test); ok {
		return test, nil
	}
	return nil, fmt.Errorf("not a http or https")
}

func (s *discovery) _httpIdent(ctx context.Context, host string, test *urlpkg.URL) (ok bool, err error) {
	cli := newHttpClient(s.timeout, nil)
	req, err := http.NewRequestWithContext(ctx, "GET", test.String(), nil)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	req.Host = host
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
