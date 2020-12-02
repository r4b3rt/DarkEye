package plugins

import (
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	"strings"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

func mysqlCheck(plg *Plugins) interface{} {
	if !plg.NoTrust && plg.TargetPort != "3306" {
		return nil
	}
	plg.SSh = make([]Account, 0)
L:
	for _, user := range mysqlUsername {
		for _, pass := range mysqlPassword {
			pass = strings.Replace(pass, "%user%", user, -1)
			if ok, stop := MysqlConn(plg, user, pass); ok {
				plg.SSh = append(plg.Mysql, Account{Username: user, Password: pass})
				plg.TargetProtocol = "[Mysql]"
				plg.highLight = true
				return &plg.SSh[0]
			} else if stop {
				if ok {
					//目标协议正确但是做了外网限制禁止登录
					break L
				}
				//非协议退出
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
	if err != nil {
		color.Red(err.Error())
		if strings.Contains(err.Error(), "password") {
			return
		}
		if strings.Contains(err.Error(), "not allowed to connect") {
			ok = true
			return
		}
		//非Mysql协议或受限制
		stop = true
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Duration(plg.TimeOut) * time.Millisecond)
	err = db.Ping()
	if err == nil {
		ok = true
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
