package plugins

import (
	"github.com/hirochachacha/go-smb2"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"golang.org/x/net/context"
	"net"
	"strings"
	"time"
)

var (
	smbUsername = make([]string, 0)
	smbPassword = make([]string, 0)
)

func init() {
	checkFuncs[SmbSrv] = msbCheck
	smbUsername = dic.DIC_USERNAME_SMB
	smbPassword = dic.DIC_PASSWORD_SMB
	supportPlugin["smb"] = "smb"
}

func msbCheck(plg *Plugins) {
	if !plg.NoTrust && plg.TargetPort != "445" {
		return
	}
	crack("smb", plg, smbUsername, smbPassword, smbConn)
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
	ck := Account{}
	ck.Shares = names
	plg.Lock()
	plg.Cracked = append(plg.Cracked, ck)
	plg.Unlock()
	return
}
