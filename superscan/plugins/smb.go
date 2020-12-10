package plugins
//没找到特别好的库
import (
	"github.com/stacktitan/smb/smb"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"strconv"
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
	port, _ := strconv.Atoi(plg.TargetPort)
	options := smb.Options{
		Host:        plg.TargetIp,
		Port:        port,
		User:        user,
		Password:    pass,
		Domain:      "",
		Workstation: "",
		TimeOut:     time.Duration(plg.TimeOut) * time.Millisecond,
	}
	session, err := smb.NewSession(options, false)
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
	defer session.Close()
	if session.IsAuthenticated {
		ok = OKDone
	}
	return
}
