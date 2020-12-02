package plugins

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"sync"
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
	wg := sync.WaitGroup{}
	wg.Add(len(sshUsername))
	var once sync.Once
	limiter := make(chan int, plg.Worker)
	ctx, cancel := context.WithCancel(context.TODO())
	for _, user := range sshUsername {
		limiter <- 1
		go func(username string) {
			defer func() {
				<-limiter
				wg.Done()
			}()
			for _, pass := range sshPassword {
				select {
				case <-ctx.Done():
					return
				default:
				}
				pass = strings.Replace(pass, "%user%", username, -1)
				plg.DescCallback(fmt.Sprintf("Cracking ssh %s:%s %s/%s",
					plg.TargetIp, plg.TargetPort, username, pass))
				ok := sshConn(plg, username, pass)
				switch ok {
				case OKDone:
					//密码正确一次退出
					plg.locker.Lock()
					plg.SSh = append(plg.SSh, Account{Username: username, Password: pass})
					plg.locker.Unlock()
					plg.highLight = true
					once.Do(func() { cancel() })
					return
				case OKWait:
					//太快了服务器限制
					color.Red("[ssh]爆破频率太快服务器受限，建议降低参数'plugin-worker'数值影响主机:%s:%s",
						plg.TargetIp, plg.TargetPort)
					once.Do(func() { cancel() })
					return
				case OKTimeOut:
					color.Red("[ssh]爆破过程中连接超时，建议提高参数'timeout'数值影响主机:%s:%s",
						plg.TargetIp, plg.TargetPort)
					once.Do(func() { cancel() })
					return
				case OKStop:
					//非协议退出
					once.Do(func() { cancel() })
					return
				default:
					//密码错误.OKNext
					plg.TargetProtocol = "[ssh]"
				}
			}
		}(user)
	}
	wg.Wait()
	//未找到密码
	if plg.TargetProtocol != "" {
		return &plg.SSh
	}
	return nil
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

func init() {
	checkFuncs[SSHSrv] = sshCheck
	sshUsername = loadDic("username_ssh.txt")
	sshPassword = loadDic("password_ssh.txt")
}
