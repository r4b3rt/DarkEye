package scan

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"net"
	"strings"
	"time"
)

type redisConf struct {
	timeout  int
	risk
}

func NewRedis(timeout int) (Scan, error) {
	s := &redisConf{
		timeout:  timeout,
	}

	return s, nil
}

func (s *redisConf) Start(parent context.Context, ip, port string) (interface{}, error) {
	addr := net.JoinHostPort(ip, port)
	return weakPass(parent, "redis", addr, s.username, s.password, s.crack)
}

func (s *redisConf) Setup(args ...interface{}) {
	setupRisk(&s.risk, args)
}

func (s *redisConf) crack(parent context.Context, addr, user, pass string) bool {
	timeOut := time.Millisecond * time.Duration(s.timeout)

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    pass,
		DB:          0,
		DialTimeout: timeOut,
	})
	ctx, _ := context.WithCancel(parent)

	cli := client.WithContext(ctx)
	r, err := cli.Ping().Result()
	if err != nil {
		s.logger.Debug("redis.ping:", err.Error())
		return false
	}
	if strings.Contains(r, "PONG") {
		return true
	}
	return false
}

func (s *redisConf) Identify(parent context.Context, ip, port string) bool {
	b, err := hello(parent, "tcp", net.JoinHostPort(ip, port), []byte("fuck\n"), s.timeout)
	if err != nil {
		s.logger.Debug("redisConf.Identify:", err.Error())
		return false
	}
	return bytes.Contains(b, []byte("-ERR unknown command"))
}

func (s *redisConf) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *redisConf) Output() interface{} {
	return nil
}
