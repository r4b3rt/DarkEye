package plugins

import (
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

func mongoCheck(plg *Plugins, f *funcDesc) {

	if mongoUnAuth(plg) {
		plg.Cracked = append(plg.Cracked, Account{Username: "空", Password: "空"})
		plg.PortOpened = true
		plg.highLight = true
		plg.TargetProtocol = f.name
		return
	}
	crack(f.name, plg, f.user, f.pass, mongodbConn)
}

func mongodbConn(plg *Plugins, user, pass string) (ok int) {
	ok = OKNext
	_, err := mgo.DialWithTimeout("mongodb://"+user+":"+pass+"@"+plg.TargetIp+":"+plg.TargetPort+"/"+"admin",
		time.Duration(plg.TimeOut)*time.Millisecond)
	if err == nil {
		ok = OKDone
	} else {
		if strings.Contains(err.Error(), "Authentication failed") {
			return
		}
		ok = OKStop
	}
	return
}

func mongoUnAuth(plg *Plugins) (ok bool) {
	session, err := mgo.DialWithTimeout(plg.TargetIp+":"+plg.TargetPort, time.Duration(plg.TimeOut)*time.Millisecond)
	if err == nil && session.Run("serverStatus", nil) == nil {
		ok = true
	}
	return ok
}
