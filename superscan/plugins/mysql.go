package plugins

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func mysqlCheck(s *Service) {
	s.crack()
}

func mysqlConn(_ context.Context, s *Service, user, pass string)  (ok int) {
	ok = OKNext
	source := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?timeout=%dms&readTimeout=%dms",
		user, pass, s.parent.TargetIp, s.parent.TargetPort, "mysql", Config.TimeOut, Config.TimeOut)
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
	db.SetConnMaxLifetime(time.Duration(Config.TimeOut) * time.Millisecond)
	err = db.Ping()
	if err == nil {
		ok = OKDone
	} else {
		//MariaDB Open时候会返回正确但这里返回错误
	}
	return
}

func init() {
	services["mysql"] = Service{
		name:    "mysql",
		port:    "3306",
		user:    dic.DIC_USERNAME_MYSQL,
		pass:    dic.DIC_PASSWORD_MYSQL,
		check:   mysqlCheck,
		connect: mysqlConn,
		thread:  1,
	}
}

//接管mysql的垃圾日志
type mysqlLogger interface {
	Print(v ...interface{})
}

type mysqlNoLogger struct {
}

func (*mysqlNoLogger) Print(v ...interface{}) {
}


