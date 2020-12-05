package plugins

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

var (
	mysqlUsername = make([]string, 0)
	mysqlPassword = make([]string, 0)
)

func init() {
	checkFuncs[MysqlSrv] = mysqlCheck
	mysqlUsername = loadDic("username_mysql.txt")
	mysqlPassword = loadDic("password_mysql.txt")
	_ = mysql.SetLogger(mysqlLogger(&mysqlNoLogger{}))
}

func mysqlCheck(plg *Plugins) {
	if !plg.NoTrust && plg.TargetPort != "3306" {
		return
	}
	crack("mysql", plg, mysqlUsername, mysqlPassword, mysqlConn)
}

func mysqlConn(plg *Plugins, user string, pass string) (ok int) {
	plg.RateWait(plg.RateLimiter) //爆破限制
	ok = OKNext
	source := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?timeout=%dms&readTimeout=%dms",
		user, pass, plg.TargetIp, plg.TargetPort, "mysql", plg.TimeOut, plg.TimeOut)
	db, err := sql.Open("mysql", source)
	if err != nil {
		if strings.Contains(err.Error(), "password") {
			return
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			//防火墙连接限制
			ok = OKWait
			return
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			ok = OKTimeOut
			return
		}
		if strings.Contains(err.Error(), "not allowed to connect") {
			//Mysql配置限制
			ok = OKForbidden
			return
		}
		//非协议
		ok = OKStop
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Duration(plg.TimeOut) * time.Millisecond)
	err = db.Ping()
	if err == nil {
		ok = OKDone
	} else {
		//MariaDB Open时候会返回正确但这里返回错误
	}
	return
}

//接管mysql的垃圾日志
type mysqlLogger interface {
	Print(v ...interface{})
}

type mysqlNoLogger struct {
}

func (*mysqlNoLogger) Print(v ...interface{}) {
}
