package plugins

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func postgresCheck(plg *Plugins, f *funcDesc) {
	crack(f.name, plg, f.user, f.pass, postgresConn)
}

func postgresConn(plg *Plugins, user, pass string) (ok int) {
	ok = OKNext
	source := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		plg.TargetIp, plg.TargetPort, user, pass)
	db, err := sql.Open("postgres", source)
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
	err = db.Ping()
	if err == nil {
		ok = OKDone
	} else {
		if strings.Contains(err.Error(), "password authentication") {
			return
		}
		//	color.Red(err.Error() + user + pass)
		ok = OKStop
	}
	return
}
