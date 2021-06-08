package plugins

import (
	"context"
	"github.com/b1gcat/DarkEye/common"
	"net"
	"strings"
	"time"
)

func memCachedCheck(s *Service) {
	if memCachedUnAuth(s) {
		s.parent.Result.Output.Set("account", `空, 空`)
	}
	s.parent.Result.ServiceName = s.name
	s.parent.Hit = true
	return
}

func memCacheConn(parent context.Context, s *Service, user, pass string) (ok int) {
	return OKStop
}

func memCachedUnAuth(s *Service) (ok bool) {
	conn, err := common.DialCtx(context.Background(), "tcp",
		net.JoinHostPort(s.parent.TargetIp, s.parent.TargetPort), time.Duration(Config.TimeOut)*time.Millisecond)
	if err != nil {
		//网络不通或墙了
		ok = false
		return
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(time.Duration(Config.TimeOut) * time.Millisecond))
	_, _ = conn.Write([]byte("stats\n"))
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		return
	}
	if strings.Contains(string(buff[:n]), "STAT") {
		s.parent.Result.Output.Set("helper",
			"apt install libmemcached-tools\nmemcdump --servers=192.168.1.33 (列出key）\nmemccat --servers=192.168.1.33 key1 （列出key1内容）\n")
		return true
	}
	return
}

func init() {
	services["memcached"] = Service{
		name:    "memcached",
		port:    "11211",
		check:   memCachedCheck,
		connect: memCacheConn,
		thread:  1,
	}
}
