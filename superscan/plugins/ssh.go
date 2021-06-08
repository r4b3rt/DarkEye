package plugins

import "C"
import (
	"context"
	"fmt"
	"github.com/melbahja/goph"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"golang.org/x/crypto/ssh"
	"net"
	"strconv"
	"strings"
	"time"
)

func sshCheck(s *Service) {
	s.crack()
}

func sshConn(parent context.Context, s *Service, user, pass string) (ok int) {
	ok = OKNext
	port, _ := strconv.Atoi(s.parent.TargetPort)
	//初始化变量
	client := &goph.Client{
		Config: &goph.Config{
			User:     user,
			Addr:     s.parent.TargetIp,
			Port:     uint(port),
			Auth:     goph.Password(pass),
			Timeout:  time.Duration(Config.TimeOut) * time.Millisecond,
			Callback: ssh.InsecureIgnoreHostKey(),
		},
	}
	//初始化连接
	conn, err := common.DialCtx(parent, "tcp",
		net.JoinHostPort(client.Config.Addr, fmt.Sprint(client.Config.Port)),
		client.Config.Timeout)
	if err != nil {
		//网络不通或墙了
		ok = OKTerm
		return
	}
	config := &ssh.ClientConfig{
		User:            client.Config.User,
		Auth:            client.Config.Auth,
		Timeout:         client.Config.Timeout,
		HostKeyCallback: client.Config.Callback,
	}
	_ = conn.SetReadDeadline(time.Now().Add(client.Config.Timeout))
	c, ch, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(client.Config.Addr, fmt.Sprint(client.Config.Port)), config)
	if err != nil {
		if strings.Contains(err.Error(), "password") {
			//密码错误
			return
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			return
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			ok = OKTerm
			return
		}
		return
	}
	_ = conn.SetReadDeadline(time.Now().Add(client.Config.Timeout))
	client.Client = ssh.NewClient(c, ch, reqs)
	defer client.Close()
	out, err := client.Run("id")
	if err == nil {
		s.parent.Result.Output.Set("helper", string(out))
	}
	ok = OKDone
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
		thread:  1,
	}
}
