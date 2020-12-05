package plugins

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"time"
)

var (
	ftpUsername = make([]string, 0)
	ftpPassword = make([]string, 0)
)

func init() {
	checkFuncs[FtpSrv] = ftpCheck
	ftpUsername = loadDic("username_ftp.txt")
	ftpPassword = loadDic("password_ftp.txt")
}

func ftpCheck(plg *Plugins) {
	if !plg.NoTrust && plg.TargetPort != "21" {
		return
	}
	crack("ftp", plg, ftpUsername, ftpPassword, ftpConn)
}

func ftpConn(plg *Plugins, user string, pass string) (ok int) {
	ok = OKNext
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", plg.TargetIp, plg.TargetPort),
		time.Duration(plg.TimeOut)*time.Millisecond)
	if err == nil {
		err = conn.Login(user, pass)
		if err == nil {
			defer conn.Logout()
			ok = OKDone
		}
	} else {
		ok = OKStop
	}
	return ok
}
