package scan

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
)

type memCacheConf struct {
	timeout  int
	username []string
	password []string
	logger   *logrus.Logger
}

func NewMemCache(timeout int) (Scan, error) {
	s := &memCacheConf{
		timeout:  timeout,
	}
	return s, nil
}

func (s *memCacheConf) Start(parent context.Context, ip, port string) (interface{}, error) {
	return fmt.Sprintf("memcached %s unauth", net.JoinHostPort(ip, port)), nil
}

func (s *memCacheConf) Setup(args ...interface{}) {
	s.username = args[0].([]string)
	s.password = args[1].([]string)
	s.logger = args[2].(*logrus.Logger)
}

func (s *memCacheConf) Identify(parent context.Context, ip, port string) bool {
	b, err := hello(parent, "tcp", net.JoinHostPort(ip, port), []byte("stats\n"), s.timeout)
	if err != nil {
		return false
	}
	return bytes.Contains(b, []byte("STAT"))
}

func (s *memCacheConf) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *memCacheConf) Output() interface{} {
	return nil
}
