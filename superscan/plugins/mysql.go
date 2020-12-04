package plugins

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"sync"
	"time"
)

func mysqlCheck(plg *Plugins) interface{} {
	if !plg.NoTrust && plg.TargetPort != "3306" {
		return nil
	}
	plg.Mysql = make([]Account, 0)
	wg := sync.WaitGroup{}
	wg.Add(len(mysqlUsername))
	limiter := make(chan int, plg.Worker)
	ctx, cancel := context.WithCancel(context.TODO())
	for _, user := range mysqlUsername {
		limiter <- 1
		go func(username string) {
			defer func() {
				<-limiter
				wg.Done()
			}()
			for _, pass := range mysqlPassword {
				select {
				case <-ctx.Done():
					return
				default:
				}
				pass = strings.Replace(pass, "%user%", username, -1)
				plg.DescCallback(fmt.Sprintf("Cracking mysql %s:%s %s/%s",
					plg.TargetIp, plg.TargetPort, username, pass))
				ok := MysqlConn(plg, username, pass)
				switch ok {
				case OKDone:
					//密码正确一次退出
					plg.locker.Lock()
					plg.Mysql = append(plg.Mysql, Account{Username: username, Password: pass})
					plg.locker.Unlock()
					plg.highLight = true
					cancel()
					return
				case OKWait:
					//太快了服务器限制
					color.Red("[mysql]爆破频率太快服务器受限，建议降低参数'plugin-worker'数值影响主机:%s:%s", plg.TargetIp, plg.TargetPort)
					cancel()
					return
				case OKTimeOut:
					color.Red("[mysql]爆破过程中连接超时，建议提高参数'timeout'数值影响主机:%s:%s", plg.TargetIp, plg.TargetPort)
					cancel()
					return
				case OKStop:
					//非协议退出
					cancel()
					return
				default:
					//密码错误.OKNext
					plg.TargetProtocol = "[mysql]"
				}
			}
		}(user)
	}
	wg.Wait()
	//未找到密码
	if plg.TargetProtocol != "" {
		return &plg.Mysql
	}
	return nil
}

func MysqlConn(plg *Plugins, user string, pass string) (ok int) {
	plg.RateWait(plg.RateLimiter) //爆破限制
	ok = OKNext
	db, err := sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?timeout=%dms&readTimeout=%dms",
			user, pass, plg.TargetIp, plg.TargetPort, "mysql", plg.TimeOut, plg.TimeOut))
	if err != nil {
		color.Green(err.Error())
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
			ok = OKStop
			return
		}
		//非协议或受限制
		ok = OKStop
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Duration(plg.TimeOut) * time.Millisecond)
	err = db.Ping()
	if err == nil {
		ok = OKDone
	} else {
		//非协议
		ok = OKStop
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
