package scan

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

func dail(parent context.Context, protocol, addr string, timeout int) (net.Conn, error) {
	timeOut := time.Millisecond * time.Duration(timeout)
	d := net.Dialer{Timeout: timeOut}
	ctx, _ := context.WithTimeout(parent, timeOut)
	return d.DialContext(ctx, protocol, addr)
}

func hello(parent context.Context, protocol, addr string, hi []byte, timeout int) ([]byte, error) {
	c, err := dail(parent, protocol, addr, timeout)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	if hi != nil {
		_, _ = c.Write(hi)
	}
	_ = c.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeout)))
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func newHttpClient(timeout int) *http.Client {
	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Millisecond,
			KeepAlive: -1,
		}).DialContext,
		ResponseHeaderTimeout: time.Duration(timeout) * time.Millisecond,
		TLSHandshakeTimeout:   time.Duration(timeout) * time.Millisecond,
	}

	return &http.Client{
		Transport: tr,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("forbidden redirects(10)")
			}
			return nil
		},
	}
}
