package plugins

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"
)

func sshCheck(plg *Plugins, f *funcDesc) {
	crack(f.name, plg, f.user, f.pass, sshConn)
}

func sshConn(plg *Plugins, user string, pass string) (ok int) {
	ok = OKNext
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout:         time.Duration(plg.TimeOut) * time.Millisecond,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := net.DialTimeout("tcp4",
		fmt.Sprintf("%v:%v", plg.TargetIp, plg.TargetPort), time.Duration(plg.TimeOut)*time.Millisecond)
	if err != nil {
		return
	}
	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(time.Duration(plg.TimeOut) * time.Millisecond))
	if err != nil {
		return
	}
	clientConn, channelCh, reqCh, err := ssh.NewClientConn(conn, fmt.Sprintf("%v:%v", plg.TargetIp, plg.TargetPort), config)
	if err != nil {
		if strings.Contains(err.Error(), "password") {
			//密码错误
			return
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			//连接限制
			ok = OKWait
			return
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			ok = OKTimeOut
			return
		}
		//color.Red(err.Error() + plg.TargetIp + plg.TargetPort + user + pass)
		//协议异常
		ok = OKStop
		return
	}
	defer clientConn.Close()
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(plg.TimeOut) * time.Millisecond))
	if err != nil {
		return
	}
	client := ssh.NewClient(clientConn, channelCh, reqCh)
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(plg.TimeOut) * time.Millisecond))
	if err != nil {
		return
	}
	defer client.Close()
	session, err := client.NewSession()
	if err == nil {
		session.Close()
		ok = OKDone
	}
	return

}
