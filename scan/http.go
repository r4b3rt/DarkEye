package scan

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type httpConf struct {
	timeout  int
	username []string
	password []string
	logger   *logrus.Logger
}

func NewHttp(timeout int, args []interface{}) (Scan, error) {
	s := &httpConf{
		timeout:  timeout,
		username: args[0].([]string),
		password: args[1].([]string),
		logger:   args[2].(*logrus.Logger),
	}

	return s, nil
}

func (s *httpConf) Start(parent context.Context, ip, port string) (interface{}, error) {
	addr := net.JoinHostPort(ip, port)
	return weakPass(parent, "ssh", addr, s.username, s.password, s.crack)
}

func (s *httpConf) crack(_ context.Context, addr, user, pass string) bool {
	timeOut := time.Millisecond * time.Duration(s.timeout)
	config := &ssh.ClientConfig{
		User:    user,
		Auth:    []ssh.AuthMethod{ssh.Password(pass)},
		Timeout: timeOut,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	cli, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return false
	}
	defer cli.Close()
	return true
}

func (s *httpConf) Identify(parent context.Context, ip, port string) bool {
	b, err := hello(parent, "tcp", net.JoinHostPort(ip, port), nil, s.timeout)
	if err != nil {
		return false
	}
	return bytes.Contains(b, []byte("SSH-"))
}

func (s *httpConf) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *httpConf) Output() interface{} {
	return nil
}
