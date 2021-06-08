package plugins

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	"github.com/b1gcat/DarkEye/superscan/dic"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func mssqlCheck(s *Service) {
	s.crack()
}

func mssqlConn(_ context.Context, s *Service, user, pass string) (ok int) {
	ok = OKNext
	source := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;encrypt=disable;timeout=%dms",
		s.parent.TargetIp, user, pass, s.parent.TargetPort, time.Duration(Config.TimeOut)*time.Millisecond)
	db, err := sql.Open("mssql", source)
	if err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			//防火墙连接限制
			ok = OKWait
			return
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			ok = OKTimeOut
			return
		}
		//非协议
		ok = OKStop
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Duration(Config.TimeOut) * time.Millisecond)
	if err = db.Ping(); err == nil {
		ok = OKDone
	} else {
		if strings.Contains(err.Error(), "login error: mssql:") {
			return
		}
		color.Red(err.Error())
	}
	return
}

func init() {
	services["mssql"] = Service{
		name:    "mssql",
		port:    "1433",
		user:    dic.DIC_USERNAME_SQLSERVER,
		pass:    dic.DIC_PASSWORD_SQLSERVER,
		check:   mssqlCheck,
		connect: mssqlConn,
		thread:  1,
	}
}
