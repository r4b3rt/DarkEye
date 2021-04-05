package plugins

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"time"
)


func ftpCheck(plg *Plugins, f *funcDesc) {
	crack(f.name, plg, f.user, f.pass, ftpConn)
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
