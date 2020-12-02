package plugins

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"
)

var (
	sshUsername = make([]string, 0)
	sshPassword = make([]string, 0)
)

func sshCheck(plg *Plugins) interface{} {
	if !plg.NoTrust && plg.TargetPort != "22" {
		return nil
	}
	plg.SSh = make([]Account, 0)
	for _, user := range sshUsername {
		for _, pass := range sshPassword {
			pass = strings.Replace(pass, "%user%", user, -1)
			if ok, stop := sshConn(plg, user, pass); ok {
				plg.SSh = append(plg.SSh, Account{Username: user, Password: pass})
				plg.TargetProtocol = "[SSH]"
				return &plg.SSh[0]
			} else if stop {
				//非SSH协议退出
				return nil
			}
		}
	}
	//未找到密码
	plg.TargetProtocol = "[SSH]"
	plg.SSh = append(plg.SSh, Account{})
	return &plg.SSh[0]
}

func sshConn(plg *Plugins, user string, pass string) (ok bool, stop bool) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout: time.Duration(plg.TimeOut) * time.Millisecond,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
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
	clientConn, channelCh, reqCh, err := ssh.NewClientConn(conn, "tcp", config)
	if err != nil {
		if strings.Contains(err.Error(), "unable to authenticate") {
			//密码错误
			return
		}
		//协议异常
		stop = true
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
		ok = true
	}
	return

}

func init() {
	checkFuncs[SSHSrv] = sshCheck
	sshUsername = loadDic("username_ssh.txt")
	sshPassword = loadDic("password_ssh.txt")
}
