package scan

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type sshConf struct {
	timeout  int
	risk
}

func NewSSh(timeout int) (Scan, error) {
	s := &sshConf{
		timeout: timeout,
	}
	return s, nil
}

func (s *sshConf) Start(parent context.Context, ip, port string) (interface{}, error) {
	addr := net.JoinHostPort(ip, port)
	return weakPass(parent, "ssh", addr, s.username, s.password, s.crack)
}

func (s *sshConf) Setup(args ...interface{}) {
	setupRisk(&s.risk, args)
}

func (s *sshConf) crack(_ context.Context, addr, user, pass string) bool {
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
		s.logger.Debug("sshConf.Dial:", err.Error())
		return false
	}
	defer cli.Close()
	return true
}

func (s *sshConf) Identify(parent context.Context, ip, port string) bool {
	b, err := hello(parent, "tcp", net.JoinHostPort(ip, port), nil, s.timeout)
	if err != nil {
		return false
	}
	return bytes.Contains(b, []byte("SSH-"))
}

func (s *sshConf) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *sshConf) Output() interface{} {
	return nil
}
