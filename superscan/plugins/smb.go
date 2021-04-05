package plugins

import (
	"github.com/hirochachacha/go-smb2"
	"golang.org/x/net/context"
	"net"
	"strings"
	"time"
)

func msbCheck(plg *Plugins, f *funcDesc) {
	crack(f.name, plg, f.user, f.pass, smbConn)
}

func smbConn(plg *Plugins, user, pass string) (ok int) {
	ok = OKNext
	conn, err := net.DialTimeout("tcp", plg.TargetIp+":"+plg.TargetPort,
		time.Duration(plg.TimeOut)*time.Millisecond)
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
	ctx, _ := context.WithTimeout(context.TODO(), time.Millisecond*time.Duration(plg.TimeOut))
	s, err := d.DialContext(ctx, conn)
	if err != nil {
		ok = OKStop
		return
	}
	defer s.Logoff()
	names, err := s.ListSharenames()
	if err != nil {
		ok = OKStop
		return
	}
	ok = OKDone
	for _, v := range names {
		plg.NetBios.Shares += "," + v
	}
	strings.TrimPrefix(",", plg.NetBios.Shares)
	return
}
