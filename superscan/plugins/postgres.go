package plugins

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func postgresCheck(s *Service) {
	s.crack()
}

func postgresConn(_ context.Context, s *Service, user, pass string) (ok int) {
	ok = OKNext
	source := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		s.parent.TargetIp, s.parent.TargetPort, user, pass)
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
	db.SetConnMaxLifetime(time.Duration(Config.TimeOut) * time.Millisecond)
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

func init() {
	services["postgres"] = Service{
		name:    "postgres",
		port:    "5432",
		user:    dic.DIC_USERNAME_POSTGRESQL,
		pass:    dic.DIC_PASSWORD_POSTGRESQL,
		check:   postgresCheck,
		connect: postgresConn,
		thread:  1,
	}
}
