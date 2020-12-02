package plugins

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func mysqlCheck(plg *Plugins) interface{} {
	if !plg.NoTrust && plg.TargetPort != "3306" {
		return nil
	}
	plg.SSh = make([]Account, 0)
	for _, user := range mysqlUsername {
		for _, pass := range mysqlPassword {
			pass = strings.Replace(pass, "%user%", user, -1)
			if ok, stop := MysqlConn(plg, user, pass); ok {
				plg.SSh = append(plg.Mysql, Account{Username: user, Password: pass})
				plg.TargetProtocol = "[Mysql]"
				return &plg.SSh[0]
			} else if stop {
				//非SSH协议退出
				return nil
			}
		}
	}
	//未找到密码
	plg.TargetProtocol = "[Account]"
	plg.SSh = append(plg.Mysql, Account{})
	return &plg.SSh[0]
}

func MysqlConn(plg *Plugins, user string, pass string) (ok bool, stop bool) {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", user, pass, plg.TargetIp, plg.TargetPort, "mysql"))
	db.SetConnMaxLifetime(time.Duration(plg.TimeOut) * time.Millisecond)
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			ok = true
		}
	} else {
		//非Mysql协议或受限制
		stop = true
	}
	return
}

var (
	mysqlUsername = make([]string, 0)
	mysqlPassword = make([]string, 0)
)

func init() {
	checkFuncs[MysqlSrv] = mysqlCheck
	mysqlUsername = loadDic("username_mysql.txt")
	mysqlPassword = loadDic("password_mysql.txt")
}
