package plugins

import (
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"github.com/b1gcat/DarkEye/superscan/dic"
	"golang.org/x/net/context"
	"net"
	"strings"
	"time"
)

func smbCheck(s *Service) {
	s.crack()
}

func smbConn(parent context.Context, s *Service, user, pass string) (ok int) {
	ok = OKNext
	c := net.Dialer{Timeout: time.Duration(Config.TimeOut) * time.Millisecond}
	ctx, _ := context.WithCancel(parent)
	conn, err := c.DialContext(ctx, "tcp",
		fmt.Sprintf("%s:%s", s.parent.TargetIp, s.parent.TargetPort))
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			//连接限制
			ok = OKWait
			return
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			ok = OKTimeOut
			return
		}
		ok = OKStop
		return
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: pass,
		},
	}
	ctx, _ = context.WithTimeout(parent, time.Millisecond*time.Duration(Config.TimeOut))
	sb, err := d.DialContext(ctx, conn)
	if err != nil {
		ok = OKStop
		return
	}
	defer sb.Logoff()
	names, err := sb.ListSharenames()
	if err != nil {
		ok = OKStop
		return
	}
	ok = OKDone
	sharedDirectory := ""
	for _, v := range names {
		sharedDirectory += "," + v
	}
	strings.TrimPrefix(",", sharedDirectory)
	s.parent.Result.Output.Set("smb_shared", sharedDirectory)
	return
}

func init() {
	services["smb"] = Service{
		name:    "smb",
		port:    "445",
		user:    dic.DIC_USERNAME_SMB,
		pass:    dic.DIC_PASSWORD_SMB,
		check:   smbCheck,
		connect: smbConn,
		thread:  1,
	}
}
