package scan

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jlaffaye/ftp"
	"net"
)

type ftpConf struct {
	timeout  int
	risk
}

func NewFtp(timeout int) (Scan, error) {
	s := &ftpConf{
		timeout:  timeout,
	}

	return s, nil
}

func (s *ftpConf) Start(parent context.Context, ip, port string) (interface{}, error) {
	addr := net.JoinHostPort(ip, port)
	return weakPass(parent, "ftp", addr, s.username, s.password, s.crack)
}

func (s *ftpConf) Setup(args ...interface{}) {
	setupRisk(&s.risk, args)
}

func (s *ftpConf) crack(parent context.Context, addr, user, pass string) bool {
	c, err := ftp.Dial(addr, ftp.DialWithContext(parent))
	if err != nil {
		s.logger.Debug("ftpConf.dail:", err.Error())
		return false
	}
	if err= c.Login(user, pass); err != nil {
		s.logger.Debug("ftpConf.dail:", err.Error())
		return false
	}
	return true
}

func (s *ftpConf) Identify(parent context.Context, ip, port string) bool {
	b, err := hello(parent, "tcp", net.JoinHostPort(ip, port), nil, s.timeout)
	if err != nil {
		return false
	}
	return bytes.Contains(b, []byte("FTP"))
}

func (s *ftpConf) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *ftpConf) Output() interface{} {
	return nil
}
