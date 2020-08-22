package common

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Http struct {
	Method      string
	Url         string
	Data        []byte
	Cookie      string
	Origin      string
	ContentType string
	Referer     string
	Agent       string
	TimeOut     time.Duration
}

func (m *Http) Http() ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli := http.Client{
		Transport: tr,
	}
	req, _ := http.NewRequest(m.Method, m.Url, bytes.NewReader(m.Data))
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
	req.Header.Add("Cookie", m.Cookie)

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(fmt.Sprintf("Bad status %d", resp.StatusCode))
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func IsAlive(ip, port string) bool {
	c, err := net.DialTimeout("tcp", ip+":"+port, time.Second*2)
	if err != nil {
		return false
	}
	defer c.Close()
	return true
}
