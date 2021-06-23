package plugins

import (
	"context"
	"fmt"
	"github.com/b1gcat/DarkEye/common"
	"github.com/b1gcat/DarkEye/superscan/dic"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"
)

func sshCheck(s *Service) {
	s.crack()
}

func sshConn(parent context.Context, s *Service, user, pass string) (ok int) {
	ok = OKNext

	timeOut := time.Millisecond * time.Duration(Config.TimeOut)
	//初始化连接
	conn, err := common.DialCtx(parent, "tcp",
		net.JoinHostPort(s.parent.TargetIp, s.parent.TargetPort), timeOut)
	if err != nil {
		//网络不通或墙了
		ok = OKTerm
		return
	}
	defer conn.Close()
	config := &ssh.ClientConfig{
		User:    user,
		Auth:    []ssh.AuthMethod{ssh.Password(pass)},
		Timeout: timeOut,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	_ = conn.SetReadDeadline(time.Now().Add(time.Millisecond*1000))
	c, ch, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(s.parent.TargetIp, s.parent.TargetPort), config)
	if err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "password") {
			//密码错误
			return
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			return
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			return
		}
		return
	}
	ok = OKDone
	_ = conn.SetReadDeadline(time.Now().Add(timeOut))
	client := ssh.NewClient(c, ch, reqs)
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()
	out, err := session.CombinedOutput("id")
	if err == nil {
		s.parent.Result.Output.Set("helper", string(out))
	}
	return
}

func init() {
	services["ssh"] = Service{
		name:    "ssh",
		port:    "22",
		user:    dic.DIC_USERNAME_SSH,
		pass:    dic.DIC_PASSWORD_SSH,
		check:   sshCheck,
		connect: sshConn,
		thread:  3,
	}
}
