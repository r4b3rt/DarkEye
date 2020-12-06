package plugins

import (
	"github.com/zsdevX/DarkEye/superscan/dic"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

var (
	mongoUsername = make([]string, 0)
	mongoPassword = make([]string, 0)
)

func init() {
	checkFuncs[MongoSrv] = mongoCheck
	mongoUsername = dic.DIC_USERNAME_MONGODB
	mongoPassword = dic.DIC_PASSWORD_MONGODB
}

func mongoCheck(plg *Plugins) {
	if !plg.NoTrust && plg.TargetPort != "27017" {
		return
	}
	if mongoUnAuth(plg) {
		plg.Cracked = append(plg.Cracked, Account{Username: "空", Password: "空"})
		plg.PortOpened = true
		plg.highLight = true
		plg.TargetProtocol = "mongodb"
		return
	}
	crack("mongodb", plg, mongoUsername, mongoPassword, mongodbConn)
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
