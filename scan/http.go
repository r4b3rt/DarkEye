package scan

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	urlpkg "net/url"
)

type httpDisco struct {
	Server     string
	StatusCode int
	Title      string
	Url        string
}

func (s *discovery) http(ctx context.Context, ip, port string) (interface{}, error) {
	url, host, err := s.httpIdent(ctx, ip, port)
	if err != nil {
		return nil, err
	}

	client := newHttpClient(s.timeout)
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

	disco := httpDisco{
		StatusCode: response.StatusCode,
		Url:        url,
	}

	for k := range response.Header {
		if k == "Server" || k == "X-Powered-By" {
			disco.Server = response.Header.Get(k)
		}
	}

	var body []byte
	switch response.StatusCode {
	default:
		return &disco, nil
	case 200:
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
	}

	return nil, nil
}

//httpIdent return url, host, error
func (s *discovery) httpIdent(ctx context.Context, ip, port string) (url string, host string, err error) {
	test := urlpkg.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(ip, port),
	}
	var ok bool
	if ok, err = s._httpIdent(ctx, &test); ok {
		return test.String(), test.Host, nil
	}
	//if error fallthrough ...
	test.Scheme = "https"
	if ok, err = s._httpIdent(ctx, &test); ok {
		return test.String(), test.Host, nil
	}
	return
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
		if !bytes.Contains(body, []byte("The plain HTTP request was sent to HTTPS port")) {
			return ok, nil
		}
	}
	return
}
