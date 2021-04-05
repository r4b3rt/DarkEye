package plugins

import (
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func mssqlCheck(plg *Plugins, f *funcDesc) {
	crack(f.name, plg, f.user, f.pass, mssqlConn)
}

func mssqlConn(plg *Plugins, user, pass string) (ok int) {
	ok = OKNext
	source := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;encrypt=disable;timeout=%dms",
		plg.TargetIp, user, pass, plg.TargetPort, time.Duration(plg.TimeOut)*time.Millisecond)
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
	db.SetConnMaxLifetime(time.Duration(plg.TimeOut) * time.Millisecond)
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
