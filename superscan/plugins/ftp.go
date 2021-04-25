package plugins

import (
	"context"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"net"
	"time"
)

func ftpCheck(s *Service) {
	s.crack()
}

func ftpConn(parent context.Context, s *Service, user, pass string) (ok int) {
	ok = OKNext

	c, err := common.DialCtx(parent, "tcp",
		net.JoinHostPort(s.parent.TargetIp, s.parent.TargetPort), time.Duration(Config.TimeOut)*time.Millisecond)
	if err != nil {
		//网络不通或墙了
		ok = OKTerm
		return
	}
	_ = c.SetReadDeadline(time.Now().Add(time.Duration(Config.TimeOut)*time.Millisecond))
	conn, err := ftp.Dial(net.JoinHostPort(s.parent.TargetIp, s.parent.TargetPort), ftp.DialWithNetConn(c))
	if err != nil {
		fmt.Println(err.Error())
		ok = OKStop
		return
	}
	err = conn.Login(user, pass)
	if err == nil {
		defer conn.Logout()
		ok = OKDone
	}
	return
}

func init() {
	services["ftp"] = Service{
		name:    "ftp",
		port:    "21",
		user:    dic.DIC_USERNAME_FTP,
		pass:    dic.DIC_PASSWORD_FTP,
		check:   ftpCheck,
		connect: ftpConn,
		thread:  1,
	}
}
